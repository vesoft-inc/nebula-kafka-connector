# Nebula Operator 5.x
Nebula Operator 5.x manages [NebulaGraph](https://github.com/vesoft-inc/nebula-ng) 5.x metads and clusters
on [Kubernetes](https://kubernetes.io) by automating tasks related to the operation of a NebulaGraph cluster.
It makes NebulaGraph a truly cloud-native database.

## Quick Start

- [Install Nebula Operator](#install-nebula-operator)
- [Create and Destroy](#create-and-destroy-a-nebula-cluster)
- [Create and destroy a nebula cluster using helm charts](#create-and-destroy-a-nebula-cluster-using-helm-charts)
- [Connecting to a nebula cluster](#connecting-to-a-nebula-cluster)
- [Resize](#resize-a-nebula-cluster)
- [Failover](#failover)


### Install nebula operator

See [install/uninstall nebula operator](doc/user/operator_guide.md) .

### Create and destroy a nebula cluster

#### Creating a nebula cluster

Create the nebula secret to reset the password for the root user. This is required in 5.x.
```bash
$ kubectl create secret generic nebula --from-literal=username=root --from-literal=password=<new_password_for_root_user>
```

Start the metad pods
```bash
$ kubectl create -f config/samples/nebulametad.yaml
```

Wait for the metad pods to finish starting up. Then start the graphd, storaged and nebula console pods.
```bash
$ kubectl create -f config/samples/nebulacluster.yaml
```

A none ha-mode nebula cluster will be created.

```bash
$ kubectl get pods -l app.kubernetes.io/cluster=nebula-metad
NAME                   READY   STATUS    RESTARTS   AGE
nebula-metad-metad-0   1/1     Running   0          4m11s
nebula-metad-metad-1   1/1     Running   0          4m11s
nebula-metad-metad-2   1/1     Running   0          4m11s
```

``` bash
$ kubectl get pods -l app.kubernetes.io/cluster=nebula-cluster
NAME                        READY   STATUS    RESTARTS   AGE
nebula-cluster-graphd-0     1/1     Running   0          1m
nebula-cluster-storaged-0   1/1     Running   0          1m
nebula-cluster-storaged-1   1/1     Running   0          1m
nebula-cluster-storaged-2   1/1     Running   0          1m
```

**Warning:**
In 5.x one metad is able to managed several different clusters. So the name specified under the metadRef field in nebulacluster.yaml 
must match either the name of the metad in nebulametad.yaml, or an exising metad in the namesapce nebulacluster.yaml is deployed in. 
Otherwise the nebula cluster will not deploy successfully.


#### Destroying a nebula cluster

Delete the graphd and storaged pods
```bash
$ kubectl delete -f config/samples/nebulacluster.yaml
```

Delete the metad pods
```bash
$ kubectl delete -f config/samples/nebulametad.yaml
```
**Warning:**
Only delete the metad pods with the above command when all nebula clusters it manages have been deleted.
Otherwise the metad pods will fail to delete.

Delete the nebula secret
```bash
$ kubectl delete secret nebula
```


### Create and destroy a nebula cluster using helm charts

See [Create/Destroy a Nebula Cluster with helm](doc/user/nebula_cluster_guide.md).


### Connecting to a nebula cluster
Nebula operator 5.x will automatically deploy a graphd service of type NodePort and a nebula console pod for each nebula cluster created,
which will help you interact with NebulaGraph using the natively supported GQL. This interaction can be both within and outside of the cluster.

To see the NodePort service run
```bash
$ kubectl get service -l app.kubernetes.io/cluster=nebula-cluster
NAME                               TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                          AGE
nebula-cluster-graphd-headless     ClusterIP   None          <none>        9669/TCP,19669/TCP               14h
nebula-cluster-graphd-svc          NodePort    10.96.84.74   <none>        9669:31291/TCP,19669:31869/TCP   14h
nebula-cluster-storaged-headless   ClusterIP   None          <none>        9779/TCP,19779/TCP               14h
```
The graphd service is the one that ends with **graphd-svc**.

To interact with the Nebula cluster run
```bash
$ kubectl exec -it nebula-cluster-console -n <nebula-cluster-namespace> -- ngql --host <host_ip> --port <service_port> --user <nebula_username> --password <nebula_password>
```
Replace the values in the angular brackets with the following:
| Value                               | Replace With                                                                                                                                            |
|:------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------|
| <nebula-cluster-namespace>          | The kubernetes namespace the nebula cluster was deployed in. Leave out along with the "-n" for the default namespace.                                   |
| <host_ip>                           | The cluster-ip of the graphd service if connecting from with in the cluster. Otherwise the ip of any node in the cluster.                               |
| <service_port>                      | The port of the graphd service that port 9669 is mapped to (**31291** in the case of this example).                                                     |
| <nebula_username>                   | The username used to connect to the Nebula Cluster (usually **root** for the initial connection).                                                       |
| <nebula_password>                   | The password used to connect to the Nebula Cluster (should be the one set in the nebula secret during cluster installation for the initial connection). |

You should see the following prompt if successful
```bash
$ kubectl exec -it nebula-cluster-console -n nebula -- ngql --host 172.18.0.4 --port 31291 --user root --password *****
Welcome to NebulaGraph 5.0, the distributed graph database offering native GQL support!
:help for help.

(root@nebula) [(none)]>
```


### Resize a nebula cluster

[Create a nebula cluster](#creating-a-nebula-cluster)

In `config/samples/nebulacluster.yaml` storaged replicas is 3 by default.  
Modify the file and change `replicas` from 3 to 5.

```yaml
  storaged:
    resources:
      requests:
        cpu: "500m"
        memory: "500Mi"
      limits:
        cpu: "2"
        memory: "2Gi"
    replicas: 5
    image: reg.vesoft-inc.com/vesoft-ng/nebula-storaged
    version: nightly
    dataVolumeClaims:
    - resources:
        requests:
          storage: 10Gi
      storageClassName: local-path
    logVolumeClaim:
      resources:
        requests:
          storage: 1Gi
      storageClassName: local-path
```

Apply the replicas change to the cluster CR:

```bash
$ kubectl apply -f config/samples/nebulacluster.yaml
```

The storaged cluster will scale to 5 members (5 pods):

```bash
$ kubectl get pods -l app.kubernetes.io/cluster=nebula-cluster
NAME                        READY   STATUS    RESTARTS   AGE
nebula-cluster-graphd-0     1/1     Running   0          63m
nebula-cluster-storaged-0   1/1     Running   0          63m
nebula-cluster-storaged-1   1/1     Running   0          63m
nebula-cluster-storaged-2   1/1     Running   0          63m
nebula-cluster-storaged-3   1/1     Running   0          66s
nebula-cluster-storaged-4   1/1     Running   0          66s
```

Similarly, you can decrease the size of the cluster from 5 back down to 3 by changing the replicas field again and reapplying
the change.

```yaml
  storaged:
    resources:
      requests:
        cpu: "500m"
        memory: "500Mi"
      limits:
        cpu: "2"
        memory: "2Gi"
    replicas: 3
    image: reg.vesoft-inc.com/vesoft-ng/nebula-storaged
    version: nightly
    dataVolumeClaims:
    - resources:
        requests:
          storage: 10Gi
      storageClassName: local-path
    logVolumeClaim:
      resources:
        requests:
          storage: 1Gi
      storageClassName: local-path
```

You can see that the storaged cluster will eventually be reduced to 3 pods:

```bash
$ kubectl get pods -l app.kubernetes.io/cluster=nebula-cluster
NAME                        READY   STATUS    RESTARTS   AGE
nebula-cluster-graphd-0     1/1     Running   0          72m
nebula-cluster-storaged-0   1/1     Running   0          72m
nebula-cluster-storaged-1   1/1     Running   0          72m
nebula-cluster-storaged-2   1/1     Running   0          72m
```


### Failover

Nebula operator will automatically heal the cluster if some of the components fail. The following outlines how this is done.

[Create a nebula cluster](#creating-a-nebula-cluster)

Wait until pods are up. Simulate a member failure by deleting a storaged pod:

```bash
$ kubectl delete pod nebula-storaged-2 --now
```

The nebula operator will recover the failure by creating a new pod `nebula-storaged-2`:

```bash
$ kubectl get pods -l app.kubernetes.io/cluster=nebula-cluster
NAME                        READY   STATUS    RESTARTS   AGE
nebula-cluster-graphd-0     1/1     Running   0          3h50m
nebula-cluster-storaged-0   1/1     Running   0          3h50m
nebula-cluster-storaged-1   1/1     Running   0          3h50m
nebula-cluster-storaged-2   1/1     Running   0          31s
```


## FAQ

Please refer to [FAQ.md](FAQ.md)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

