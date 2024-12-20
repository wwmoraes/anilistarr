package anilist

//go:generate go run github.com/Khan/genqlient@latest

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	telemetry "github.com/wwmoraes/gotell"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	// interval reflects the current Anilist API rate limits
	// https://anilist.gitbook.io/anilist-apiv2-docs/overview/rate-limiting
	interval time.Duration = time.Minute
	// requests reflects the current Anilist API rate limits
	// https://anilist.gitbook.io/anilist-apiv2-docs/overview/rate-limiting
	requests int = 90
)

// Tracker abstracts an Anilist GraphQL client and provides the common requests
// needed by MediaLister
type Tracker struct {
	Client   graphql.Client
	PageSize int
}

// New creates an Anilist Tracker and its GraphQL client
func New(anilistEndpoint string, pageSize int) usecases.Tracker {
	return &Tracker{
		Client:   NewGraphQLClient(anilistEndpoint),
		PageSize: pageSize,
	}
}

func NewGraphQLClient(anilistEndpoint string) graphql.Client {
	return graphql.NewClient(anilistEndpoint, NewRatedClient(interval, requests, nil))
}

func (tracker *Tracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := GetUserByName(ctx, tracker.Client, name)
	if err != nil {
		return "", span.Assert(fmt.Errorf(usecases.FailedGetUserErrorTemplate, err))
	}

	return strconv.Itoa(res.User.Id), span.Assert(nil)
}

func (tracker *Tracker) GetMediaListIDs(ctx context.Context, userId string) ([]string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	log := telemetry.Logr(ctx)

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf(usecases.ConvertUserIDErrorTemplate, err))
	}

	page := 1
	anilistIds := make([]string, 0, tracker.PageSize)
	span.SetAttributes(attribute.Int("page.size", tracker.PageSize))

	for {
		if ctx.Err() != nil {
			break
		}

		log.Info("requesting media list", "page", page)
		extCtx, extSpan := telemetry.StartNamed(
			ctx,
			"anilist.GetWatching",
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attribute.Int("page", page)),
		)
		res, err := GetWatching(extCtx, tracker.Client, userIdInt, page, tracker.PageSize)
		err = extSpan.Assert(err)
		extSpan.End()

		if err != nil {
			return nil, span.Assert(fmt.Errorf(usecases.FailedMediaErrorTemplate, err))
		}

		if len(res.Page.MediaList) == 0 {
			break
		}

		// far from optimal, I know, yet it works fine unless the use has thousands
		// of entries...
		for _, entry := range res.Page.MediaList {
			anilistIds = append(anilistIds, strconv.Itoa(entry.Media.Id))
		}

		page++
	}

	return anilistIds, span.Assert(nil)
}

func (tracker *Tracker) Close() error {
	return nil
}
