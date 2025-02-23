package main

import (
	"fmt"
	"net"
	"net/http"
	"notes-go/internal/config"
	"notes-go/internal/user"
	"notes-go/pkg/logging"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("register user nadler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()

	logger.Info("Start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Infof("Listen unix socket: %s", socketPath)
		listener, listenErr = net.Listen("unix", socketPath)
	} else {
		logger.Info("listen tcp")
		address := fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
		listener, listenErr = net.Listen("tcp", address)
		logger.Infof("Server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatalln(server.Serve(listener))
}
