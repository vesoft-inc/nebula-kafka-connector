## Install Guide

Instead of manually installing, scaling, upgrading, and uninstalling NebulaGraph in a production environment, you can
use the Nebula [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) to manage NebulaGraph automatically.

Follow this guide to install Nebula Operator using Helm for an in-depth evaluation.

### Requirements

* Kubernetes >= 1.24
* [RBAC](https://kubernetes.io/docs/admin/authorization/rbac) enabled (optional)
* [CoreDNS](https://github.com/coredns/coredns) >= 1.10.0
* [Helm](https://helm.sh) >= 3.2.0

## Add-ons

See [add-ons](add-ons.md) for how to install the add-ons.

### Get Repo Info

```shell script
$ helm repo add nebula-ng-operator https://vesoft-inc.github.io/nebula-ng-tools/operator/charts
$ helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for the command documentation._

### Install Operator

```shell script
# helm install [NAME] [CHART] [flags]
$ helm install nebula-ng-operator nebula-ng-operator/nebula-operator --namespace=nebula-operator-system --version=${chart_version}
```

Note:   
If the corresponding nebula-operator-system namespace does not exist, you can create the namespace first by running the _kubectl
create namespace nebula-operator-system_ command.

${chart_version} represents the chart version of Nebula Operator. For example, v1.0.0. You can view the currently
supported versions by running the _helm search repo -l nebula-ng-operator_ command.

_See [configuration](#configure-operator) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for the command documentation._

### Configure Operator

See [Customizing the Chart Before Installing](https://helm.sh/docs/intro/using_helm/#customizing-the-chart-before-installing)
. To see all configurable options with detailed comments, visit the
chart's [values.yaml](https://github.com/vesoft-inc/nebula-ng-tools/blob/master/operator/charts/nebula-operator/values.yaml), or
run the following commands:

```shell script
$ helm show values nebula-ng-operator/nebula-operator
```

### Upgrade Operator

If you need to upgrade the Nebula Operator, modify the ${HOME}/nebula-operator/values.yaml file, and then execute the
following command to upgrade:

```shell script
$ helm upgrade nebula-ng-operator nebula-ng-operator/nebula-operator --namespace=nebula-operator-system -f `${HOME}/nebula-operator/values.yaml`
```

### Uninstall Operator

```shell script
$ helm uninstall nebula-ng-operator --namespace=nebula-operator-system
$ kubectl delete crd nebulaclusters.apps.nebula-graph.io
```
