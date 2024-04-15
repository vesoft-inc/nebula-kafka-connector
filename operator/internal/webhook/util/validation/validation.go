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

package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
)

const (
	minReplicasGraphdNotInHaMode   = 1
	minReplicasGraphdInHaMode      = 2
	minReplicasMetadNotInHaMode    = 1
	minReplicasMetadInHaMode       = 3
	minReplicasStoragedNotInHaMode = 1
	minReplicasStoragedInHaMode    = 3

	fmtNotHaModeErrorDetail = "should be at least %d not in HA mode"
	fmtHaModeErrorDetail    = "should be at least %d in HA mode"
)

// ValidateMinReplicas validates replicas min value
func ValidateMinReplicas(fldPath *field.Path, actualValue, minValue int, bHaMode bool) *field.Error {
	if actualValue < minValue {
		detail := fmt.Sprintf(fmtNotHaModeErrorDetail, minValue)
		if bHaMode {
			detail = fmt.Sprintf(fmtHaModeErrorDetail, minValue)
		}
		return field.Invalid(fldPath, actualValue, detail)
	}
	return nil
}

// ValidateMinReplicasGraphd validates replicas min value for Graphd
func ValidateMinReplicasGraphd(fldPath *field.Path, replicas int, bHaMode bool) (allErrs field.ErrorList) {
	minReplicas := minReplicasGraphdNotInHaMode
	if bHaMode {
		minReplicas = minReplicasGraphdInHaMode
	}
	if fieldErr := ValidateMinReplicas(fldPath, replicas, minReplicas, bHaMode); fieldErr != nil {
		allErrs = append(allErrs, fieldErr)
	}

	return allErrs
}

// ValidateMinReplicasMetad validates replicas min value for Metad
func ValidateMinReplicasMetad(fldPath *field.Path, replicas int, bHaMode bool) (allErrs field.ErrorList) {
	minReplicas := minReplicasMetadNotInHaMode
	if bHaMode {
		minReplicas = minReplicasMetadInHaMode
	}

	if fieldErr := ValidateMinReplicas(fldPath, replicas, minReplicas, bHaMode); fieldErr != nil {
		allErrs = append(allErrs, fieldErr)
	}
	if replicas&1 == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath, replicas, "should be odd number"))
	}

	return allErrs
}

// ValidateMinReplicasStoraged validates replicas min value for Storaged
func ValidateMinReplicasStoraged(fldPath *field.Path, replicas int, bHaMode bool) (allErrs field.ErrorList) {
	minReplicas := minReplicasStoragedNotInHaMode
	if bHaMode {
		minReplicas = minReplicasStoragedInHaMode
	}

	if fieldErr := ValidateMinReplicas(fldPath, replicas, minReplicas, bHaMode); fieldErr != nil {
		allErrs = append(allErrs, fieldErr)
	}

	return allErrs
}

func IsNebulaClusterHA(nc *v2alpha1.NebulaCluster) bool {
	if pointer.Int32Deref(nc.Spec.Graphd.Replicas, 0) < minReplicasGraphdInHaMode {
		return false
	}
	if pointer.Int32Deref(nc.Spec.Storaged.Replicas, 0) < minReplicasStoragedInHaMode {
		return false
	}
	return true
}

func IsNebulaMetadHA(nm *v2alpha1.NebulaMetad) bool {
	return pointer.Int32Deref(nm.Spec.Replicas, 0) < minReplicasMetadInHaMode
}
