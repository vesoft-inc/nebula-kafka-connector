package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Fatalf("unable to load client cert and key: %v", err)
	}

	caCert, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Fatalf("unable to read ca cert: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("unable to append ca cert to pool")
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://localhost:8888/health")
	if err != nil {
		log.Fatalf("unable to send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}
