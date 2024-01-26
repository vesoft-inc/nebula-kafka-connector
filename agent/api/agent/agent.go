package main

import (
	"compress/gzip"
	"embed"
	"flag"
	"fmt"
	"net/http"

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

//go:embed assets/*
var assetsFS embed.FS

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	svcCtx := svc.NewServiceContext(c)

	gzipHandler := gziphandler.MustNewGzipLevelHandler(gzip.DefaultCompression)
	server := rest.MustNewServer(c.RestConf,
		rest.WithNotFoundHandler(gzipHandler(
			gopkgmiddleware.NewAssetsHandler(gopkgmiddleware.AssetsConfig{
				Prefix:     "/assets/",
				Root:       "./assets/",
				SPA:        true,
				Filesystem: http.FS(assetsFS),
			}),
		)),
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
