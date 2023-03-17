package main

import (
	"app/internal"
	"app/internal/config"
	"app/internal/db"
	"app/pkg/logging"
	"app/pkg/mongodb"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create config")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password,
		cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	task1 := internal.Task{
		ID:       "",
		Task:     "Write code",
		Deadline: "ever",
		Status:   true,
	}

	task1Id, err := storage.Create(context.Background(), task1)
	logger.Info(task1Id)

	if err != nil {
		panic(err)
	}
	task2 := internal.Task{
		ID:       "",
		Task:     "готовьdddd",
		Deadline: "скоро",
		Status:   true,
	}

	task2Id, err := storage.Create(context.Background(), task2)
	if err != nil {
		panic(err)
	}
	logger.Info(task2Id)
	userFound, err := storage.FindOne(context.Background(), task2Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(userFound)

	userFound.Status = false

	err = storage.UpdateTask(context.Background(), userFound)
	if err != nil {
		panic(err)
	}

	err = storage.DeleteTask(context.Background(), task2Id)
	if err != nil {
		panic(err)
	}
	_, err = storage.FindOne(context.Background(), task2Id)
	if err != nil {
		panic(err)
	}

	logger.Info("register author handler")
	handler := internal.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
