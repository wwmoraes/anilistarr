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
	"github.com/wwmoraes/anilistarr/internal/api"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/time/rate"
)

const (
	apiInboundRateBurst     = 1000
	refreshInterval         = time.Hour * 24 * 7
	gracefulShutdownTimeout = 5 * time.Second
)

//nolint:funlen // TODO tidy handler main fn
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := telemetry.InstrumentAll(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	log := telemetry.DefaultLogger()
	if errors.Is(err, telemetry.ErrNoEndpoint) {
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

	dataPath := os.Getenv("DATA_PATH")

	store, err := NewStore(dataPath)
	assert(err)

	cache, err := NewCache(dataPath)
	assert(err)

	mediaLister, err := NewAnilistMediaLister(
		os.Getenv("ANILIST_GRAPHQL_ENDPOINT"),
		store,
		cache,
	)
	assert(err)

	ctx = telemetry.ContextWithLogger(ctx)

	r := chi.NewRouter()
	r.Use(otelchi.Middleware(telemetry.NAME, otelchi.WithChiRoutes(r)))
	r.Use(telemetry.NewHandlerMiddleware)
	r.Use(Limiter(rate.NewLimiter(rate.Every(time.Minute), apiInboundRateBurst)))

	service, err := api.NewService(mediaLister)
	assert(err)

	api.HandlerFromMux(service, r)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}

	// update mapping every week
	go scheduledRefresh(ctx, mediaLister, refreshInterval)
	//nolint:errcheck // ignore listen errors
	go server.ListenAndServe()

	log.Info("server listening", "address", server.Addr)

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

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	assert(server.Shutdown(ctx))
}
