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
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/audit"

	"github.com/NYTimes/gziphandler"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	configFile   = flag.String("f", "", "the config file")
	host         = flag.String("H", "", "the host, default: 0.0.0.0")
	port         = flag.Uint("P", 0, "the port, default: 6688")
	caFile       = flag.String("ca", "", "the ca file, default: certs/ca.crt")
	certFile     = flag.String("cert", "", "the cert file, default: certs/server.crt")
	keyFile      = flag.String("key", "", "the key file, default: certs/server.key")
	auditLogFile = flag.String("audit", "", "the audit log file, default: audit.log")
)

// 优先加载 config file
// 用 flag 的值替换掉 config file 的值
// 如果值为空，设置默认值
func parseFlags() config.Config {
	flag.Parse()

	var c config.Config
	c.Port = 6688
	c.Host = "0.0.0.0"
	c.CAFile = "certs/ca.crt"
	c.CertFile = "certs/server.crt"
	c.KeyFile = "certs/server.key"
	c.AuditLogFile = "audit.log"
	setFile := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "f" {
			setFile = true
		}
	})
	if setFile {
		conf.MustLoad(*configFile, &c, conf.UseEnv())
	}
	if *host != "" {
		c.Host = *host
	}
	if *port != 0 {
		c.Port = int(*port)
	}
	if *caFile != "" {
		c.CAFile = *caFile
	}
	if *certFile != "" {
		c.CertFile = *certFile
	}
	if *keyFile != "" {
		c.KeyFile = *keyFile
	}
	if *auditLogFile != "" {
		c.AuditLogFile = *auditLogFile
	}
	return c
}

func main() {
	c := parseFlags()
	// init audit log file
	if err := audit.InitLogFile(c.AuditLogFile); err != nil {
		log.Fatal(err)
	}

	// init tlsConfig
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
