package main

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/json"
)

// Version represents the application version which is set on compile time
var Version = "unknown"

func main() {
	app := &cli.App{
	  Name:      "{{ .command | dasherize | lowercase }}",
		HelpName:  "{{ .command | dasherize | lowercase }}",
		Usage:     "{{ .command | dasherize | lowercase }} http server",
		UsageText: "{{ .command | dasherize | lowercase }} [global options]",
		Version:   Version,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Action:    run,
		OnSignal:  signal,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen-addr, addr",
				Usage: "address on which the http server is listening on",
				Value: ":8080",
				EnvVar: "{{ .command | underscore | uppercase }}_LISTEN_ADDR",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "the log message level",
				Value: "info",
				EnvVar: "{{ .command | underscore | uppercase }}_LOG_LEVEL",
			},
		},
	}

	app.Run(os.Args)
}

func before(ctx *cli.Context) error {
	level, err := log.ParseLevel(ctx.String("log-level"))
	if err != nil {
		return err
	}

	log.SetLevel(level)
	log.SetHandler(json.New(ctx.Writer))

	return nil
}

func run(ctx *cli.Context) error {
	server := service.NewServer(&service.Config{
		Addr: ctx.String("listen-addr"),
	})

	// keep the server in the metadata in order to be accessible on signal
	ctx.Metadata["server"] = server

	logger := log.WithFields(log.Map{
		"addr": server.Addr,
	})

	logger.Info("http server is listening")
	return server.ListenAndServe()
}

func signal(ctx *cli.Context, term os.Signal) error {
	server, ok := ctx.Metadata["server"].(*http.Server)
	if !ok {
		return nil
	}

	logger := log.WithFields(log.Map{
		"addr":   server.Addr,
		"signal": term,
	})

	switch term {
	case syscall.SIGTERM:
		logger.Info("shutting down the server gracefully")

		if err := server.Shutdown(context.TODO()); err != nil {
			log.WithError(err).Error("failed to shutdown the http server gracefully")
			return err
		}
	default:
		logger.Info("unhandled signal occurred")
	}

	return nil
}
