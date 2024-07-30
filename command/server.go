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

type ServerConfig struct {
	// The port number to listen on.
	Port int
	// The listen address to listen on.
	// Has precedence over Port.
	ListenAddress string
	// BasePath sets the prefix for the paths of the service.
	BasePath string
	// The logger.
	// It not set, it tries to get the default slog logger.
	// It it cannot find a default logger, it will log to null.
	Logger *slog.Logger
}

type serverCmd struct {
	listenAddress string
	baseURL       string
	logger        *slog.Logger
}

func (c *serverCmd) Run(ctx context.Context) error {

	baseContext := func(_ net.Listener) context.Context {
		return ctx
	}

	c.logger = c.logger.With(slog.String("listenAddress", c.listenAddress))

	serviceConfig := service.Config{
		BaseUrl:             c.baseURL,
		ErrorResponseFormat: service.ErrorResponseJSON,
	}

	// Create the server
	server := http.Server{
		Handler:           service.LogRequestMiddleware(service.NewService(&serviceConfig, c.logger).ServeHTTP, c.logger),
		Addr:              c.listenAddress,
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

func (c *serverCmd) setListenAddress(cfg *ServerConfig) {
	if cfg.ListenAddress != "" {
		c.listenAddress = cfg.ListenAddress
		return
	}
	if cfg.Port != 0 {
		c.listenAddress = fmt.Sprintf(":%d", cfg.Port)
		return
	}
}

func (s *serverCmd) setLogger(cfg *ServerConfig) {
	if cfg.Logger == nil {
		s.logger = slog.Default().With(slog.String("command", "server"))
		return
	}
	s.logger = cfg.Logger
}

func NewServerCmd(cfg *ServerConfig) *serverCmd {
	cmd := &serverCmd{}
	cmd.setListenAddress(cfg)
	cmd.setLogger(cfg)
	return cmd
}
