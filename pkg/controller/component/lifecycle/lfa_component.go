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

package lifecycle

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type postProvision struct {
	// namespace   string
	// clusterName string
	// compName    string
}

var _ lifecycleAction = &postProvision{}

func (a *postProvision) name() string {
	return "postProvision"
}

func (a *postProvision) parameters(ctx context.Context, cli client.Reader) (map[string]string, error) {
	// TODO: refine and implement
	// The container executing this action has access to following variables:
	//
	// - KB_CLUSTER_POD_IP_LIST: Comma-separated list of the cluster's pod IP addresses (e.g., "podIp1,podIp2").
	// - KB_CLUSTER_POD_NAME_LIST: Comma-separated list of the cluster's pod names (e.g., "pod1,pod2").
	// - KB_CLUSTER_POD_HOST_NAME_LIST: Comma-separated list of host names, each corresponding to a pod in
	//   KB_CLUSTER_POD_NAME_LIST (e.g., "hostName1,hostName2").
	// - KB_CLUSTER_POD_HOST_IP_LIST: Comma-separated list of host IP addresses, each corresponding to a pod in
	//   KB_CLUSTER_POD_NAME_LIST (e.g., "hostIp1,hostIp2").
	//
	// - KB_CLUSTER_COMPONENT_POD_NAME_LIST: Comma-separated list of all pod names within the component
	//   (e.g., "pod1,pod2").
	// - KB_CLUSTER_COMPONENT_POD_IP_LIST: Comma-separated list of pod IP addresses,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "podIp1,podIp2").
	// - KB_CLUSTER_COMPONENT_POD_HOST_NAME_LIST: Comma-separated list of host names for each pod,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "hostName1,hostName2").
	// - KB_CLUSTER_COMPONENT_POD_HOST_IP_LIST: Comma-separated list of host IP addresses for each pod,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "hostIp1,hostIp2").
	//
	// - KB_CLUSTER_COMPONENT_LIST: Comma-separated list of all cluster components (e.g., "comp1,comp2").
	// - KB_CLUSTER_COMPONENT_DELETING_LIST: Comma-separated list of components that are currently being deleted
	//   (e.g., "comp1,comp2").
	// - KB_CLUSTER_COMPONENT_UNDELETED_LIST: Comma-separated list of components that are not being deleted
	//   (e.g., "comp1,comp2").
	return nil, nil
}

type preTerminate struct {
	// namespace   string
	// clusterName string
	// compName    string
}

var _ lifecycleAction = &preTerminate{}

func (a *preTerminate) name() string {
	return "preTerminate"
}

func (a *preTerminate) parameters(ctx context.Context, cli client.Reader) (map[string]string, error) {
	// TODO: refine and implement
	// The container executing this action has access to following variables:
	//
	// - KB_CLUSTER_POD_IP_LIST: Comma-separated list of the cluster's pod IP addresses (e.g., "podIp1,podIp2").
	// - KB_CLUSTER_POD_NAME_LIST: Comma-separated list of the cluster's pod names (e.g., "pod1,pod2").
	// - KB_CLUSTER_POD_HOST_NAME_LIST: Comma-separated list of host names, each corresponding to a pod in
	//   KB_CLUSTER_POD_NAME_LIST (e.g., "hostName1,hostName2").
	// - KB_CLUSTER_POD_HOST_IP_LIST: Comma-separated list of host IP addresses, each corresponding to a pod in
	//   KB_CLUSTER_POD_NAME_LIST (e.g., "hostIp1,hostIp2").
	//
	// - KB_CLUSTER_COMPONENT_POD_NAME_LIST: Comma-separated list of all pod names within the component
	//   (e.g., "pod1,pod2").
	// - KB_CLUSTER_COMPONENT_POD_IP_LIST: Comma-separated list of pod IP addresses,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "podIp1,podIp2").
	// - KB_CLUSTER_COMPONENT_POD_HOST_NAME_LIST: Comma-separated list of host names for each pod,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "hostName1,hostName2").
	// - KB_CLUSTER_COMPONENT_POD_HOST_IP_LIST: Comma-separated list of host IP addresses for each pod,
	//   matching the order of pods in KB_CLUSTER_COMPONENT_POD_NAME_LIST (e.g., "hostIp1,hostIp2").
	//
	// - KB_CLUSTER_COMPONENT_LIST: Comma-separated list of all cluster components (e.g., "comp1,comp2").
	// - KB_CLUSTER_COMPONENT_DELETING_LIST: Comma-separated list of components that are currently being deleted
	//   (e.g., "comp1,comp2").
	// - KB_CLUSTER_COMPONENT_UNDELETED_LIST: Comma-separated list of components that are not being deleted
	//   (e.g., "comp1,comp2").
	//
	// - KB_CLUSTER_COMPONENT_IS_SCALING_IN: Indicates whether the component is currently scaling in.
	//   If this variable is present and set to "true", it denotes that the component is undergoing a scale-in operation.
	//   During scale-in, data rebalancing is necessary to maintain cluster integrity.
	//   Contrast this with a cluster deletion scenario where data rebalancing is not required as the entire cluster
	//   is being cleaned up.
	return nil, nil
}
