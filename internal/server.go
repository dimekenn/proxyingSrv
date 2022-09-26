package internal

import (
	"context"
	"crypto/tls"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"task_axxon/internal/configs"
	"task_axxon/internal/handler"
	"task_axxon/internal/service"
	"time"
)

//staring server
func StartHTTPServer(ctx context.Context, errCh chan<- error, cfg *configs.Configs) {
	app := echo.New()

	//creating http client
	cli := createHTTPClient()

	//creating service
	srv := service.NewService(cli, cfg)

	//creating handler
	srvHandler := handler.NewHandler(srv)

	//routes
	app.POST("/api/v1/redirect", srvHandler.RedirectApi)

	//staring http server
	errCh <- app.Start(cfg.SrvConfig.Port)
}

func createHTTPClient() *http.Client{
	return &http.Client{
		Transport:     &http.Transport{
			DialContext:            (&net.Dialer{
				Timeout:       30*time.Second,
				KeepAlive:     30*time.Second,
			}).DialContext,
			TLSClientConfig:        &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:    10*time.Second,
			MaxIdleConns:           100,
			IdleConnTimeout:        90*time.Second,
			ExpectContinueTimeout:  1*time.Second,
		},
	}
}
