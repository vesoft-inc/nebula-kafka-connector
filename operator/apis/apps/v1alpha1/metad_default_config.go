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

const MetadDefaultConfig = `
########## basics ##########
# Whether to run as a daemon process
--daemonize=true

# The file to host the process id
--pid_file=pids/nebula-metad.pid

# License manager address
--license_manager_url=license.vesoft-inc.com:9119

########## plugins ##########
--metad_builtin_plugins=auth_password,file_audit,log_rotate

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
--log_directory=logs
# Log level, 0, 1, 2, 3 for INFO, WARNING, ERROR, FATAL respectively
--min_log_level=0
# Verbose log level, the higher of the level, the more verbose of the logging
--v_log_level=0
# Maximum seconds to buffer the log messages
--log_buf_secs=0
# Whether to redirect stdout and stderr to separate output files
--redirect_stdout=true
# Destination filename of stdout and stderr, which will also reside in log_dir.
--stdout_log_file=metad-stdout.log
--stderr_log_file=metad-stderr.log

########## log rotator ##########
# Enable log compress or not, effect both db log and audit log, default false
--log_compress=false

##### db log rotator #####
# Maximum single log file size in MiB, default to 1024 MiB
--max_log_size=1024

# Info log rotate total size in MiB, 
# Only keep the most recent logs whose total size less then this configuration,
# Older log files will be removed, default to 10240 MiB
--info_log_rotate_size=10240

# Info log rotate days,
# Log files older then this will be removed, default to 7 days
--info_log_rotate_days=7

##### audit log rotator #####
# Maximum audit single log file size in MiB, default to 1024 MiB
--audit_log_max_size=1024

# Audit log rotate total size in MiB, 
# Only keep the most recent logs whose total size less then this configuration,
# Older log files will be removed, default -1 (no limit)
--audit_log_rotate_size=-1

# Audit log rotate days,
# Log file older then this will be removed, default -1 (no limit)
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

# TLS
# Enable TLS
--tls_enable=false
# Enable mTLS
--tls_client_verify=false
# Server certificate
--tls_cert=etc/tls/current/certs/meta-server-cert.pem
# Server private key
--tls_key=etc/tls/current/private/meta-server-key.pem
# Password file for encrypted private key
--tls_key_passfile=
# CA certificate to authenticate the client certificate
# Only necessary if --tls_client_verify enabled
--tls_ca=etc/tls/current/certs/ca/nebula-ca-cert.pem

# Enable TLS for storage client
--storage_tls_enable=false
# Client certificate, only necessary if mTLS enabled in the server side
--storage_tls_cert=etc/tls/current/certs/storage-client-cert.pem
# Client private key, only necessary if mTLS enabled in the server side
--storage_tls_key=etc/tls/current/private/storage-client-key.pem
# Password file for encrypted private key
--storage_tls_key_passfile=
# CA certificate to authenticate the server certificate
--storage_tls_ca=etc/tls/current/certs/ca/nebula-ca-cert.pem
# Expected peer name to check against CN or SAN in the server certificate
--storage_tls_peer_name=storage.server.vesoft.com
# Enable the peer name checking
--storage_tls_peer_name_verify=true

# Enable HTTPS
--http_tls_enable=false
# Enable mTLS
--http_tls_client_verify=false
# Certificate for the HTTP service
--http_tls_cert=etc/tls/current/certs/http-server-cert.pem
# Private key for the HTTP service
--http_tls_key=etc/tls/current/private/http-server-key.pem
# Password file for encrypted private key
--http_tls_key_passfile=
# CA certificate to verify the client certificate
# Only necessary if mTLS enabled
--http_tls_ca=etc/tls/current/certs/ca/nebula-ca-cert.pem

# Interval in seconds to check for cert rotation
--tls_cert_check_interval=60
# Interval for cert rotation check of gRPC.
# This will be deprecated once the configuration module is ready
--grpc_tls_cert_check_interval=60

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
