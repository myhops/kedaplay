package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"

	"kedaplay/command"
	"kedaplay/signalx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultLogger *slog.Logger

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kdp [ server | worker ]",
		Run: func(cmd *cobra.Command, args []string) {
			initDefaultLogger()
			cmd.Usage()
		},
		TraverseChildren: true,
	}
	cmd.Flags().String("loglevel", "info", "log level")
	cmd.Flags().String("logformat", "json", "log format")

	return cmd
}

func runServercmd(ccmd *cobra.Command, args []string) {
	initDefaultLogger()
	cfg := &command.ServerConfig{}
	if la := viper.GetString("listen-address"); la != "" {
		cfg.ListenAddress = la
	}
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = fmt.Sprintf(":%d", viper.GetUint("port"))
	}
	cfg.BasePath = viper.GetString("base-path")
	cfg.Logger = slog.Default().With(slog.String("command", "server"))
	command.NewServerCmd(cfg).Run(ccmd.Context())
}

func runWorkercmd(ccmd *cobra.Command, args []string) {
	initDefaultLogger()
	cfg := &command.WorkerConfig{}
	cfg.Resource = viper.GetString("resource")
	cfg.Sleep = int(viper.GetUint("sleep"))

	cfg.Logger = slog.Default().With(slog.String("command", "worker"))
	command.NewWorkerCmd(cfg).Run(ccmd.Context())
}

func newServerCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "server",
		Run:              runServercmd,
		TraverseChildren: true,
	}
	cmd.Flags().StringP("listen-address", "A", "", "Listen address. Use --port if you want to listen on all interfaces.")
	cmd.Flags().UintP("port", "p", 8080, "Listen port.")
	cmd.Flags().StringP("base-path", "B", "", "base prefix for url path.")

	cmd.Flags().AddFlagSet(parent.Flags())
	parent.AddCommand(cmd)
	return cmd
}

func newWorkerCmd(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: "worker",
		Run: runWorkercmd,
		TraverseChildren: true,
	}
	cmd.Flags().StringP("resource", "R", "http://localhost:8080/tasks/first", "Resource to pull tasks from")
	cmd.Flags().UintP("sleep", "S", 2, "Sleep interval")
	cmd.Flags().AddFlagSet(parent.Flags())
	parent.AddCommand(cmd)
	return cmd
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func initCobra() *cobra.Command {
	rootCmd := newRootCmd()
	serverCmd := newServerCmd(rootCmd)
	workerCmd := newWorkerCmd(rootCmd)

	must(viper.BindPFlags(rootCmd.Flags()))
	must(viper.BindPFlags(serverCmd.Flags()))
	must(viper.BindPFlags(workerCmd.Flags()))
	viper.AutomaticEnv()
	return rootCmd
}

// initialize the logger from viper.
func initDefaultLogger() {
	if defaultLogger != nil {
		return
	}

	var ll slog.Level
	switch strings.ToLower(viper.GetString("loglevel")) {
	case "debug":
		ll = slog.LevelDebug
	case "info":
		ll = slog.LevelInfo
	case "error":
		ll = slog.LevelError
	case "warn":
		ll = slog.LevelWarn
	}

	var lh slog.Handler
	lf := strings.ToLower(viper.GetString("logformat"))
	switch lf {
	case "text":
		lh = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: ll})
	case "json":
		lh = slog.NewJSONHandler(os.Stdout, nil)
	}
	defaultLogger := slog.New(lh)
	slog.SetDefault(defaultLogger)
}

func run(args []string) {
	nctx, cancel := signalx.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cmd := initCobra()
	cmd.SetArgs(args[1:])
	cmd.ExecuteContext(nctx)
}

func main() {
	run(os.Args)
}
