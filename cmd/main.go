package main

import (
	"context"
	"github.com/ashtanko/octo-server/app"
	"github.com/ashtanko/octo-server/app/config"
	"github.com/ashtanko/octo-server/scheduler"
	"github.com/ashtanko/octo-server/store/sqlstore"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(); err != nil {
		logrus.Info("Failed to load dotenv file ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.LoadConfig()

	if err != nil {
		logrus.Fatalf("Failed initializing config: %s", err)
	}

	interruptChan := make(chan os.Signal, 1)

	sqlStore := sqlstore.CreateNewAndConnect(cfg)
	server, err := app.NewServer(sqlStore)

	if err != nil {
		logrus.Fatalf("Error run init: %s", err.Error())
		panic(err.Error())
	}

	defer server.Shutdown()

	serverErr := server.Start(cfg.Port)
	if serverErr != nil {
		logrus.Fatalf("Error run server: %s", serverErr.Error())
	}

	scheduler.Init(ctx, sqlStore)

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-interruptChan:
		cancel()
	}
}
