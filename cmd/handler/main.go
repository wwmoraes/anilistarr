package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/go-chi/chi/v5"
	"github.com/go-logr/logr"
	telemetry "github.com/wwmoraes/gotell"
	"github.com/wwmoraes/gotell/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/time/rate"

	"github.com/wwmoraes/anilistarr/internal/api"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/functional"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

const (
	NAMESPACE = "github.com/wwmoraes/anilistarr"
	NAME      = "oteller"

	apiInboundRateBurst     = 1000
	refreshInterval         = time.Hour * 24 * 7
	gracefulShutdownTimeout = 5 * time.Second
)

var version = "0.2.0-8-g6002c65"

//nolint:funlen // TODO tidy handler main fn
func main() {
	defer process.HandleExit()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logr.New(logging.NewStandardLogSink())
	ctx = logr.NewContext(ctx, log)

	err := telemetry.Initialize(ctx, resource.NewSchemaless(
		attribute.String("service.name", NAME),
		attribute.String("service.namespace", NAMESPACE),
		attribute.String("service.version", version),
		attribute.String("host.id", functional.Unwrap(machineid.ProtectedID(NAMESPACE+NAME))),
	))
	if err != nil {
		log.Error(err, "failed to initialize telemetry")
	}

	log.Info("staring up", "name", NAME, "version", version)

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
	process.Assert(err)

	cache, err := NewCache(dataPath)
	process.Assert(err)

	mediaLister, err := NewAnilistMediaLister(
		os.Getenv("ANILIST_GRAPHQL_ENDPOINT"),
		store,
		cache,
	)
	process.Assert(err)

	r := chi.NewRouter()
	r.Use(telemetry.WithInstrumentationMiddleware)
	r.Use(Limiter(rate.NewLimiter(rate.Every(time.Minute), apiInboundRateBurst)))

	service, err := api.NewService(mediaLister)
	process.Assert(err)

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

func scheduledRefresh(ctx context.Context, linker usecases.MediaLister, interval time.Duration) {
	log := telemetry.Logr(ctx)

	for {
		log.Info("refreshing linker metadata")

		err := linker.Refresh(ctx, http.DefaultClient)
		process.Assert(err)

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
	log := logr.New(logging.NewStandardLogSink())

	log.Info("shutting down, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	process.Assert(server.Shutdown(ctx))
}
