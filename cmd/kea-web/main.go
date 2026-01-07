package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rannday/kea-web/internal/utils"
	"github.com/rannday/kea-web/internal/web"
)

const (
  shutdownTimeout = 30 * time.Second
)

func main() {
  utils.LoadEnv()
  env := utils.GetEnv()
  utils.ParseCLI(&env)
  utils.ValidateEnv(env)

  srv := web.NewServer("127.0.0.1:" + env.PORT)

  // Graceful shutdown signal handling
  shutdownChan := make(chan os.Signal, 1)
  signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

  go func() {
    utils.Info("Server running on http://localhost:%s", env.PORT)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      utils.Fatal("Server failed: %v", err)
    }
  }()

  sig := <-shutdownChan
  utils.Info("Received signal: %v. Shutting down...", sig)

  ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
  defer cancel()

  if err := srv.Shutdown(ctx); err != nil {
    utils.Error("Server shutdown failed: %v", err)
  } else {
    utils.Info("Server stopped.")
  }
}
