package cmd

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/eviltomorrow/blog-deploy-tool/internal/server"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/system"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/zlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "blog-deploy-server",
	Short: "",
	Long:  "  \r\nblog-deploy-server is a tool for blog html publish",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.Load(path, nil); err != nil {
			log.Fatalf("[Fatal] Load config failure, nest error: %v\r\n", err)
		}

		setVars := func() {
			server.Host = cfg.Server.Host
			server.Port = cfg.Server.Port
			server.CertsDir = filepath.Join(system.Pwd, "certs")
		}
		setLog()
		setVars()

		zlog.Info("Try to startup blog-deploy-server", zap.String("conf", cfg.String()))
		if err := server.StartupGRPC(); err != nil {
			zlog.Fatal("Startup GRPC server failure", zap.Error(err))
		}
		zlog.Info("blog-deploy-server startup complete", zap.String("status", "OK"))

		registerCleanFunc(server.ShutdownGRPC)
		blockingUntilTermination()
	},
}

var (
	cleanFuncs []func() error
	isServer   bool

	ServerName = "www.roigo.me"
	timeout    = 10 * time.Second
)

func init() {
	serverCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: false,
	}
	serverCmd.Flags().StringVarP(&path, "config", "c", "config.toml", "blog-deploy-server's config file")
}

func setLog() {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            cfg.Log.Level,
		Format:           cfg.Log.Format,
		DisableTimestamp: cfg.Log.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename:   filepath.Join(system.Pwd, "log", "data.log"),
			MaxSize:    cfg.Log.MaxSize,
			MaxDays:    30,
			MaxBackups: 30,
		},
		DisableStacktrace: true,
	})
	if err != nil {
		log.Fatalf("[Fatal] Setup log config failure, nest error: %v\r\n", err)
	}
	zlog.ReplaceGlobals(global, prop)
}

func registerCleanFunc(clean func() error) {
	cleanFuncs = append(cleanFuncs, clean)
}

func blockingUntilTermination() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}

	for _, f := range cleanFuncs {
		f()
	}
}

func NewServer() {
	isServer = true
	cobra.CheckErr(serverCmd.Execute())
}
