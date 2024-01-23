# Nebula Graph meta-console

This directory contains meta-console, a tool to connect meta and do cluster management.


## Meta command

Use `meta-console [command] [flag]` to execute meta command.


```bash
Use meta-console to control meta.

Usage:
  meta-console [flags]
  meta-console [command]

Available Commands:
  cluster          Process cluster command
  completion       Generate the autocompletion script for the specified shell
  connect          Execute meta command in cli mode.
  help             Help about any command
  service          Process service command
```

## 连接Meta

使用使用meta-console向meta集群发送命令需要先使用connect命令连接到meta(当前主要记录一个ip地址和leader信息，会缓存在~/nebula_meta_session文件中)

例：连接到指定meta集群

```bash
meta-console connect -m "meta1:9559,meta2:9559,meta3:9559"
```

```bash
Usage:
  meta-console connect [flags]

Flags:
  -h, --help           help for connect
  -m, --metas string   meta server address list separated by comma like "xx:xx,xx:xx"
```

## 集群操作

用以操作多个存储计算集群的相关命令,主要包含：

```
Available Commands:
  create      Create cluster in server.
  init        Init cluster storage part.
  show        Show cluster, show all if no cluster name specified.
```

### 创建集群

创建一个集群的常用命令如下：

`meta-console cluster create --name=<name> --replica=3 --zones="z1,z2,z3"`

创建集群需要指定：

- 集群名称
- 集群的副本特性
- 集群的zone列表，如果不指定，则会生成一个默认zone.


```
Usage:
  meta-console cluster create [flags]
Flags:
  -h, --help            help for create
  -n, --name string     cluster name
  -r, --replica int     replica number
  -z, --zones strings   zone list
```

### 删除集群

*TODO*

## 服务

执行服务管理的相关操作。

```bash
Available Commands:
  add         Add service into assigned cluster.
  show        Show service in cluster.
```

### 增加服务

向集群中添加一条服务的命令如下：

```bash
Usage:
  meta-console service add [flags]

Flags:
  -c, --cluster string   cluster name
  -h, --help             help for add
  -i, --ip string        service ip
  -p, --port uint32      service port
  -t, --type string      service type
```

### 删除服务

*TODO*

### 替换服务

*TODO*


## 触发任务

用以向meta触发对应graph的统计任务。

```bash
nebula-console meta triggerstatstask --graph [graphname]

Usage:
  meta-console triggerstatstask [flags]

Flags:
  -g, --graph string   graph name
  -h, --help           help for triggerstatstask
```