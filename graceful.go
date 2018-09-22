// Package graceful provides graceful shutdown for http.Server
package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Shutdown waits for interrupt signals from OS (SIGTERM/SIGINT).
// As soon as the signal arrives, it closes sockets, all
// incoming connections are being rejected. When all
// already came in requests are finished, Shutdown
// sends signal on shutdown channel
func Shutdown(server *http.Server) <-chan struct{} {
	osSigs := make(chan os.Signal)

	signal.Notify(osSigs, syscall.SIGTERM)
	signal.Notify(osSigs, syscall.SIGINT)

	shutdown := make(chan struct{})

	go func() {
		<-osSigs

		log.Println("Got OS interrupt signal")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Server Shutdown: %v", err)
		}

		shutdown <- struct{}{}
	}()

	return shutdown
}
