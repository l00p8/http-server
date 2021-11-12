package http_server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/valve"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger"github.com/l00p8/log"
)

// Config describes server configuration
type Config struct {
	Addr            string        `envconfig:"addr" mapstructure:"addr" default:":8080"`
	ShutdownTimeout time.Duration `envconfig:"shutdown_timeout" mapstructure:"shutdown_timeout" default:"20"`
	GracefulTimeout time.Duration `envconfig:"graceful_timeout" mapstructure:"graceful_timeout" default:"21"`
	HealthUri       string        `envconfig:"health_uri" mapstructure:"health_uri" default:"/_health"`
	ApiVersion      string        `envconfig:"api_version" mapstructure:"api_version" default:"v1"`
	Timeout         time.Duration `envconfig:"timeout" mapstructure:"timeout" default:"20"`
	RateLimit       int64         `envconfig:"rate_limit" mapstructure:"rate_limit" default:"1"` // TODO: change in future :)
	Logger          logger.Logger
}

// Listen starts a http server on specified address and defines gateway routes
// Server implements a graceful shutdown pattern for better handling of rolling k8s updates
func Listen(cfg Config, router Muxer, cleanUp func()) error {
	valv := valve.New()
	log := cfg.Logger

	srv := http.Server{
		Addr:    cfg.Addr,
		Handler: router.Mux(),
	}

	router.Mux().Handle("/_metrics", promhttp.Handler())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Info("Shutting down a http server...\n")

			shutdown := cfg.ShutdownTimeout

			// first valv
			if err := valv.Shutdown(shutdown); err != nil {
				log.Error("Error shutting down a Valve: ", zap.String("error", err.Error()))
				return
			}

			// create a context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), shutdown)
			defer cancel()

			// start http server shutdown
			if err := srv.Shutdown(ctx); err != nil {
				log.Error("Error shutting down a http server: ", zap.String("error", err.Error()))
				return
			}

			// cleanUp before shutDown
			cleanUp()

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(cfg.GracefulTimeout):
				log.Info("Not all connections are done")
			case <-ctx.Done():

			}
		}
	}()

	log.Info("Starting a new server on address: ", zap.String("addr", cfg.Addr))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("A server listener error: ", zap.String("error", err.Error()))
		return err
	}
	log.Info("Server is down")
	return nil
}
