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

type anilistTracker struct {
	client   graphql.Client
	pageSize int
}

func New(anilistEndpoint string, pageSize int) (usecases.Tracker, error) {
	return &anilistTracker{
		client:   graphql.NewClient(anilistEndpoint, NewRatedClient(time.Minute, 90, nil)),
		pageSize: pageSize,
	}, nil
}

func (tracker *anilistTracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	res, err := GetUserByName(ctx, tracker.client, name)
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get user by name: %w", err))
	}

	return strconv.Itoa(res.User.Id), span.Assert(nil)
}

func (tracker *anilistTracker) GetMediaListIDs(ctx context.Context, userId string) ([]string, error) {
	log := telemetry.LoggerFromContext(ctx)
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to convert user ID to integer: %w", err))
	}

	page := 1
	anilistIds := make([]string, 0, tracker.pageSize)
	telemetry.Int(span, "page.size", tracker.pageSize)
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
		res, err := GetWatching(extCtx, tracker.client, userIdInt, page, tracker.pageSize)
		err = extSpan.Assert(err)
		extSpan.End()
		if err != nil {
			return nil, span.Assert(fmt.Errorf("failed to fetch media list: %w", err))
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

func (tracker *anilistTracker) Close() error {
	return nil
}
