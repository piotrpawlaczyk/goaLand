package database

import "context"

func (c *Client) ClearSeason(ctx context.Context, competitonID, startYear int) error {

	if err := c.DeleteSeasonMatches(ctx, competitonID, startYear); err != nil {
		return err
	}
	return c.DeleteSeasonPlayers(ctx, competitonID, startYear)
}

func (c *Client) DeleteSeasonMatches(ctx context.Context, competitonID, startYear int) error {
	filter := Filter{
		"competition_id":    competitonID,
		"start_season_year": startYear,
	}
	var matches []Match
	return c.db.WithContext(ctx).Where(filter).Delete(&matches).Error
}

func (c *Client) DeleteSeasonPlayers(ctx context.Context, competitonID, startYear int) error {
	filter := Filter{
		"competition_id":    competitonID,
		"start_season_year": startYear,
	}
	var players []SeasonPlayer
	return c.db.WithContext(ctx).Where(filter).Delete(&players).Error
}
