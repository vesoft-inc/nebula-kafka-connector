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
| Parameter                                  | Description                                                                               | Default                                                                                               |
|:-------------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nameOverride`                             | Override the name of the chart                                                            | `nil`                                                                                                 |

**Nebula Metad Chart Only**
| Parameter                                  | Description                                                                               | Default                                                                                               |
|:-------------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nebulaMetad.affinity`                     | Metad affinity                                                                            | `{}`                                                                                                  |
| `nebulaMetad.credentialSecret`             | The secret that contains the credentials for the nebula graph root user                   | `"nebula"`                                                                                            |
| `nebulaMetad.dataVolume`                   | Metad data volume                                                                         | `[]`                                                                                                  |
| `nebulaMetad.enablePVReclaim`              | Flag to enable/disable PV reclaim while the Nebula cluster deleted                        | `false`                                                                                               |
| `nebulaMetad.env`                          | Metad container environment variables                                                     | `[]`                                                                                                  |
| `nebulaMetad.image`                        | Metad container image without tag, and use `nebula.version` as tag                        | `vesoft/nebula-metad`                                                                                 |
| `nebulaMetad.imagePullPolicy`              | Nebula image pull policy                                                                  | `Always`                                                                                              |
| `nebulaMetad.imagePullSecrets`             | The secret that contains the credentials for pulling the images                           | `[]`                                                                                                  |
| `nebulaMetad.initContainers`               | Metad pod init containers                                                                 | `[]`                                                                                                  |
| `nebulaMetad.livenessProbe`                | Metad pod livenessProbe                                                                   | `{}`                                                                                                  |
| `nebulaMetad.logVolume`                    | Metad log volume                                                                          | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebulaMetad.nodeSelector`                 | Metad nodeSelector                                                                        | `{}`                                                                                                  |
| `nebulaMetad.podAnnotations`               | Metad pod annotations                                                                     | `{}`                                                                                                  |
| `nebulaMetad.podLabels`                    | Metad pod labels                                                                          | `{}`                                                                                                  |
| `nebulaMetad.readinessProbe`               | Metad pod readinessProbe                                                                  | `{}`                                                                                                  |
| `nebulaMetad.replicas`                     | Metad replica number                                                                      | `3`                                                                                                   |
| `nebulaMetad.resources`                    | Metad resources                                                                           | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebulaMetad.schedulerName`                | Scheduler for nebula component                                                            | `default-scheduler`                                                                                   |
| `nebulaMetad.sidecarContainers`            | Metad pod sidecar containers                                                              | `[]`                                                                                                  |
| `nebulaMetad.storageClassName`             | StorageClass object name                                                                  | `""`                                                                                                  |
| `nebulaMetad.tolerations`                  | Metad pod tolerations                                                                     | `[]`                                                                                                  |
| `nebulaMetad.topologySpreadConstraints`    | Topology spread constraints to control how Pods are spread across failure-domains         | `[]`                                                                                                  |
| `nebulaMetad.version`                      | Nebula image tag                                                                          | `v3.6.0`                                                                                              |
| `nebulaMetad.volumeMounts`                 | Metad pod volume mounts                                                                   | `[]`                                                                                                  |
| `nebulaMetad.volumes`                      | Metad pod volumes                                                                         | `[]`                                                                                                  |


**Nebula Cluster Chart Only**
| Parameter                                  | Description                                                                               | Default                                                                                               |
|:-------------------------------------------|:------------------------------------------------------------------------------------------|:------------------------------------------------------------------------------------------------------|
| `nebulaCluster.console.image`              | nebula console container image without tag                                                | `vesoft/nebula-console`                                                                               |
| `nebulaCluster.console.version`            | nebula console container image tag                                                        | `latest`                                                                                              |
| `nebulaCluster.credentialSecret`           | The secret that contains the credentials for the nebula graph root user                   | `"nebula"`                                                                                            |
| `nebulaCluster.enablePVReclaim`            | Flag to enable/disable PV reclaim while the Nebula cluster deleted                        | `false`                                                                                               |
| `nebulaCluster.graphd.affinity`            | Graphd affinity                                                                           | `{}`                                                                                                  |
| `nebulaCluster.graphd.env`                 | Graphd container environment variables                                                    | `[]`                                                                                                  |
| `nebulaCluster.graphd.image`               | Graphd container image without tag                                                        | `vesoft/nebula-graphd`                                                                                |
| `nebulaCluster.graphd.initContainers`      | Graphd pod init containers                                                                | `[]`                                                                                                  |
| `nebulaCluster.graphd.livenessProbe`       | Graphd pod livenessProbe                                                                  | `{}`                                                                                                  |
| `nebulaCluster.graphd.logVolume`           | Graphd log volume                                                                         | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebulaCluster.graphd.nodeSelector`        | Graphd nodeSelector                                                                       | `{}`                                                                                                  |
| `nebulaCluster.graphd.podAnnotations`      | Graphd pod annotations                                                                    | `{}`                                                                                                  |
| `nebulaCluster.graphd.podLabels`           | Graphd pod labels                                                                         | `{}`                                                                                                  |
| `nebulaCluster.graphd.readinessProbe`      | Graphd pod readinessProbe                                                                 | `{}`                                                                                                  |
| `nebulaCluster.graphd.replicas`            | Graphd replica number                                                                     | `2`                                                                                                   |
| `nebulaCluster.graphd.resources`           | Graphd resources                                                                          | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebulaCluster.graphd.serviceType`         | Type of kubernetes service to create for graphd                                           | `"NodePort"`                                                                                          |
| `nebulaCluster.graphd.sidecarContainers`   | Graphd pod sidecar containers                                                             | `[]`                                                                                                  |
| `nebulaCluster.graphd.storageClassName`    | Graphd StorageClass object name                                                           | `""`                                                                                                  |
| `nebulaCluster.graphd.tolerations`         | Graphd pod tolerations                                                                    | `[]`                                                                                                  |
| `nebulaCluster.graphd.version`             | Nebula Graphd image tag                                                                   | `v3.6.0`                                                                                              |
| `nebulaCluster.graphd.volumeMounts`        | Graphd pod volume mounts                                                                  | `[]`                                                                                                  |
| `nebulaCluster.graphd.volumes`             | Graphd pod volumes                                                                        | `[]`                                                                                                  |
| `nebulaCluster.imagePullPolicy`            | Nebula Cluster image pull policy                                                          | `Always`                                                                                              |
| `nebulaCluster.imagePullSecrets`           | The secret that contains the credentials for pulling the images                           | `[]`                                                                                                  |
| `nebulaCluster.metadRef`                   | Reference to the metad that manages the nebula cluster                                    | `{"name": "metad-sample"}`                                                                            |
| `nebulaCluster.replicaFactor`              | The number or replicas for each partition                                                 | `3`                                                                                                   |
| `nebulaCluster.schedulerName`              | Scheduler for nebula component                                                            | `default-scheduler`                                                                                   |
| `nebulaCluster.storaged.affinity`          | Storaged affinity                                                                         | `{}`                                                                                                  |
| `nebulaCluster.storaged.dataVolumes`       | Storaged data volumes                                                                     | `[]`                                                                                                  |
| `nebulaCluster.storaged.enableAutoBalance` | Enable auto load balance for storaged                                                     | `false`                                                                                               |
| `nebulaCluster.storaged.env`               | Storaged container environment variables                                                  | `[]`                                                                                                  |
| `nebulaCluster.storaged.image`             | Storaged container image without tag                                                      | `vesoft/nebula-storaged`                                                                              |
| `nebulaCluster.storaged.initContainers`    | Stroaged pod init containers                                                              | `[]`                                                                                                  |
| `nebulaCluster.storaged.livenessProbe`     | Storaged pod livenessProbe                                                                | `{}`                                                                                                  |
| `nebulaCluster.storaged.logVolume`         | Storaged log volume                                                                       | `{"enable":true,"storage":"500Mi"}`                                                                   |
| `nebulaCluster.storaged.nodeSelector`      | Storaged nodeSelector                                                                     | `{}`                                                                                                  |
| `nebulaCluster.storaged.podAnnotations`    | Storaged pod annotations                                                                  | `{}`                                                                                                  |
| `nebulaCluster.storaged.podLabels`         | Storaged pod labels                                                                       | `{}`                                                                                                  |
| `nebulaCluster.storaged.readinessProbe`    | Storaged pod readinessProbe                                                               | `{}`                                                                                                  |
| `nebulaCluster.storaged.replicas`          | Storaged replica number                                                                   | `3`                                                                                                   |
| `nebulaCluster.storaged.resources`         | Storaged resources                                                                        | `{"resources":{"requests":{"cpu":"500m","memory":"500Mi"},"limits":{"cpu":"1","memory":"1Gi"}}}`      |
| `nebulaCluster.storaged.sidecarContainers` | Storaged pod sidecar containers                                                           | `[]`                                                                                                  |
| `nebulaCluster.storaged.storageClassName`  | StorageClass object name                                                                  | `""`                                                                                                  |
| `nebulaCluster.storaged.tolerations`       | Storaged pod tolerations                                                                  | `[]`                                                                                                  |
| `nebulaCluster.storaged.version`           | Nebula image tag for storaged                                                             | `v3.6.0`                                                                                              |
| `nebulaCluster.storaged.volumeMounts`      | Storaged pod volume mounts                                                                | `[]`                                                                                                  |
| `nebulaCluster.storaged.volumes`           | Storaged pod volumes                                                                      | `[]`                                                                                                  |
| `nebulaCluster.topologySpreadConstraints`  | Topology spread constraints to control how Pods are spread across failure-domains         | `[topologyKey: "kubernetes.io/hostname", whenUnsatisfiable: "ScheduleAnyway"]`                        |
| `nebulaCluster.zones`                      | Names of the zones to deploy the replicas into                                            | `[az1, az2, az3]`                                                                                     |
