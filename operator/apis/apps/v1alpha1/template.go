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

package v1alpha1

const (
	GraphdConfigTemplate = `
########## basics ##########
# Whether to run as a daemon process
--daemonize=true
# The file to host the process id
--pid_file=pids/nebula-graphd.pid
# The file to host clusterid and service id.
--sid_file=sids/nebula-graphd.sid

########## plugins ##########
--plugins=plugin_manager.so,dbms.so,file_audit.so,log_rotate.so

########## audit ##########
# This variable is used to enable audit. The value can be 'true' or 'false'.
--enable_audit=true
# Audit event categories to be audited, e.g. LOGIN,SIGNOUT, supported categories includes:
# [LOGIN | SIGNOUT | AUTHENTICATION | AUTHORIZATION | DDL | DQL | DCL | DML | DML_INSERT | DML_SET | DML_REMOVE |  DML_DELETE | PROCEDURE | ERROR | JOB | CONFIGURATION | GENERAL ]
--audit_log_categories=LOGIN,SIGNOUT,AUTHENTICATION,AUTHORIZATION,DDL,DCL,PROCEDURE,JOB,CONFIGURATION,ERROR

########## file audit plugin ##########
# This variable is used to specify the directory that’s going to store the audit log.
# It can contain the path relative to the install dir or absolute path.
--audit_log_dir=./logs/audit/
# This variable is used to specify the audit log strategy, Optional：[ asynchronous｜ synchronous ]
# asynchronous: log using memory buffer, do not block the main thread
# synchronous: log directly to file, flush and sync every event
# Caution: For performance reasons, when the buffer is full and has not been flushed to the disk,
# the 'asynchronous' mode will discard subsequent requests.
--audit_log_strategy=synchronous
# This variable is used to specify the size of memory buffer used for audit log,
# used when audit_log_strategy variable is set to 'asynchronous' values, Uint: Byte
--audit_log_max_buffer_size=1048576

########## logging ##########
# The directory to host logging files
--log_dir=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0
# Verbose log level, the higher of the level, the more verbose of the logging
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
# Enable log compress or not, default false
--log_compress=false

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
# Comma separated Meta Server Addresses
--meta_server_addrs=127.0.0.1:9559
# Local IP used to identify the nebula-graphd process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=127.0.0.1
# Port to listen on
--port=9669
# The number of seconds before idle sessions expire
# The range should be in [1, 604800]
--session_idle_timeout_secs=28800
# HTTP service ip/hostname
--http_host=0.0.0.0
# HTTP service port
--http_port=19669

########## threads ##########
# The number of networking IO threads, 0 for # of CPU cores
--num_netio_threads=0
# The number of threads to execute user queries, 0 for # of CPU cores
--num_worker_threads=0
# The number of threads used for computing algorithms
--num_computing_threads=8

########## memory ##########
# process max memory in MiB, negative value means unlimited, default: -1
--process_max_memory_mib=-1
# enable print memory stats into log in check_memory_interval_in_secs interval, default false
--memory_stats_log=false
# background checking memory interval, changes made to graphd_max_memory_mib will take effect in this interval, default 1s
--check_memory_interval_in_secs=1

########## heartbeat ##########
# heartbeat report interval in seconds
--heartbeat_interval_secs=10
`

	MetadhConfigTemplate = `
########## basics ##########

# Whether to run as a daemon process
--daemonize=true

# The file to host the process id
--pid_file=pids/nebula-metad.pid

# License manager address
--license_manager_url=license.vesoft-inc.com:9119

########## plugins ##########
--plugins=file_audit.so,log_rotate.so,auth_password.so

########## audit ##########
# This variable is used to enable audit. The value can be 'true' or 'false'.
--enable_audit=true
# Audit event categories to be audited, e.g. LOGIN,SIGNOUT, supported categories includes:
# [LOGIN | SIGNOUT | AUTHENTICATION | AUTHORIZATION | DDL | DQL | DCL | DML | DML_INSERT | DML_SET | DML_REMOVE |  DML_DELETE | PROCEDURE | ERROR | JOB | CONFIGURATION | GENERAL ]
--audit_log_categories=LOGIN,SIGNOUT,AUTHENTICATION,AUTHORIZATION,DDL,DCL,PROCEDURE,JOB,CONFIGURATION,ERROR

########## file audit plugin ##########
# This variable is used to specify the directory that’s going to store the audit log.
# It can contain the path relative to the install dir or absolute path.
--audit_log_dir=./logs/audit/
# This variable is used to specify the audit log strategy, Optional：[ asynchronous｜ synchronous ]
# asynchronous: log using memory buffer, do not block the main thread
# synchronous: log directly to file, flush and sync every event
# Caution: For performance reasons, when the buffer is full and has not been flushed to the disk,
# the 'asynchronous' mode will discard subsequent requests.
--audit_log_strategy=synchronous
# This variable is used to specify the size of memory buffer used for audit log,
# used when audit_log_strategy variable is set to 'asynchronous' values, Uint: Byte
--audit_log_max_buffer_size=1048576

########## logging ##########
# The directory to host logging files
--log_dir=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0
# Verbose log level, the higher of the level, the more verbose of the logging
--v=3
# Maximum seconds to buffer the log messages
--logbufsecs=0
# Whether to redirect stdout and stderr to separate output files
--redirect_stdout=true
# Destination filename of stdout and stderr, which will also reside in log_dir.
--stdout_log_file=metad-stdout.log
--stderr_log_file=metad-stderr.log
# Copy log messages at or above this level to stderr in addition to logfiles. The numbers of severity levels INFO, WARNING, ERROR, and FATAL are 0, 1, 2, and 3, respectively.
--stderrthreshold=2

########## log rotator ##########
# Enable log compress or not, default false
--log_compress=false

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

########## meta service ##########
# Comma separated Meta Server addresses
--meta_server_addrs=127.0.0.1:9559

# Local IP used to identify the nebula-metad process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=127.0.0.1

# Meta daemon listening port
--port=9559
# HTTP service ip/hostname
--http_host=0.0.0.0
# HTTP service port
--http_port=19559

########## storage ##########
# Root data path, here should be only single path for metad
--data_path=data/meta

########## threads ##########
# The number of threads to execute CPU bound tasks, 0 for # of CPU cores
--num_worker_threads=32

########## heartbeat ##########
# heartbeat report interval in seconds
--heartbeat_interval_secs=10
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
--plugins=file_audit.so,log_rotate.so

########## audit ##########
# This variable is used to enable audit. The value can be 'true' or 'false'.
--enable_audit=true
# Audit event categories to be audited, e.g. LOGIN,SIGNOUT, supported categories includes:
# [LOGIN | SIGNOUT | AUTHENTICATION | AUTHORIZATION | DDL | DQL | DCL | DML | DML_INSERT | DML_SET | DML_REMOVE |  DML_DELETE | PROCEDURE | ERROR | JOB | CONFIGURATION | GENERAL ]
--audit_log_categories=LOGIN,SIGNOUT,AUTHENTICATION,AUTHORIZATION,DDL,DCL,PROCEDURE,JOB,CONFIGURATION,ERROR

########## file audit plugin ##########
# This variable is used to specify the directory that’s going to store the audit log.
# It can contain the path relative to the install dir or absolute path.
--audit_log_dir=./logs/audit/
# This variable is used to specify the audit log strategy, Optional：[ asynchronous｜ synchronous ]
# asynchronous: log using memory buffer, do not block the main thread
# synchronous: log directly to file, flush and sync every event
# Caution: For performance reasons, when the buffer is full and has not been flushed to the disk,
# the 'asynchronous' mode will discard subsequent requests.
--audit_log_strategy=synchronous
# This variable is used to specify the size of memory buffer used for audit log,
# used when audit_log_strategy variable is set to 'asynchronous' values, Uint: Byte
--audit_log_max_buffer_size=1048576

########## logging ##########
# The directory to host logging files
--log_dir=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--minloglevel=0
# Verbose log level, the higher of the level, the more verbose of the logging
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
# Enable log compress or not, default false
--log_compress=false

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
--meta_server_addrs=127.0.0.1:9559
# Local IP used to identify the nebula-storaged process.
# Change it to an address other than loopback if the service is distributed or
# will be accessed remotely.
--local_ip=127.0.0.1
# Storage daemon listening port
--port=9779
# HTTP service ip/hostname
--http_host=0.0.0.0
# HTTP service port
--http_port=19779

########## Disk ##########
# Root data path. Split by comma. e.g. --data_path=/disk1/path1/,/disk2/path2/
# One partition per Rocksdb instance.
--data_path=data/storage

########## memory ##########
# process max memory in MiB, negative value means unlimited, default: -1
--process_max_memory_mib=-1
# enable print memory stats into log in check_memory_interval_in_secs interval, default false
--memory_stats_log=false
# background checking memory interval, changes made to storaged_max_memory_mib will take effect in this interval, default 1s
--check_memory_interval_in_secs=1

########## threads ##########
# The number of threads to accomplish network IO
--num_io_threads=16
# The num of raft worker threads
--num_raft_worker_threads=32
# The num of threads to handle request and response
--num_storage_worker_threads=32
# The num of threads to run execution plan
--num_storage_query_threads=32

########## heartbeat ##########
# heartbeat report interval in seconds
--heartbeat_interval_secs=10
`
)
