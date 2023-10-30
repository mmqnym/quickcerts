package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Gracefully shutdown the server.
func WaitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	Logger.Info("The Server is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		Logger.Fatal("Something wrong happened when shutting down the server: " + err.Error())
	}

	Logger.Info("The Server has exited successfully.")

}