package server

import (
	"context"
	"net/http"
	"time"

	"github.com/dafsic/gambler/config"
	"github.com/dafsic/gambler/lib/mylog"
	"github.com/dafsic/gambler/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Server interface {
	RegisterHandler(method, path string, h gin.HandlerFunc)
}

type ServerImpl struct {
	listenAddr string
	l          *utils.Logger
	srv        *http.Server
	gin        *gin.Engine
}

func NewServer(lc fx.Lifecycle, cfg config.ConfigI, log mylog.Logging) Server {
	s := ServerImpl{
		l:          log.GetLogger("http"),
		listenAddr: cfg.GetElem("listen").(string),
	}
	s.l.Info("Init...")

	s.gin = gin.New()
	//s.gin.Use(s.l)
	s.srv = &http.Server{
		Addr:         s.listenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.gin,
	}

	lc.Append(fx.Hook{
		// app.start调用
		OnStart: func(ctx context.Context) error {
			// 这里不能阻塞
			go func() {
				if err := s.srv.ListenAndServe(); err != nil {
					s.l.Error(err)
				}
			}()
			return nil
		},
		// app.stop调用，收到中断信号的时候调用app.stop
		OnStop: func(ctx context.Context) error {
			s.srv.Shutdown(ctx)
			return nil
		},
	})

	return &s
}

func (s ServerImpl) RegisterHandler(method, path string, h gin.HandlerFunc) {
	switch method {
	case "get":
		s.gin.GET(path, h)
	case "post":
		s.gin.POST(path, h)
	default:
	}
}

// Module for fx
var ServerModule = fx.Options(fx.Provide(NewServer))
