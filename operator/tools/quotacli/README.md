## quota-cli

```shell
A tool for k8s user and quota management. Use 'quota-cli -h' to see usage.

Usage:
  quota-cli
  quota-cli [command]

Available Commands:
  create      create user and resource quota for nebula graph
  help        Help about any command
  list        list all user and resource quotas used for nebula graph

Flags:
  -c, --config-path string   Path to the kubeconfig file to use for CLI requests.
  -h, --help                 help for quota-cli

Use "quota-cli [command] --help" for more information about a command.
```

### create
```shell
$ kubectl config current-context 
kubernetes-admin@kubernetes

$ quotacli create --quota-user john --quota-namespace john-ns --cluster-name kubernetes --resource-requests cpu=3,memory=6Gi --resource-limits cpu=4,memory=8Gi

$ ls -lt certs/
total 8
-rw-r--r-- 1 root root  883 Oct 12 18:24 john.csr
-rw-r--r-- 1 root root 1675 Oct 12 18:24 john.key

$ ls -lt kube/
total 8
-rw-r--r-- 1 root root 5593 Oct 12 18:24 john-kubeconfig

$ kubectl --kubeconfig kube/john-kubeconfig get resourcequotas 
NAME               AGE    REQUEST                                     LIMIT
compute-resource   2d7h   requests.cpu: 0/3, requests.memory: 0/6Gi   limits.cpu: 0/4, limits.memory: 0/8Gi
```

### list

```shell
$ quotacli list
User:            iris
Name:            compute-resource
Namespace:       iris-ns
Resource         Used  Hard
--------         ----  ----
limits.cpu       0     4
limits.memory    0     8Gi
requests.cpu     0     3
requests.memory  0     6Gi


User:            john
Name:            compute-resource
Namespace:       john-ns
Resource         Used  Hard
--------         ----  ----
limits.cpu       0     4
limits.memory    0     8Gi
requests.cpu     0     3
requests.memory  0     6Gi
```