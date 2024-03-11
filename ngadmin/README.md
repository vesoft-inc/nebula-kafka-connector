# ngadmin

## linux software dependencies
+ timeout

## build

```bash
make
```

# test
## make ca cert

be sure to put the cert in the certs directory in the bin directory
```bash
mkdir certs
```

## server cert
```bash
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
openssl req -new -x509 -key server.key -out ca.crt -days 365
```

## client cert
```bash
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey server.key -set_serial 01 -out client.crt
```

## start agent
```bash
cd agent/api/agent && go run agent.go
```

## install 
```bash
make
./bin/ngadmin -f ./examples/nebula.yaml 
```

## run ngadmin
```bash
make test-install
make test-uninstall
```

## registe product

edit the file `ngadmin/yamlparser/config.go` and add the product like this:
```go
"license-manager": {
		Name:          "license-manager",                                          // the name of the process for install directory & product name
		ExecShellPath: "./nebula-license-manager/scripts/license-manager.service", // for shell start
		ExecStartPath: "./nebula-license-manager/nebula-license-manager",          // for systemd start
		WorkingDir:    "./nebula-license-manager/",                                // for systemd start
		ConfigPath:    "./nebula-license-manager/etc/nebula-license-manager.yaml", // for merge user config to your product config
	},
```