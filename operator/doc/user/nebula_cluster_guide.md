## Install NebulaGraph cluster with helm

Please install [nebula-operator](install_guide.md) before installing NebulaGraph cluster.

### Get the repo

```shell script
# If you have already added it, please skip this step.
$ helm repo add nebula-ng-operator https://vesoft-inc.github.io/nebula-ng-tools/operator/charts
$ helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

### Install with helm

#### Install nebula metad
If you have an existing nebula 5.0 metad cluster in the namespace you want to deploy the nebula cluster in, please skip this step.
```shell script
export NEBULA_METAD_NAME=nebula-metad   # the name of the nebula metad cluster
export NEBULA_METAD_NAMESPACE=nebula    # the namespace you want to install the nebula metad to
export STORAGE_CLASS_NAME=gp2           # the storage class for the nebula metad

$ kubectl create namespace "${NEBULA_METAD_NAMESPACE}" # If you have already created it, please skip this step.
$ kubectl create secret generic nebula --from-literal=username=root --from-literal=password=<new_password_for_root_user> -n ${NEBULA_METAD_NAMESPACE}
$ helm install "${NEBULA_METAD_NAME}" nebula-ng-operator/nebula-metad \
    --namespace="${NEBULA_METAD_NAMESPACE}" \
    --set nameOverride=${NEBULA_METAD_NAME} \
    --set nebula.storageClassName="${STORAGE_CLASS_NAME}"

# Please wait a while for the nebula metad to be ready.
$ kubectl -n "${NEBULA_METAD_NAMESPACE}" get pod -l "app.kubernetes.io/cluster=${NEBULA_METAD_NAME}"
NAME                   READY   STATUS    RESTARTS   AGE
nebula-metad-metad-0   1/1     Running   0          69s
nebula-metad-metad-1   1/1     Running   0          69s
nebula-metad-metad-2   1/1     Running   0          69s
```

#### Install nebula cluster
```shell script
export NEBULA_CLUSTER_NAME=nebula-cluster  # the name for nebula cluster
export NEBULA_CLUSTER_NAMESPACE=nebula     # the namespace you want to install the nebula cluster. Must be the same as the one for the metad that will manage this cluster.
export STORAGE_CLASS_NAME=gp2              # the storage class for the nebula cluster

$ helm install "${NEBULA_CLUSTER_NAME}" nebula-ng-operator/nebula-cluster \
    --namespace="${NEBULA_CLUSTER_NAMESPACE}" \
    --set nameOverride=${NEBULA_CLUSTER_NAME} \
    --set nebula.storageClassName="${STORAGE_CLASS_NAME}"

# Please wait a while for the cluster to be ready.
$ kubectl -n "${NEBULA_CLUSTER_NAMESPACE}" get pod -l "app.kubernetes.io/cluster=${NEBULA_CLUSTER_NAME}"
NAME                        READY   STATUS    RESTARTS   AGE
nebula-cluster-graphd-0     1/1     Running   0          2m33s
nebula-cluster-storaged-0   1/1     Running   0          2m33s
nebula-cluster-storaged-1   1/1     Running   0          2m33s
nebula-cluster-storaged-2   1/1     Running   0          2m33s
```

### Uninstall with helm

```shell
$ helm uninstall "${NEBULA_CLUSTER_NAME}" --namespace="${NEBULA_CLUSTER_NAMESPACE}"
```

Wait for the nebula cluster to finish uninstalling. Then uninstall nebula metad
```shell
$ helm uninstall "${NEBULA_METAD_NAME}" --namespace="${NEBULA_CLUSTER_NAMESPACE}"
```
**Warning** Only uninstall the nebula metad after uninstalling all nebula clusters it manages. Otherwise uninstallation will fail.

Delete the nebula secret
```shell
kubectl delete secret nebula -n"${NEBULA_METAD_NAMESPACE}"
```

### Optional: chart parameters

The following tables list the configurable parameters for the nebula metad and nebula cluster charts. Their default values are also listed.

**Global**
| Parameter                           | Description                                                                               | Default                                                                                               |
|:------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nameOverride`                      | Override the name of the chart                                                            | `nil`                                                                                                 |
| `nebula.version`                    | Nebula image tag                                                                          | `v3.6.0`                                                                                              |
| `nebula.imagePullPolicy`            | Nebula image pull policy                                                                  | `Always`                                                                                              |
| `nebula.enablePVReclaim`            | Flag to enable/disable PV reclaim while the Nebula cluster deleted                        | `false`                                                                                               |
| `nebula.storageClassName`           | StorageClass object name                                                                  | `""`                                                                                                  |
| `nebula.schedulerName`              | Scheduler for nebula component                                                            | `default-scheduler`                                                                                   |
| `nebula.topologySpreadConstraints`  | Topology spread constraints to control how Pods are spread across failure-domains         | `[]`                                                                                                  |
| `nebula.imagePullSecrets`           | The secret that contains the credentials for pulling the images                           | `[]`                                                                                                  |
| `credentialSecret`                  | The secret that contains the credentials for the nebula graph root user                   | "nebula"                                                                                              |



**Nebula Metad Chart Only**
| Parameter                           | Description                                                                               | Default                                                                                               |
|:------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nebula.metad.image`                | Metad container image without tag, and use `nebula.version` as tag                        | `vesoft/nebula-metad`                                                                                 |
| `nebula.metad.replicas`             | Metad replica number                                                                      | `3`                                                                                                   |
| `nebula.metad.env`                  | Metad container environment variables                                                     | `[]`                                                                                                  |
| `nebula.metad.resources`            | Metad resources                                                                           | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebula.metad.logVolume`            | Metad log volum                                                                           | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebula.metad.dataVolume`           | Metad data volume                                                                         | `[]`                                                                                                  |
| `nebula.metad.podLabels`            | Metad pod labels                                                                          | `{}`                                                                                                  |
| `nebula.metad.podAnnotations`       | Metad pod annotations                                                                     | `{}`                                                                                                  |
| `nebula.metad.SecurityContext`      | Metad security context                                                                    | `{}`                                                                                                  |
| `nebula.metad.nodeSelector`         | Metad nodeSelector                                                                        | `{}`                                                                                                  |
| `nebula.metad.tolerations`          | Metad pod tolerations                                                                     | `[]`                                                                                                  |
| `nebula.metad.affinity`             | Metad affinity                                                                            | `{}`                                                                                                  |
| `nebula.metad.readinessProbe`       | Metad pod readinessProbe                                                                  | `{}`                                                                                                  |
| `nebula.metad.livenessProbe`        | Metad pod livenessProbe                                                                   | `{}`                                                                                                  |
| `nebula.metad.initContainers`       | Metad pod init containers                                                                 | `[]`                                                                                                  |
| `nebula.metad.sidecarContainers`    | Metad pod sidecar containers                                                              | `[]`                                                                                                  |
| `nebula.metad.volumes`              | Metad pod volumes                                                                         | `[]`                                                                                                  |
| `nebula.metad.volumeMounts`         | Metad pod volume mounts                                                                   | `[]`                                                                                                  |



**Nebula Cluster Chart Only**
| Parameter                           | Description                                                                               | Default                                                                                               |
|:------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nebula.metadRef`                   | Reference to the metad that manages the nebula cluster                                    | `{"name": "metad-sample"}`                                                                            |
| `nebula.graphd.image`               | Graphd container image without tag, and use `nebula.version` as tag                       | `vesoft/nebula-graphd`                                                                                |
| `nebula.graphd.replicas`            | Graphd replica number                                                                     | `2`                                                                                                   |
| `nebula.graphd.env`                 | Graphd container environment variables                                                    | `[]`                                                                                                  |
| `nebula.graphd.resources`           | Graphd resources                                                                          | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebula.graphd.logVolume`           | Graphd log volume                                                                         | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebula.graphd.podLabels`           | Graphd pod labels                                                                         | `{}`                                                                                                  |
| `nebula.graphd.podAnnotations`      | Graphd pod annotations                                                                    | `{}`                                                                                                  |
| `nebula.graphd.SecurityContext`     | Graphd security context                                                                   | `{}`                                                                                                  |
| `nebula.graphd.nodeSelector`        | Graphd nodeSelector                                                                       | `{}`                                                                                                  |
| `nebula.graphd.tolerations`         | Graphd pod tolerations                                                                    | `[]`                                                                                                  |
| `nebula.graphd.affinity`            | Graphd affinity                                                                           | `{}`                                                                                                  |
| `nebula.graphd.readinessProbe`      | Graphd pod readinessProbe                                                                 | `{}`                                                                                                  |
| `nebula.graphd.livenessProbe`       | Graphd pod livenessProbe                                                                  | `{}`                                                                                                  |
| `nebula.graphd.initContainers`      | Graphd pod init containers                                                                | `[]`                                                                                                  |
| `nebula.graphd.sidecarContainers`   | Graphd pod sidecar containers                                                             | `[]`                                                                                                  |
| `nebula.graphd.volumes`             | Graphd pod volumes                                                                        | `[]`                                                                                                  |
| `nebula.graphd.volumeMounts`        | Graphd pod volume mounts                                                                  | `[]`                                                                                                  |
| `nebula.storaged.image`             | Storaged container image without tag, and use `nebula.version` as tag                     | `vesoft/nebula-storaged`                                                                              |
| `nebula.storaged.replicas`          | Storaged replica number                                                                   | `3`                                                                                                   |
| `nebula.storaged.env`               | Storaged container environment variables                                                  | `[]`                                                                                                  |
| `nebula.storaged.resources`         | Storaged resources                                                                        | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebula.storaged.logVolume`         | Storaged log volume                                                                       | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebula.storaged.dataVolumes`       | Storaged data volumes                                                                     | `[]`                                                                                                  |
| `nebula.storaged.podLabels`         | Storaged pod labels                                                                       | `{}`                                                                                                  |
| `nebula.storaged.podAnnotations`    | Storaged pod annotations                                                                  | `{}`                                                                                                  |
| `nebula.storaged.SecurityContext`   | Storaged security context                                                                 | `{}`                                                                                                  |
| `nebula.storaged.nodeSelector`      | Storaged nodeSelector                                                                     | `{}`                                                                                                  |
| `nebula.storaged.tolerations`       | Storaged pod tolerations                                                                  | `[]`                                                                                                  |
| `nebula.storaged.affinity`          | Storaged affinity                                                                         | `{}`                                                                                                  |
| `nebula.storaged.readinessProbe`    | Storaged pod readinessProbe                                                               | `{}`                                                                                                  |
| `nebula.storaged.livenessProbe`     | Storaged pod livenessProbe                                                                | `{}`                                                                                                  |
| `nebula.storaged.initContainers`    | Stroaged pod init containers                                                              | `[]`                                                                                                  |
| `nebula.storaged.sidecarContainers` | Storaged pod sidecar containers                                                           | `[]`                                                                                                  |
| `nebula.storaged.volumes`           | Storaged pod volumes                                                                      | `[]`                                                                                                  |
| `nebula.storaged.volumeMounts`      | Storaged pod volume mounts                                                                | `[]`                                                                                                  |
| `nebula.console.image`              | nebula console container image without tag                                                | `vesoft/nebula-console`                                                                               |
| `nebula.console.version`            | nebula console container image tag                                                        | `latest`                                                                                              |                                                                                          |                                                                                                  |
