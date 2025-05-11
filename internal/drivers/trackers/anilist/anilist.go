// Package anilist provides an [usecases.Tracker] that communicates with
// Anilist's GraphQL API.
package anilist

//go:generate go run github.com/Khan/genqlient

import (
	"context"
	"errors"
	"io"
	"iter"
	"net/http"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	telemetry "github.com/wwmoraes/gotell"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"

	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/with"
)

const (
	// interval reflects the current Anilist API rate limits
	// https://anilist.gitbook.io/anilist-apiv2-docs/overview/rate-limiting
	interval time.Duration = time.Minute
	// requests reflects the current Anilist API rate limits
	// https://anilist.gitbook.io/anilist-apiv2-docs/overview/rate-limiting
	requests        int = 30
	defaultPageSize int = 10
)

var _ io.Closer = (*Tracker)(nil)

// Tracker abstracts an Anilist GraphQL client and provides the common requests
// needed by MediaLister
type Tracker struct {
	Client   graphql.Client
	PageSize int
}

// Options contains optional settings for tracker instances.
type Options struct {
	Client   usecases.Doer
	PageSize int
}

// NewOptions generates an [Options] value ready to use. It starts with defaults
// and then applies all [Option] in order.
func NewOptions(opts ...with.Option[Options]) Options {
	return with.Apply(Options{
		PageSize: defaultPageSize,
		Client:   http.DefaultClient,
	}, opts...)
}

// WithClient sets a custom HTTP client.
func WithClient(client usecases.Doer) with.Functor[Options] {
	return with.Functor[Options](func(options *Options) {
		options.Client = client
	})
}

// WithPageSize sets a custom page size for paginated requests.
func WithPageSize(size int) with.Functor[Options] {
	return with.Functor[Options](func(options *Options) {
		options.PageSize = size
	})
}

// New creates an Anilist client that uses a [RatedClient] that respects the
// upstream API limits.
func New(endpoint string, opts ...with.Option[Options]) *Tracker {
	options := NewOptions(opts...)

	return &Tracker{
		Client: graphql.NewClient(
			endpoint,
			&RatedClient{
				Doer:    options.Client,
				Limiter: rate.NewLimiter(rate.Every(interval), requests),
			},
		),
		PageSize: options.PageSize,
	}
}

// GetUserID retrieves an user ID using their profile name.
func (tracker *Tracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := GetUserByName(ctx, tracker.Client, name)

	gqlErrorList := gqlerror.List{}
	if errors.As(err, &gqlErrorList) && len(gqlErrorList) > 0 && gqlErrorList[0].Message == "Not Found." {
		return "", span.Assert(usecases.ErrStatusNotFound)
	}

	if err != nil {
		return "", span.Assert(errors.Join(usecases.ErrStatusUnavailable, err))
	}

	return strconv.Itoa(res.User.Id), span.Assert(nil)
}

// GetMediaListIDs retrieves a list of medias from a user ID.
func (tracker *Tracker) GetMediaListIDs(ctx context.Context, userID string) ([]string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, span.Assert(errors.Join(usecases.ErrStatusInvalidArgument, err))
	}

	anilistIDs := make([]string, 0, tracker.PageSize)
	span.SetAttributes(attribute.Int("page.size", tracker.PageSize))

	for mediaID := range tracker.getMediaListIDs(ctx, userIDInt) {
		anilistIDs = append(anilistIDs, strconv.Itoa(mediaID))
	}

	return anilistIDs, span.Assert(nil)
}

// Close terminates the client to the upstream API.
func (tracker *Tracker) Close() error {
	tracker.Client = nil

	return nil
}

//nolint:gocognit // gotta have those short-circuit returns ¯\_(ツ)_/¯
func (tracker *Tracker) getMediaListIDs(ctx context.Context, userID int) iter.Seq[int] {
	span := telemetry.SpanFromContext(ctx)

	return func(yield func(int) bool) {
		for page := 1; ; page++ {
			res, err := tracker.getWatchingPage(ctx, userID, page)
			if err != nil {
				//nolint:errcheck // logged in span
				span.Assert(errors.Join(usecases.ErrStatusUnknown, err))

				return
			}

			for _, entry := range res.Page.MediaList {
				if !yield(entry.Media.Id) {
					return
				}

				if ctx.Err() != nil {
					return
				}
			}

			// avoid doing requests after exhausting records
			if len(res.Page.MediaList) < tracker.PageSize {
				return
			}
		}
	}
}

func (tracker *Tracker) getWatchingPage(ctx context.Context, userID, page int) (*GetWatchingResponse, error) {
	ctx, span := telemetry.StartNamed(
		ctx,
		"anilist.GetWatching",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.Int("page", page)),
	)
	defer span.End()

	log := telemetry.Logr(ctx)

	log.Info("requesting media list", "page", page)
	res, err := GetWatching(ctx, tracker.Client, userID, page, tracker.PageSize)

	return res, span.Assert(err)
}
