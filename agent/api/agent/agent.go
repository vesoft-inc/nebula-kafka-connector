package main

import (
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"

	gopkgmiddleware "github.com/vesoft-inc/go-pkg/middleware"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/config"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/handler"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/svc"

	"github.com/NYTimes/gziphandler"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/agent-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	caCert, err := os.ReadFile(c.CAFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	svcCtx := svc.NewServiceContext(c)

	gzipHandler := gziphandler.MustNewGzipLevelHandler(gzip.DefaultCompression)
	server := rest.MustNewServer(c.RestConf,
		rest.WithTLSConfig(tlsConfig),
	)
	defer server.Stop()

	server.Use(rest.ToMiddleware(gzipHandler))
	server.Use(rest.ToMiddleware(gopkgmiddleware.ReserveRequest(gopkgmiddleware.ReserveRequestConfig{})))
	server.Use(rest.ToMiddleware(gopkgmiddleware.ReserveResponseWriter(gopkgmiddleware.ReserveResponseWriterConfig{})))

	handler.RegisterHandlers(server, svcCtx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		return svcCtx.ResponseHandler.GetStatusBody(nil, nil, err)
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
