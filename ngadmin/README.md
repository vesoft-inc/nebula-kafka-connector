# ngadmin

## build

```bash
make
```

# make ca cert

be sure to put the cert in the certs directory in the bin directory
```bash
cd bin 
mkdir cert
cd cert
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