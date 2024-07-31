## Installing Add-ons

**Caution:**
This section links to third party projects that provide functionality required by nebula-operator. The nebula-operator
project authors aren't responsible for these projects.

## coredns

[CoreDNS](https://coredns.io/) is a flexible, extensible DNS server which can
be [installed](https://github.com/coredns/deployment/tree/master/kubernetes) as the in-cluster DNS for pods.

Each component of NebulaGraph can communicate via DNS using an address like _x.default.svc.cluster.local_. Coredns is used for address
resolution.

## cert-manager

**Note:**
If you set helm chart nebula-operator _.Values.admissionWebhook.create_ to false, cert-manager is not needed.

[cert-manager](https://cert-manager.io/) is a tool that automates certificate management. It makes use of extending the
Kubernetes API server using a Webhook server to provide dynamic admission control over cert-manager resources.

Refer to the [cert-manager installation documentation](https://cert-manager.io/docs/installation/) to get
started.

Cert-manager is used for validating the replicas of each NebulaGraph component. If you run NebulaGraph in a production environment and
need high availability make sure to set _.Values.admissionWebhook.create_ to true and install cert-manager.
