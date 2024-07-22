package command

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"kedaplay/service"
)

type serverCmd struct{}

type ServerOptions struct{}

func (c *serverCmd) Run(ctx context.Context, args []string, logger *slog.Logger) error {

	baseContext := func(_ net.Listener) context.Context {
		return ctx
	}

	// Create the server
	server := http.Server{
		Handler:           service.LogRequestMiddleware(service.NewService().ServeHTTP, logger),
		Addr:              ":8080",
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext:       baseContext,
	}
	go server.ListenAndServe()
	log.Print("listening on :8080")

	<-ctx.Done()
	log.Print("Done closed")
	if ctx.Err() != nil {
		log.Printf("error: %s", ctx.Err().Error())
	}

	// shutdown the server
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := server.Shutdown(sctx)
	if err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}
	log.Print("server shut down normally")
	return nil
}

func NewServerCmd() *serverCmd {
	return &serverCmd{}
}
