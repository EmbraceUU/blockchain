package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"web-blockchain/api/controller"
	"web-blockchain/internal/common/log"
	"web-blockchain/internal/config"
)

type Server struct {
	port string

	srv    *http.Server
	logger *logrus.Logger

	handler *controller.Handler
}

func NewServer(conf config.Config, handler *controller.Handler, log *logrus.Logger) *Server {
	var s Server
	s.logger = log
	s.port = conf.APIServer.Port
	s.handler = handler
	return &s
}

func (s *Server) Run() {
	//使用默认路由创建 http server
	s.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.port),
		Handler:      InitRouter(s.handler),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(err.Error())
		}
	}()
}

func (s *Server) Stop() {
	// 关闭server
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("shutdown server :" + err.Error())
	}
	select {
	case <-ctx.Done():
		s.logger.Error("server stop timeout of 5 seconds.")
	}
	s.logger.Info("server stop")
}
