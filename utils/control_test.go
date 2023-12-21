package utils

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitForShutdown(t *testing.T) {
	// Create a simulated HTTP server.
	server := &http.Server{Addr: ":0"}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			t.Logf("HTTP server ListenAndServe: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		WaitForShutdown(server)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	quit <- os.Interrupt

	time.Sleep(100 * time.Millisecond)

	// Try to connect to the server which should be shutdown.
	_, err := net.Dial("tcp", server.Addr)
	assert.NotNil(t, err, "Expected server to be shutdown, but it was still accessible")
}
