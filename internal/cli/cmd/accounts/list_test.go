/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

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

package accounts

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
	clientfake "k8s.io/client-go/rest/fake"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"

	"github.com/apecloud/kubeblocks/internal/cli/testing"
	"github.com/apecloud/kubeblocks/internal/cli/types"
	"github.com/apecloud/kubeblocks/internal/sqlchannel"
)

var _ = Describe("List Account Options", func() {
	const (
		namespace   = "test"
		clusterName = "apple"
	)

	var (
		streams genericclioptions.IOStreams
		tf      *cmdtesting.TestFactory
		cluster = testing.FakeCluster(clusterName, namespace)
		pods    = testing.FakePods(3, namespace, clusterName)
	)

	BeforeEach(func() {
		streams, _, _, _ = genericclioptions.NewTestIOStreams()
		tf = testing.NewTestFactory(namespace)
		codec := scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)
		httpResp := func(obj runtime.Object) *http.Response {
			return &http.Response{StatusCode: http.StatusOK, Header: cmdtesting.DefaultHeader(), Body: cmdtesting.ObjBody(codec, obj)}
		}

		tf.UnstructuredClient = &clientfake.RESTClient{
			GroupVersion:         schema.GroupVersion{Group: types.AppsAPIGroup, Version: types.AppsAPIVersion},
			NegotiatedSerializer: resource.UnstructuredPlusDefaultContentConfig().NegotiatedSerializer,
			Client: clientfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
				urlPrefix := "/api/v1/namespaces/" + namespace
				mapping := map[string]*http.Response{
					urlPrefix + "/pods":                       httpResp(pods),
					urlPrefix + "/pods/" + pods.Items[0].Name: httpResp(&pods.Items[0]),
				}
				return mapping[req.URL.Path], nil
			}),
		}

		tf.Client = tf.UnstructuredClient
		tf.FakeDynamicClient = testing.FakeDynamicClient(cluster, testing.FakeClusterDef(), testing.FakeClusterVersion())
	})

	AfterEach(func() {
		tf.Cleanup()
	})

	Context("new options", func() {
		It("new option", func() {
			o := NewListUserOptions(tf, streams)
			Expect(o).ShouldNot(BeNil())
			Expect(o.AccountBaseOptions).ShouldNot(BeNil())
			Expect(o.AccountBaseOptions.AccountOp).Should(Equal(sqlchannel.ListUsersOp))
		})

		It("validate options", func() {
			o := NewListUserOptions(tf, streams)
			Expect(o).ShouldNot(BeNil())
			args := []string{}
			Expect(o.Validate(args)).Should(MatchError(errClusterNameorInstName))

			// add two elements
			By("add two args")
			args = []string{"foo", "bar"}
			Expect(o.Validate(args)).Should(MatchError(errClusterNameNum))

			// add one element
			By("add one more args, should fail")
			args = []string{"foo"}
			Expect(o.Validate(args)).Should(Succeed())

			// set pod name
			o.PodName = "pod1"
			Expect(o.Validate(args)).Should(MatchError(errClusterNameorInstName))
			// set component name
			o.ComponentName = "foo-component"
			Expect(o.Validate(args)).Should(MatchError(errCompNameOrInstName))
			// set both
			o.PodName = ""
			Expect(o.Validate(args)).Should(Succeed())
		})

		It("complete option", func() {
			o := NewListUserOptions(tf, streams)
			Expect(o).ShouldNot(BeNil())
			o.PodName = pods.Items[0].Name
			o.ClusterName = clusterName
			Expect(o.Complete(tf)).Should(Succeed())

			Expect(o.Client).ShouldNot(BeNil())
			Expect(o.Dynamic).ShouldNot(BeNil())
			Expect(o.Namespace).Should(Equal(namespace))
			Expect(o.Pod.Name).Should(Equal(o.PodName))
		})
	})
})
