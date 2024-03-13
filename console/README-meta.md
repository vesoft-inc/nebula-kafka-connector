# Nebula Graph meta-ctl

This directory contains meta-ctl, a tool to connect meta and control meta management.

## Meta command

Use `meta-ctl [command] [flag]` to execute meta command.

```bash
Execute meta command in cli mode. Use 'meta-ctl -h' to see usage.

Usage:
  meta-ctl [flags]
  meta-ctl [command]

Available Commands:
  cluster     Process cluster command
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Login meta server.
  service     Process service command
```

## 连接Meta

使用 meta-ctl 向 metad 集群发送命令需要先使用 login 命令连接到 meta(当前主要记录一个ip地址和leader信息，会缓存在~/nebula_meta_session文件中，后面只有 God 用户可以登陆)

例：连接到指定meta集群

```bash
meta-ctl login -H 192.168.8.6 -P 19565
```

```bash
Usage:
  meta-ctl login [flags]

Flags:
  -h, --help              help for login
  -H, --host string       meta server host (default "127.0.0.1")
  -p, --password string   password
  -P, --port uint32       meta server port (default 9559)
  -u, --user string       user name
```

## 集群操作

用以操作多个存储计算集群的相关命令,主要包含：

```bash
meta-ctl cluster -h

Execute cluster command in cli mode.

Usage:
  meta-ctl cluster [flags]
  meta-ctl cluster [command]

Available Commands:
  create      Create cluster in meta server.
  init        Init cluster storage part.
  show        Show cluster, show all if no cluster name specified.
```

### 创建集群

创建一个集群的常用命令如下：

```bash
meta-ctl cluster create -c=<clusterName> --replica=3 --zones="z1,z2,z3"
```

创建集群需要指定：

- 集群名称
- 集群的副本特性
- 集群的zone列表，如果不指定，则会生成一个默认zone.

```bash
Usage:
  meta-ctl cluster create [flags]

Flags:
  -h, --help                 help for create
  -r, --replica-factor int   replica number, default: 3 (default 3)
  -z, --zones stringArray    zones

Global Flags:
  -c, --cluster string   Cluster name
```

### 显示集群

```bash
meta-ctl cluster show  -h                                                                                                                                                     fix_nightly_workflow [7c66e60] (!) modified untracked
nebula-meta cluster show --cluster [clustername]

Usage:
  meta-ctl cluster show [flags]
```

示例

```bash
# 显示所有集群
meta-ctl cluster show
# 只显示一个集群
meta-ctl cluster show -c test_cluster
```

### 初始化集群

添加 storaged 后，需要初始化集群，主要是初始化 partition 分布

```bash
meta-console cluster init --cluster [clustername]

Usage:
  meta-ctl cluster init [flags]

Flags:
  -h, --help   help for init

Global Flags:
  -c, --cluster string   Cluster name
```

```bash
meta-ctl cluster init -c testcluster
```

### 删除集群

（待实现）

*TODO*

## 服务

执行服务管理的相关操作。

```bash
meta-ctl service -h

Execute service command in cli mode.

Usage:
  meta-ctl service [flags]
  meta-ctl service [command]

Available Commands:
  add         Add service into assigned cluster.
  drop        Drop service from assigned cluster.
  show        Show service in cluster.

Flags:
  -c, --cluster string   Cluster name
  -h, --help             help for service
```

### 增加服务

向集群中添加一条服务的命令如下：

```bash
meta-ctl service add --type [graphd|storaged] --host [host] --port [port] --cluster [clustername]

Usage:
  meta-ctl service add [flags]

Flags:
  -c, --cluster string   cluster name
  -h, --help             help for add
  -H, --host string      service host
  -P, --port uint32      service port
  -t, --type string      service type
```

示例

```bash
meta-ctl service add --type graphd -H 192.168.8.6 -P 9669 -c test_cluster
```

### 删除服务

（内核暂时只有接口定义，未实现逻辑）

```bash
meta-ctl service drop --type [graphd|storaged] --host [host] --port [port] --cluster [clustername]

Usage:
  meta-ctl service drop [flags]

Flags:
  -c, --cluster string   cluster name
  -h, --help             help for drop
  -H, --host string      service host
  -P, --port uint32      service port
  -t, --type string      service type
```

```bash
meta-ctl service drop --type graphd -H 192.168.8.6 -P 9669 -c testcluster
```

### 显示服务

```bash
meta-console service show --cluster [clustername]

Usage:
  meta-ctl service show [flags]

Flags:
  -h, --help   help for show

Global Flags:
  -c, --cluster string   Cluster name
```

```bash
meta-ctl service show -c testcluster
```
