package anilist

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
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

type Tracker struct {
	Client   graphql.Client
	PageSize int
}

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
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	res, err := GetUserByName(ctx, tracker.Client, name)
	if err != nil {
		return "", span.Assert(fmt.Errorf(usecases.FailedGetUserErrorTemplate, err))
	}

	return strconv.Itoa(res.User.Id), span.Assert(nil)
}

func (tracker *Tracker) GetMediaListIDs(ctx context.Context, userId string) ([]string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()
	log := telemetry.LoggerFromContext(ctx)

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf(usecases.ConvertUserIDErrorTemplate, err))
	}

	page := 1
	anilistIds := make([]string, 0, tracker.PageSize)
	telemetry.Int(span, "page.size", tracker.PageSize)
	for {
		if ctx.Err() != nil {
			break
		}

		log.Info("requesting media list", "page", page)
		extCtx, extSpan := telemetry.Start(
			ctx,
			"anilist.GetWatching",
			telemetry.WithSpanKindClient(),
			telemetry.WithInt("page", page),
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
