/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package component

import (
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/builder"
	kbagent "github.com/apecloud/kubeblocks/pkg/kbagent"
	"github.com/apecloud/kubeblocks/pkg/kbagent/proto"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

const (
	kbAgentContainerName     = "kbagent"
	kbAgentInitContainerName = "init-kbagent"
	kbAgentCommand           = "/bin/kbagent"
	kbAgentPortName          = "http"

	kbAgentSharedMountPath      = "/kubeblocks"
	kbAgentCommandOnSharedMount = "/kubeblocks/kbagent"
)

var (
	sharedVolumeMount = corev1.VolumeMount{Name: "kubeblocks", MountPath: kbAgentSharedMountPath}
)

func buildKBAgentContainer(synthesizedComp *SynthesizedComponent) error {
	if synthesizedComp.LifecycleActions == nil {
		return nil
	}

	envVars, err := buildKBAgentStartupEnv(synthesizedComp)
	if err != nil {
		return err
	}

	port := 3501 // TODO: port

	container := builder.NewContainerBuilder(kbAgentContainerName).
		SetImage(viper.GetString(constant.KBToolsImage)).
		SetImagePullPolicy(corev1.PullIfNotPresent).
		AddCommands(kbAgentCommand).
		AddArgs("--port", strconv.Itoa(port)).
		AddEnv(envVars...).
		AddPorts(corev1.ContainerPort{
			ContainerPort: int32(port),
			Name:          kbAgentPortName,
			Protocol:      "TCP",
		}).
		SetStartupProbe(corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				TCPSocket: &corev1.TCPSocketAction{Port: intstr.FromInt(port)},
			}}).
		GetObject()

	if err = adaptKBAgentIfCustomImageNContainerDefined(synthesizedComp, container); err != nil {
		return err
	}

	// set kb-agent container ports to host network
	if synthesizedComp.HostNetwork != nil {
		if synthesizedComp.HostNetwork.ContainerPorts == nil {
			synthesizedComp.HostNetwork.ContainerPorts = make([]appsv1alpha1.HostNetworkContainerPort, 0)
		}
		synthesizedComp.HostNetwork.ContainerPorts = append(
			synthesizedComp.HostNetwork.ContainerPorts,
			appsv1alpha1.HostNetworkContainerPort{
				Container: container.Name,
				Ports:     []string{kbAgentPortName},
			})
	}

	synthesizedComp.PodSpec.Containers = append(synthesizedComp.PodSpec.Containers, *container)
	return nil
}

func buildKBAgentStartupEnv(synthesizedComp *SynthesizedComponent) ([]corev1.EnvVar, error) {
	var (
		actions []proto.Action
		probes  []proto.Probe
	)

	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.PostProvision, "postProvision"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.PreTerminate, "preTerminate"); a != nil {
		actions = append(actions, *a)
	}
	if synthesizedComp.LifecycleActions.Switchover != nil {
		if a := buildAction4KBAgentLow(synthesizedComp.LifecycleActions.Switchover.WithoutCandidate, "switchover"); a != nil {
			actions = append(actions, *a)
		}
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.MemberJoin, "memberJoin"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.MemberLeave, "memberLeave"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.Readonly, "readonly"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.Readwrite, "readwrite"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.DataDump, "dataDump"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.DataLoad, "dataLoad"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.Reconfigure, "reconfigure"); a != nil {
		actions = append(actions, *a)
	}
	if a := buildAction4KBAgent(synthesizedComp.LifecycleActions.AccountProvision, "accountProvision"); a != nil {
		actions = append(actions, *a)
	}

	if a, p := buildProbe4KBAgent(synthesizedComp.LifecycleActions.RoleProbe, "roleProbe"); a != nil && p != nil {
		actions = append(actions, *a)
		probes = append(probes, *p)
	}

	return kbagent.BuildEnvVars(actions, probes)
}

func buildAction4KBAgent(handler *appsv1alpha1.LifecycleActionHandler, name string) *proto.Action {
	if handler == nil {
		return nil
	}
	return buildAction4KBAgentLow(handler.CustomHandler, name)
}

func buildAction4KBAgentLow(action *appsv1alpha1.Action, name string) *proto.Action {
	if action == nil || action.Exec == nil {
		return nil
	}
	a := &proto.Action{
		Name: name,
		Exec: &proto.ExecAction{
			Commands: action.Exec.Command,
			Args:     action.Exec.Args,
			// Env:       action.Exec.Env,
			Container: action.Exec.Container,
		},
		TimeoutSeconds: action.TimeoutSeconds,
	}
	if action.RetryPolicy != nil {
		a.RetryPolicy = &proto.RetryPolicy{
			MaxRetries:    action.RetryPolicy.MaxRetries,
			RetryInterval: action.RetryPolicy.RetryInterval,
		}
	}
	return a
}

func buildProbe4KBAgent(probe *appsv1alpha1.Probe, name string) (*proto.Action, *proto.Probe) {
	if probe == nil || probe.Exec == nil {
		return nil, nil
	}
	a := buildAction4KBAgentLow(&probe.Action, name)
	p := &proto.Probe{
		Action:              name,
		InitialDelaySeconds: probe.InitialDelaySeconds,
		PeriodSeconds:       probe.PeriodSeconds,
		SuccessThreshold:    probe.SuccessThreshold,
		FailureThreshold:    probe.FailureThreshold,
		ReportPeriodSeconds: nil, // TODO
	}
	return a, p
}

func adaptKBAgentIfCustomImageNContainerDefined(synthesizedComp *SynthesizedComponent, container *corev1.Container) error {
	image, _, err := customExecActionImageNContainer(synthesizedComp)
	if err != nil {
		return err
	}
	if len(image) == 0 {
		return nil
	}

	// init-container to copy binaries to the shared mount point /kubeblocks
	initContainer := buildKBAgentInitContainer()
	synthesizedComp.PodSpec.InitContainers = append(synthesizedComp.PodSpec.InitContainers, *initContainer)

	container.Image = image
	container.Command[0] = kbAgentCommandOnSharedMount
	container.VolumeMounts = append(container.VolumeMounts, sharedVolumeMount)

	// TODO: exec container resources
	return nil
}

func customExecActionImageNContainer(synthesizedComp *SynthesizedComponent) (string, string, error) {
	if synthesizedComp.LifecycleActions == nil {
		return "", "", nil
	}

	handlers := []*appsv1alpha1.LifecycleActionHandler{
		synthesizedComp.LifecycleActions.PostProvision,
		synthesizedComp.LifecycleActions.PreTerminate,
		synthesizedComp.LifecycleActions.MemberJoin,
		synthesizedComp.LifecycleActions.MemberLeave,
		synthesizedComp.LifecycleActions.Readonly,
		synthesizedComp.LifecycleActions.Readwrite,
		synthesizedComp.LifecycleActions.DataDump,
		synthesizedComp.LifecycleActions.DataLoad,
		synthesizedComp.LifecycleActions.Reconfigure,
		synthesizedComp.LifecycleActions.AccountProvision,
	}
	if synthesizedComp.LifecycleActions.RoleProbe != nil && synthesizedComp.LifecycleActions.RoleProbe.Exec != nil {
		handlers = append(handlers, &appsv1alpha1.LifecycleActionHandler{
			CustomHandler: &synthesizedComp.LifecycleActions.RoleProbe.Action,
		})
	}

	var image, container string
	for _, handler := range handlers {
		if handler == nil || handler.CustomHandler == nil || handler.CustomHandler.Exec == nil {
			continue
		}
		if handler.CustomHandler.Exec.Image != "" {
			if len(image) > 0 && image != handler.CustomHandler.Exec.Image {
				return "", "", fmt.Errorf("only one exec image is allowed in lifecycle actions")
			}
			image = handler.CustomHandler.Exec.Image
		}
		if handler.CustomHandler.Exec.Container != "" {
			if len(container) > 0 && container != handler.CustomHandler.Exec.Container {
				return "", "", fmt.Errorf("only one exec container is allowed in lifecycle actions")
			}
			container = handler.CustomHandler.Exec.Container
		}
	}

	var c *corev1.Container
	if len(container) > 0 {
		for _, cc := range synthesizedComp.PodSpec.Containers {
			if cc.Name == container {
				c = &cc
				break
			}
		}
		if c == nil {
			return "", "", fmt.Errorf("exec container %s not found", container)
		}
	}
	if len(image) > 0 && len(container) > 0 {
		if c.Image == image {
			return image, container, nil
		}
		return "", "", fmt.Errorf("exec image and container must be the same")
	}
	if len(image) == 0 && len(container) > 0 {
		image = c.Image
	}
	return image, container, nil
}

func buildKBAgentInitContainer() *corev1.Container {
	return builder.NewContainerBuilder(kbAgentInitContainerName).
		SetImage(viper.GetString(constant.KBToolsImage)).
		SetImagePullPolicy(corev1.PullIfNotPresent).
		AddCommands([]string{"cp", "-r", kbAgentCommand, "/bin/curl", kbAgentSharedMountPath + "/"}...).
		AddVolumeMounts(sharedVolumeMount).
		GetObject()
}
