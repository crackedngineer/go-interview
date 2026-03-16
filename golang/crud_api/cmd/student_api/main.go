package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crackedngineer/go-interview/internal/config"
	"github.com/crackedngineer/go-interview/internal/http/handlers/student"
	"github.com/crackedngineer/go-interview/internal/storage/sqlite"
)

func main() {
	// load Config
	cfg := config.MustLoadConfig()

	// DB setup
	db, err := sqlite.NewDb(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(db))
	router.HandleFunc("GET /api/students/{id}", student.GetById(db))
	router.HandleFunc("GET /api/students", student.GetAll(db))
	// setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}
	fmt.Println("Server started at", cfg.HttpServer.Addr)

	// graceful shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("shutting down server")

}
