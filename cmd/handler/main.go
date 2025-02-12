/*
Handler processes requests to map media list IDs to specific services and formats.

The sample implementation maps Anilist IDs to TVDB ones and formats the result as a Sonarr Custom List.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/textproto"
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
	"golang.org/x/time/rate"
	_ "modernc.org/sqlite"

	"github.com/wwmoraes/anilistarr/internal/adapters/cachedtracker"
	"github.com/wwmoraes/anilistarr/internal/adapters/chaincache"
	"github.com/wwmoraes/anilistarr/internal/adapters/sources"
	"github.com/wwmoraes/anilistarr/internal/api"
	"github.com/wwmoraes/anilistarr/internal/drivers/animelists"
	"github.com/wwmoraes/anilistarr/internal/drivers/redis"
	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

const (
	serviceNamespace = "github.com/wwmoraes/anilistarr"
	serviceName      = "anilistarr"

	apiInboundRateBurst      = 1000
	apiInboundRateInterval   = time.Minute
	refreshInterval          = time.Hour * 24 * 7
	gracefulShutdownTimeout  = 5 * time.Second
	requestReadHeaderTimeout = 5 * time.Second
)

var version = "0.0.0-unknown"

//nolint:funlen,maintidx // TODO tidy handler main fn
func main() {
	defer process.HandleExit()

	flags := struct {
		version bool
	}{
		version: false,
	}

	flag.BoolVar(&flags.version, "version", false, "shows version and exits")
	flag.Parse()

	if flags.version {
		fmt.Fprintln(os.Stdout, version)

		return
	}

	ctx := context.Background()

	log := logr.New(logging.NewStandardLogSink())
	ctx = logr.NewContext(ctx, log)

	err := telemetry.Initialize(ctx, resource.NewSchemaless(
		attribute.String("service.name", serviceName),
		attribute.String("service.namespace", serviceNamespace),
		attribute.String("service.version", version),
		attribute.String("host.id", getHostID(ctx)),
	))
	if err != nil {
		log.Error(err, "failed to initialize telemetry")
	}

	defer telemetry.Shutdown(ctx)
	defer telemetry.ForceFlush(ctx)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.Info("staring up", "name", serviceName, "version", version)

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dataPath := os.Getenv("DATA_PATH")

	store, err := newStore(dataPath)
	process.Assert(err)

	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisUsername := os.Getenv("REDIS_USERNAME")

	redisCache, err := redis.New(&redis.Options{
		Addr:       redisAddress,
		ClientName: serviceName,
		Password:   redisPassword,
		Username:   redisUsername,
	})
	process.Assert(err)

	fileCache, err := newCache(dataPath)
	process.Assert(err)

	cache := chaincache.ChainCache{
		redisCache,
		fileCache,
	}
	defer process.AssertClose(cache, "failed to close cache")

	tracker := cachedtracker.CachedTracker{
		Cache: cache,
		Tracker: anilist.New(
			os.Getenv("ANILIST_GRAPHQL_ENDPOINT"),
			anilist.WithPageSize(anilistPageSize),
		),
		TTL: cachedtracker.TTLs{
			UserID:       cacheUserTTL,
			MediaListIDs: cacheMediaListTTL,
		},
	}
	defer process.AssertClose(&tracker, "failed to close tracker")

	source := sources.JSON[animelists.Anilist2TVDBMetadata](
		"https://github.com/Fribb/anime-lists/raw/master/anime-list-full.json",
	)

	mediaLister := usecases.MediaList{
		Tracker: &tracker,
		Source:  source,
		Store:   store,
	}

	router := chi.NewRouter()
	router.Use(telemetry.WithInstrumentationMiddleware)
	router.Use(setHeaders(http.Header{
		"Cross-Origin-Resource-Policy": []string{"same-origin"},
		"X-Content-Type-Options":       []string{"nosniff"},
		"X-Frame-Options":              []string{"DENY"},
	}))
	router.Use(Limiter(rate.NewLimiter(
		rate.Every(apiInboundRateInterval),
		apiInboundRateBurst,
	)))

	service := api.Service{
		MediaLister: &mediaLister,
	}

	api.HandlerFromMux(&service, router)

	server := http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           router,
		ReadHeaderTimeout: requestReadHeaderTimeout,
	}

	// update mapping every week
	go scheduledRefresh(ctx, &mediaLister, refreshInterval)
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

		err := linker.Refresh(ctx, usecases.HTTPGetter(http.DefaultClient))
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

func setHeaders(headers http.Header) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header()[textproto.CanonicalMIMEHeaderKey(key)] = value
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getHostID(ctx context.Context) string {
	log := telemetry.Logr(ctx)

	id, err := machineid.ProtectedID(serviceNamespace + serviceName)
	if err != nil {
		log.Error(err, "failed to generate host ID")
	}

	return id
}
