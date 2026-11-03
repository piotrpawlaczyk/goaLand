package database

import (
	"backend/internal/utils"
	"cmp"
	"context"
	"slices"
	"time"
)

func (c *Client) GetEditionMatches(ctx context.Context, competitionID int, startYear int) ([]Match, error) {
	statuses := slices.Concat(utils.UpcomingMatchStatuses, utils.LiveMatchStatuses)

	return c.GetEditionMatchesWithStatus(ctx, competitionID, startYear, statuses)
}

func (c *Client) GetEditionUpcommingMatches(ctx context.Context, competitionID int, startYear int) ([]Match, error) {
	return c.GetEditionMatchesWithStatus(ctx, competitionID, startYear, utils.UpcomingMatchStatuses)
}

func (c *Client) GetEditionLiveMatches(ctx context.Context, competitionID int, startYear int) ([]Match, error) {
	return c.GetEditionMatchesWithStatus(ctx, competitionID, startYear, utils.LiveMatchStatuses)
}

func (c *Client) GetEditionResult(ctx context.Context, competitionID int, startYear int) ([]Match, error) {
	matches, err := c.GetEditionMatchesWithStatus(ctx, competitionID, startYear, utils.FinishedMatchStatuses)
	if err != nil {
		return nil, err
	}
	slices.Reverse(matches)

	return matches, nil
}

func (c *Client) GetEditionMatchesWithStatus(ctx context.Context, competitionID int, startYear int,
	statuses []string) ([]Match, error) {
	return c.GetMatches(ctx, Filter{
		"start_season_year": startYear,
		"competition_id":    competitionID,
		"status":            statuses,
	})
}

func (c *Client) GetAllMatchesWithStatus(ctx context.Context, startYear int,
	statuses []string) ([]Match, error) {
	return c.GetMatches(ctx, Filter{
		"start_season_year": startYear,
		"status":            statuses,
	})
}

func (c *Client) GetMatches(ctx context.Context, filter Filter) ([]Match, error) {
	var matches []Match
	if err := c.List(ctx, &matches, filter, preloadTeams...); err != nil {
		return nil, err
	}

	sortMatches(matches)
	return matches, nil
}

func (c *Client) FeaturedMatchesQuery(from time.Time, to time.Time) func(
	ctx context.Context, competitionID int, startYear int) ([]Match, error) {
	return func(ctx context.Context, competitionID int, startYear int) ([]Match, error) {
		var matches []Match
		err := c.buildQuery(ctx, preloadTeams).Where(Filter{
			"competition_id":    competitionID,
			"start_season_year": startYear,
			"status":            utils.DisplayableMatchStatuses,
		}).Where("start_time >= ? AND start_time < ?", from, to).Find(&matches).Error
		if err != nil {
			return nil, err
		}

		sortMatches(matches)
		return matches, nil
	}
}

func sortMatches(matches []Match) {
	slices.SortFunc(matches, func(a, b Match) int {
		if result := a.StartTime.Compare(b.StartTime); result != 0 {
			return result
		}

		return cmp.Compare(a.MatchID, b.MatchID)
	})
}

func (c *Client) GetEditionPlayers(ctx context.Context, competitionID int, startYear int) ([]SeasonPlayer, error) {
	var goalScorers []SeasonPlayer
	err := c.List(ctx, &goalScorers, Filter{
		"competition_id":    competitionID,
		"start_season_year": startYear,
	},
	)
	return goalScorers, err
}

func (c *Client) GetEditionGoalScorers(ctx context.Context, competitionID int, startYear int) ([]SeasonPlayer, error) {
	var goalScorers []SeasonPlayer
	err := c.buildQuery(ctx, []string{"Player", "Team"}).
		Where(Filter{
			"competition_id":    competitionID,
			"start_season_year": startYear,
		}).
		Where("goals IS NOT NULL").
		Find(&goalScorers).Error

	return goalScorers, err
}

func (c *Client) GetEditions(ctx context.Context, competitionID int, preloads ...string) ([]Edition, error) {
	var edition []Edition
	err := c.List(ctx, &edition, Filter{"competition_id": competitionID}, preloads...)
	return edition, err
}

func (c *Client) GetEdition(ctx context.Context, competitionID int, startYear int, preloads ...string) (Edition, error) {
	var edition Edition
	err := c.Get(ctx, &edition, Filter{
		"competition_id": competitionID,
		"start_year":     startYear,
	}, preloads...)
	return edition, err
}
