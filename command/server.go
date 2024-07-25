package command

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"kedaplay/service"
)

type serverCmd struct {
	listenAddress string
	logger        *slog.Logger
}

var getenv = func(name string) string {
	return os.Getenv(name)
}

func (c *serverCmd) Run(ctx context.Context, args []string, logger *slog.Logger) error {

	baseContext := func(_ net.Listener) context.Context {
		return ctx
	}

	var port = "8080"
	if p := getenv("PORT"); p != "" {
		port = p
	}
	port = ":" + port

	c.logger = logger.With(slog.String("module", "server"), "listenAddress", port)

	serviceConfig := service.Config{
		BaseUrl:     "",
		ErrorResponseFormat: service.ErrorResponseJSON,
	}

	// Create the server
	server := http.Server{
		Handler:           service.LogRequestMiddleware(service.NewService(&serviceConfig, c.logger).ServeHTTP, logger),
		Addr:              ":" + port,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext:       baseContext,
	}
	go server.ListenAndServe()
	c.logger.Info("listen started", slog.Duration("readHeaderTimeout", server.ReadHeaderTimeout))

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
	cmd := &serverCmd{}
	return cmd
}
