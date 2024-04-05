package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/meta-node-blockchain/meta-node/cmd/usdtnoti/internal/app"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
)

const (
	defaultConfigPath = "config.yaml"
	defaultLogLevel   = logger.FLAG_INFO
)

var (
	// flags
	CONFIG_FILE_PATH string
	LOG_LEVEL        int
)

// application flags

func main() {
	// init flags
	flag.StringVar(&CONFIG_FILE_PATH, "config", defaultConfigPath, "Config path")
	flag.StringVar(&CONFIG_FILE_PATH, "c", defaultConfigPath, "Config path (shorthand)")

	flag.IntVar(&LOG_LEVEL, "log-level", defaultLogLevel, "Log level")
	flag.IntVar(&LOG_LEVEL, "ll", defaultLogLevel, "Log level (shorthand)")

	flag.Parse()
	
	// init run app
	app, err := app.NewApp(CONFIG_FILE_PATH, LOG_LEVEL)
	if err != nil {
		panic(err)
	}
	go func() {
		app.Run()
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Stop()
		close(done)
	}()
	<-done
}
