/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

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

package v1alpha1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
)

// ConvertTo converts this ClusterDefinition to the Hub version (v1).
func (r *ClusterDefinition) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*appsv1.ClusterDefinition)
	// ObjectMeta
	dst.ObjectMeta = r.ObjectMeta
	// Spec
	// Status
	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (r *ClusterDefinition) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*appsv1.ClusterDefinition)
	// ObjectMeta
	r.ObjectMeta = src.ObjectMeta
	// Spec
	// Status
	return nil
}
