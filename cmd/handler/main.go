package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"golang.org/x/time/rate"
)

type ServerContext string

const ListenerAddressKey ServerContext = "listener-address"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := telemetry.DefaultLogger()
	log.Info("testing new logger")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := "127.0.0.1:" + port

	mapper, err := NewAnilistLinker(
		os.Getenv("ANILIST_GRAPHQL_ENDPOINT"),
		&caches.RedisOptions{
			Addr:       os.Getenv("REDIS_ADDRESS"),
			Username:   os.Getenv("REDIS_USERNAME"),
			Password:   os.Getenv("REDIS_PASSWORD"),
			ClientName: "anilistarr-handler",
		},
	)
	assert(err)

	api, err := NewRestAPI(mapper)
	assert(err)

	shutdown, err := telemetry.InstrumentAll(ctx, os.Getenv("OTLP_ENDPOINT"))
	assert(err)
	defer shutdown(context.Background())

	ctx = telemetry.ContextWithLogger(ctx)

	r := chi.NewRouter()
	r.Use(telemetry.NewHandlerMiddleware)
	r.Use(Limiter(rate.NewLimiter(rate.Every(time.Minute), 10)))
	r.Get("/list", api.GetList)
	r.Get("/map", api.GetMap)
	r.Get("/user", api.GetUser)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// update mapping every week
	go scheduledRefresh(ctx, mapper, time.Hour*24*7)
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

func scheduledRefresh(ctx context.Context, linker *usecases.MediaLinker, interval time.Duration) {
	log := telemetry.LoggerFromContext(ctx)

	for {
		log.Info("refreshing linker metadata")
		err := linker.Refresh(ctx)
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
