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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		Logger.Fatal("Something wrong happened when shutting down the server: " + err.Error())
	}

	Logger.Info("The Server has exited successfully.")
}

// Ensure that the current working directory is the root directory of the project.
func Change2RootDir() bool {
	if _, err := os.Stat("go.mod"); !os.IsNotExist(err) {
		return false
	}

	for {
		if _, err := os.Stat("go.mod"); !os.IsNotExist(err) {
			break
		}

		if err := os.Chdir(".."); err != nil {
			panic("can not find root directory of the project")
		}

		curr, err := os.Getwd()
		if err != nil {
			panic("can not find root directory of the project")
		}

		if curr == "/" {
			panic("can not find go.mod file")
		}
	}

	root, _ := os.Getwd()
	os.Chdir(root)
	return true
}
