package main

import (
	"os"
	"os/signal"
	"syscall"
	"web-blockchain/api/controller"
	"web-blockchain/api/server"
	svc "web-blockchain/api/svc"
	"web-blockchain/internal/common/log"
	"web-blockchain/internal/config"
	"web-blockchain/internal/consts"
)

func main() {
	conf := config.NewConfig()

	serviceLogger := log.WithLoggerName(consts.ServiceLoggerName)
	outLogger := log.WithLoggerName(consts.OutLoggerName)

	handler := controller.NewHandler(svc.NewServiceContext(conf, serviceLogger), serviceLogger, outLogger)

	s := server.NewServer(conf, handler, serviceLogger)
	s.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Stop()
}
