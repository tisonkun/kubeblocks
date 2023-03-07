/*
Copyright ApeCloud, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package extensions

import (
	"context"
	"runtime"
	"time"

	ctrlerihandler "github.com/authzed/controller-idioms/handler"
	"github.com/spf13/viper"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	extensionsv1alpha1 "github.com/apecloud/kubeblocks/apis/extensions/v1alpha1"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
)

// AddonReconciler reconciles a Addon object
type AddonReconciler struct {
	client.Client
	Scheme     *k8sruntime.Scheme
	Recorder   record.EventRecorder
	RestConfig *rest.Config
}

const (
	// settings keys
	maxConcurrentReconcilesKey = "MAXCONCURRENTRECONCILES_ADDON"
	addonJobImagePullPolicyKey = "ADDON_JOB_IMAGE_PULL_POLICY"
	addonSANameKey             = "KUBEBLOCKS_ADDON_SA_NAME"
)

func init() {
	viper.SetDefault(maxConcurrentReconcilesKey, runtime.NumCPU()*2)
}

// +kubebuilder:rbac:groups=extensions.kubeblocks.io,resources=addons,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions.kubeblocks.io,resources=addons/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=extensions.kubeblocks.io,resources=addons/finalizers,verbs=update

// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;deletecollection

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AddonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqCtx := intctrlutil.RequestCtx{
		Ctx:      ctx,
		Req:      req,
		Log:      log.FromContext(ctx).WithValues("addon", req.NamespacedName),
		Recorder: r.Recorder,
	}

	buildStageCtx := func(next ...ctrlerihandler.Handler) stageCtx {
		return stageCtx{
			reqCtx:     &reqCtx,
			reconciler: r,
			next:       ctrlerihandler.Handlers(next).MustOne(),
		}
	}

	fetchNDeletionCheckStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&fetchNDeletionCheckStage{
			stageCtx: buildStageCtx(next...),
			deletionStage: deletionStage{
				stageCtx: buildStageCtx(ctrlerihandler.NoopHandler),
			},
		})
	}

	genIDProceedStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&genIDProceedCheckStage{stageCtx: buildStageCtx(next...)})
	}

	installableCheckStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&installableCheckStage{stageCtx: buildStageCtx(next...)})
	}

	autoInstallCheckStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&autoInstallCheckStage{stageCtx: buildStageCtx(next...)})
	}

	progressingStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&progressingHandler{stageCtx: buildStageCtx(next...)})
	}

	terminalStateStageBuilder := func(next ...ctrlerihandler.Handler) ctrlerihandler.Handler {
		return ctrlerihandler.NewTypeHandler(&terminalStateStage{stageCtx: buildStageCtx(next...)})
	}

	handlers := ctrlerihandler.Chain(
		fetchNDeletionCheckStageBuilder,
		genIDProceedStageBuilder,
		installableCheckStageBuilder,
		autoInstallCheckStageBuilder,
		progressingStageBuilder,
		terminalStateStageBuilder,
	).Handler("")

	handlers.Handle(ctx)
	res, ok := reqCtx.Ctx.Value(resultValueKey).(*ctrl.Result)
	if ok && res != nil {
		err, ok := reqCtx.Ctx.Value(errorValueKey).(error)
		if ok {
			return *res, err
		}
		return *res, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AddonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&extensionsv1alpha1.Addon{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: viper.GetInt(maxConcurrentReconcilesKey),
		}).
		Owns(&batchv1.Job{}). // TODO: cannot owns a namespaced object
		Complete(r)
}

func (r *AddonReconciler) deleteExternalResources(reqCtx intctrlutil.RequestCtx, addon *extensionsv1alpha1.Addon) (*ctrl.Result, error) {

	if addon.Annotations != nil && addon.Annotations[NoDeleteJobs] == "true" {
		return nil, nil
	}

	deleteJobIfExist := func(jobName string) error {
		key := client.ObjectKey{
			Namespace: viper.GetString("CM_NAMESPACE"),
			Name:      jobName,
		}
		job := &batchv1.Job{}
		if err := r.Get(reqCtx.Ctx, key, job); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		if !job.DeletionTimestamp.IsZero() {
			return nil
		}

		if err := r.Delete(reqCtx.Ctx, job); err != nil {
			return client.IgnoreNotFound(err)
		}
		return nil
	}
	for _, j := range []string{getInstallJobName(addon), getUninstallJobName(addon)} {
		if err := deleteJobIfExist(j); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

const (
	resultValueKey  = "result"
	errorValueKey   = "err"
	operandValueKey = "operand"
)

type stageCtx struct {
	reqCtx     *intctrlutil.RequestCtx
	reconciler *AddonReconciler
	next       ctrlerihandler.Handler
}

func (r *stageCtx) setReconciled() {
	res, err := intctrlutil.Reconciled()
	r.updateResultNErr(&res, err)
}

func (r *stageCtx) setRequeue() {
	res, err := intctrlutil.Requeue(r.reqCtx.Log, "")
	r.updateResultNErr(&res, err)
}

func (r *stageCtx) setRequeueAfter(duration time.Duration, msg string) {
	res, err := intctrlutil.RequeueAfter(time.Second, r.reqCtx.Log, msg)
	r.updateResultNErr(&res, err)
}

func (r *stageCtx) setRequeueWithErr(err error, msg string) {
	res, err := intctrlutil.CheckedRequeueWithError(err, r.reqCtx.Log, msg)
	r.updateResultNErr(&res, err)
}

func (r *stageCtx) updateResultNErr(res *ctrl.Result, err error) {
	r.reqCtx.UpdateCtxValue(errorValueKey, err)
	r.reqCtx.UpdateCtxValue(resultValueKey, res)
}

func (r *stageCtx) doReturn() (*ctrl.Result, error) {
	res, _ := r.reqCtx.Ctx.Value(resultValueKey).(*ctrl.Result)
	err, _ := r.reqCtx.Ctx.Value(errorValueKey).(error)
	return res, err
}

func (r *stageCtx) process(processor func(*extensionsv1alpha1.Addon)) {
	res, _ := r.doReturn()
	if res != nil {
		return
	}
	addon := r.reqCtx.Ctx.Value(operandValueKey).(*extensionsv1alpha1.Addon)
	processor(addon)
}
