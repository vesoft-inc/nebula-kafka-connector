/*
Copyright 2024 Vesoft Inc.

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

package nebulametad

import (
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	apivalidation "k8s.io/kubernetes/pkg/apis/core/validation"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/webhook/util/validation"
)

// validateNebulaMetadReplicas validates Metad replicas.
func validateNebulaMetadReplicas(nm *v2alpha1.NebulaMetad) (allErrs field.ErrorList) {
	replicas := *nm.Spec.Replicas
	bHaMode := annotation.IsInHaMode(nm.Annotations)

	allErrs = append(allErrs, validation.ValidateMinReplicasMetad(
		field.NewPath("spec").Child("replicas"),
		int(replicas),
		bHaMode,
	)...)

	return allErrs
}

// validateNebulaMetadCreate validates a NebulaMetad on create.
func validateNebulaMetadCreate(nm *v2alpha1.NebulaMetad) (allErrs field.ErrorList) {
	allErrs = append(allErrs, validateNebulaMetadReplicas(nm)...)

	return allErrs
}

// ValidateNebulaCluster validates a NebulaMetad on Update.
func validateNebulaMetadUpdate(nm, oldNM *v2alpha1.NebulaMetad) (allErrs field.ErrorList) {
	name := nm.Name
	namespace := nm.Namespace

	klog.Infof("receive admission with resource [%s/%s], GVK %s, operation %s", namespace, name,
		nm.GroupVersionKind().String(), admissionv1.Update)

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaUpdate(
		&nm.ObjectMeta,
		&oldNM.ObjectMeta,
		field.NewPath("metadata"),
	)...)

	if !validation.IsNebulaMetadHA(oldNM) {
		allErrs = append(allErrs, apivalidation.ValidateImmutableAnnotation(
			nm.Annotations[annotation.AnnHaModeKey],
			oldNM.Annotations[annotation.AnnHaModeKey],
			annotation.AnnHaModeKey,
			field.NewPath("metadata"),
		)...)
	}

	allErrs = append(allErrs, validateNebulaMetadReplicas(nm)...)

	return allErrs
}

// validateNebulaMetadDelete validates a NebulaMetad on Delete.
func validateNebulaMetadDelete(nm *v2alpha1.NebulaMetad) (allErrs field.ErrorList) {
	name := nm.Name
	namespace := nm.Namespace

	klog.Infof("receive admission with resource [%s/%s], GVK %s, operation %s", namespace, name,
		nm.GroupVersionKind().String(), admissionv1.Delete)

	if annotation.IsDeleteProtected(nm.Annotations) {
		fldPath := field.NewPath("metadata")
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("annotations").Key(annotation.AnnDeleteProtection),
			"protected metad cannot be deleted"))
	}
	return allErrs
}
