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

package v2alpha1

const (
	GraphdConfigTemplate = `
########## basics ##########
# Whether to run as a daemon process
--daemonize=true
# The file to host the process id
--pid_file=pids/nebula-graphd.pid
# The file to host clusterid and service id.
--sid_file=sids/nebula-graphd.sid
# Whether to enable optimizer
--enable_optimizer=true
# The default charset when a space is created
--default_charset=utf8
# The default collate when a space is created
--default_collate=utf8_bin
# Whether to use the configuration obtained from the configuration file
--local_config=true

########## plugins ##########
--plugins=dbms.so,file_audit.so,logrotate.so

########## audit ##########
# This variable is used to enable audit. The value can be 'true' or 'false'.
--enable_audit=true

# This variable is used to specify the filename that’s going to store the audit log.
# It can contain the path relative to the install dir or absolute path.
# This variable has effect only when audit_log_handler is set to 'file'.
--audit_log_dir=./logs/audit/

# This variable is used to specify the audit log strategy, Optional：[ asynchronous｜ synchronous ]
# asynchronous: log using memory buffer, do not block the main thread
# synchronous: log directly to file, flush and sync every event
# Caution: For performance reasons, when the buffer is full and has not been flushed to the disk,
# the 'asynchronous' mode will discard subsequent requests.
--audit_log_strategy=synchronous

# This variable can be used to specify the size of memory buffer used for logging,
# used when audit_log_strategy variable is set to 'asynchronous' values.
# This variable has effect only when audit_log_handler is set to 'file'. Uint: B
--audit_log_max_buffer_size=1048576

# Event categories to be audited
--audit_log_categories=LOGIN,SIGNOUT,DDL,DCL,PROCEDURE,JOB,CONFIGURATION,RETURN_ERROR

########## logging ##########
# The directory to host logging files
--log_dir=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0
# Verbose log level, 1, 2, 3, 4, the higher of the level, the more verbose of the logging
--v=0
# Maximum seconds to buffer the log messages
--logbufsecs=0
# Whether to redirect stdout and stderr to separate output files
--redirect_stdout=true
# Destination filename of stdout and stderr, which will also reside in log_dir.
--stdout_log_file=graphd-stdout.log
--stderr_log_file=graphd-stderr.log
# Copy log messages at or above this level to stderr in addition to logfiles. The numbers of severity levels INFO, WARNING, ERROR, and FATAL are 0, 1, 2, and 3, respectively.
--stderrthreshold=2

########## log rotator ##########
# Maximum dbms log file size in MiB, default to 1024 MiB
--max_log_size=1024
# Info log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--info_log_rotate_size=-1
# Info log rotate days, log file older then this will be removed, default -1 (no limit)
--info_log_rotate_days=-1

# Maximum audit log file size in MiB, default to 1024 MiB
--audit_log_max_size=1024
# Audit log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--audit_log_rotate_size=-1
# Audit log rotate days, log file older then this will be removed, default -1 (no limit)
--audit_log_rotate_days=-1

########## query ##########
# Whether to treat partial success as an error.
# This flag is only used for Read-only access, and Modify access always treats partial success as an error.
--accept_partial_success=false
# Maximum sentence length, unit byte
--max_allowed_query_size=4194304

########## networking ##########
# Comma separated Meta Server Addresses
--meta_server_addrs=192.168.8.176:9559
# Local IP used to identify the nebula-graphd process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=192.168.8.176
# Network device to listen on
--listen_netdev=any
# Port to listen on
--port=9669
# To turn on SO_REUSEPORT or not
--reuse_port=false
# Backlog of the listen socket, adjust this together with net.core.somaxconn
--listen_backlog=1024
# The number of seconds Nebula service waits before closing the idle connections
--client_idle_timeout_secs=28800
# The number of seconds before idle sessions expire
# The range should be in [1, 604800]
--session_idle_timeout_secs=28800
# The number of threads to accept incoming connections
--num_accept_threads=1
# The number of networking IO threads, 0 for # of CPU cores
--num_netio_threads=0
# The number of threads to execute user queries, 0 for # of CPU cores
--num_worker_threads=0
# HTTP service ip
--ws_ip=0.0.0.0
# HTTP service port
--ws_http_port=19669
# storage client timeout
--storage_client_timeout_ms=60000
# Port to listen on Meta with HTTP protocol, it corresponds to ws_http_port in metad's configuration file
--ws_meta_http_port=19559

########## authentication ##########
# Enable authorization
--enable_authorize=false
# User login authentication type, password for nebula authentication, ldap for ldap authentication, cloud for cloud authentication
--auth_type=password

########## memory ##########
# process max memory in MiB, negative value means unlimited, default: -1
--process_max_memory_mib=-1
# enable print memory stats into log in check_memory_interval_in_secs interval, default false
--memory_stats_log=false
# background checking memory interval, changes made to graphd_max_memory_mib will take effect in this interval, default 1s
--check_memory_interval_in_secs=1


########## metrics ##########
--enable_space_level_metrics=false
`

	MetadhConfigTemplate = `
########## basics ##########

# Whether to run as a daemon process
--daemonize=true

# The file to host the process id
--pid_file=pids/nebula-metad.pid

########## plugins ##########
--plugins=logrotate.so

########## logging ##########
# The directory to host logging files
--log_dir=logs

########## log rotator ##########
# Maximum dbms log file size in MiB, default to 1024 MiB
--max_log_size=1024
# Info log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--info_log_rotate_size=-1
# Info log rotate days, log file older then this will be removed, default -1 (no limit)
--info_log_rotate_days=-1

# Maximum audit log file size in MiB, default to 1024 MiB
--audit_log_max_size=1024
# Audit log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--audit_log_rotate_size=-1
# Audit log rotate days, log file older then this will be removed, default -1 (no limit)
--audit_log_rotate_days=-1

# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0

# Verbose log level, 1, 2, 3, 4, the higher of the level, the more verbose of the logging
--v=0

# Maximum seconds to buffer the log messages
--logbufsecs=0

# Whether to redirect stdout and stderr to separate output files
--redirect_stdout=true

# Destination filename of stdout and stderr, which will also reside in log_dir.
--stdout_log_file=metad-stdout.log
--stderr_log_file=metad-stderr.log

# Copy log messages at or above this level to stderr in addition to logfiles. The numbers of severity levels INFO, WARNING, ERROR, and FATAL are 0, 1, 2, and 3, respectively.
--stderrthreshold=2

########## meta service ##########
# Comma separated Meta Server addresses
--meta_server_addrs=192.168.8.176:9559

# Local IP used to identify the nebula-metad process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=192.168.8.176

# Meta daemon listening port
--port=9559
# (To be deleted) Http port to execute meta command.
--admin_port=19559

########## storage ##########
# Root data path, here should be only single path for metad
--data_path=data/meta
`

	StoragedConfigTemplate = `
########## basics ##########
# Whether to run as a daemon process
--daemonize=true
# The file to host the process id
--pid_file=pids/nebula-storaged.pid
# The file to host clusterid and service id.
--sid_file=sids/nebula-storaged.sid

########## plugins ##########
--plugins=file_audit.so,logrotate.so

########## logging ##########
# The directory to host logging files
--log_dir=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0
# Verbose log level, 1, 2, 3, 4, the higher of the level, the more verbose of the logging
--v=0
# Maximum seconds to buffer the log messages
--logbufsecs=0
# Whether to redirect stdout and stderr to separate output files
--redirect_stdout=true
# Destination filename of stdout and stderr, which will also reside in log_dir.
--stdout_log_file=storaged-stdout.log
--stderr_log_file=storaged-stderr.log
# Copy log messages at or above this level to stderr in addition to logfiles. The numbers of severity levels INFO, WARNING, ERROR, and FATAL are 0, 1, 2, and 3, respectively.
--stderrthreshold=2

########## log rotator ##########
# Maximum dbms log file size in MiB, default to 1024 MiB
--max_log_size=1024
# Info log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--info_log_rotate_size=-1
# Info log rotate days, log file older then this will be removed, default -1 (no limit)
--info_log_rotate_days=-1

# Maximum audit log file size in MiB, default to 1024 MiB
--audit_log_max_size=1024
# Audit log rotate size in MiB, old log file exceed this size will be removed, default -1 (no limit)
--audit_log_rotate_size=-1
# Audit log rotate days, log file older then this will be removed, default -1 (no limit)
--audit_log_rotate_days=-1

########## networking ##########
# Comma separated Meta server addresses
--meta_server_addrs=192.168.8.176:9559
# Local IP used to identify the nebula-storaged process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=192.168.8.176
# Storage daemon listening port
--port=9779
# HTTP service ip
--ws_ip=0.0.0.0
# HTTP service port
--ws_http_port=19779
# heartbeat with meta service
--heartbeat_interval_secs=10

########## Disk ##########
# Root data path. Split by comma. e.g. --data_path=/disk1/path1/,/disk2/path2/
# One path per Rocksdb instance.
--data_path=data/storage

########## memory ##########
# process max memory in MiB, negative value means unlimited, default: -1
--process_max_memory_mib=-1
# enable print memory stats into log in check_memory_interval_in_secs interval, default false
--memory_stats_log=false
# background checking memory interval, changes made to storaged_max_memory_mib will take effect in this interval, default 1s
--check_memory_interval_in_secs=1
`
)
