/*
Copyright 2023 Vesoft Inc.

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

package annotation

const (
	AnnDeploymentRevision = "deployment.kubernetes.io/revision"
	// AnnPVCDeferDeletingKey is pvc defer deletion annotation key used in PVC for defer deleting PVC
	AnnPVCDeferDeletingKey = "nebula-graph-ng.io/pvc-defer-deleting"
	// AnnPodNameKey is pod name annotation key used in PV/PVC for synchronizing nebula cluster meta info
	AnnPodNameKey = "nebula-graph-ng.io/pod-name"
	// AnnLastSyncTimestampKey is annotation key to indicate the last timestamp the operator sync the workload
	AnnLastSyncTimestampKey = "nebula-graph-ng.io/sync-timestamp"
	// AnnHaModeKey is annotation key to indicate whether in HA mode
	AnnHaModeKey = "nebula-graph-ng.io/ha-mode"
	// AnnLastReplicas is annotation key to indicate the last replicas
	AnnLastReplicas = "nebula-graph-ng.io/last-replicas"
	// AnnLastAppliedDynamicFlagsKey is annotation key to indicate the last applied custom dynamic flags
	AnnLastAppliedDynamicFlagsKey = "nebula-graph-ng.io/last-applied-dynamic-flags"
	// AnnLastAppliedConfigKey is annotation key to indicate the last applied configuration
	AnnLastAppliedConfigKey = "nebula-graph-ng.io/last-applied-configuration"
	// AnnPodConfigMapHash is pod configmap hash key to update configuration dynamically
	AnnPodConfigMapHash = "nebula-graph-ng.io/cm-hash"
	// AnnPvReclaimKey is annotation key that indicate whether reclaim persistent volume
	AnnPvReclaimKey = "nebula-graph-ng.io/enable-pv-reclaim"

	// AnnHaModeVal is annotation value to indicate whether in HA mode
	AnnHaModeVal = "true"

	// AnnDeleteProtection is an annotation key used to prevent the deletion of a nebula cluster that has been annotated by this key
	AnnDeleteProtection = "nebula-graph-ng.io/delete-protection"

	// AnnDeleteProtectionVal is annotation value to indicate whether nebula cluster is protected
	AnnDeleteProtectionVal = "true"
)

// IsInHaMode check whether in HA mode
func IsInHaMode(ann map[string]string) bool {
	if ann != nil {
		val, ok := ann[AnnHaModeKey]
		if ok && val == AnnHaModeVal {
			return true
		}
	}
	return false
}

func CopyAnnotations(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func IsDeleteProtected(ann map[string]string) bool {
	if ann != nil {
		val, ok := ann[AnnDeleteProtection]
		if ok && val == AnnDeleteProtectionVal {
			return true
		}
	}
	return false
}
