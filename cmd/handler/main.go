package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/time/rate"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := telemetry.InstrumentAll(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	log := telemetry.DefaultLogger()
	if errors.Is(err, telemetry.NoEndpointError) {
		log.Error(err, "skipping instrumentation")
		err = nil
	} else {
		defer shutdown(context.Background())
	}
	assert(err)

	log.Info("staring up", "name", telemetry.NAME, "version", telemetry.VERSION)

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	mediaLister, err := NewAnilistMediaLister(
		os.Getenv("ANILIST_GRAPHQL_ENDPOINT"),
		os.Getenv("DATA_PATH"),
	)
	assert(err)

	api := &RestAPI{mediaLister}
	ctx = telemetry.ContextWithLogger(ctx)

	r := chi.NewRouter()
	r.Use(otelchi.Middleware(telemetry.NAME, otelchi.WithChiRoutes(r)))
	r.Use(telemetry.NewHandlerMiddleware)
	r.Use(Limiter(rate.NewLimiter(rate.Every(time.Minute), 1000)))
	r.Get("/list", api.GetList)
	r.Get("/map", api.GetMap)
	r.Get("/user", api.GetUser)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// update mapping every week
	go scheduledRefresh(ctx, mediaLister, time.Hour*24*7)
	go server.ListenAndServe() //nolint:errcheck
	log.Info("server listening", "address", addr)

	<-ctx.Done()
	cancel()

	gracefulShutdown(&server)
}

func assert(err error) {
	if err == nil {
		return
	}

	log := telemetry.DefaultLogger()

	log.Error(err, "assertion failed")
	os.Exit(1)
}

func scheduledRefresh(ctx context.Context, linker usecases.MediaLister, interval time.Duration) {
	log := telemetry.LoggerFromContext(ctx)

	for {
		log.Info("refreshing linker metadata")
		err := linker.Refresh(ctx, http.DefaultClient)
		assert(err)
		log.Info("linker metadata refreshed")

		select {
		case <-ctx.Done():
			log.Info("scheduled refresh stopped")
			return
		case <-time.After(interval):
			continue
		}
	}
}

func gracefulShutdown(server *http.Server) {
	log := telemetry.DefaultLogger()

	log.Info("shutting down, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	assert(server.Shutdown(ctx))
}
