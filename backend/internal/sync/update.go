package sync

import (
	"backend/internal/api"
	"context"
	"errors"
)

func updateTransaction(ctx context.Context, originalSync *Sync, databaseClient Database, season Season, apiTeams []api.Team) (bool, error) {
	var change bool
	err := databaseClient.DoTransaction(ctx, func(databaseClient Database) error {
		transactionSync := originalSync.replaceDatabase(databaseClient)
		if err := databaseClient.ClearSeason(ctx, season.CompetitionID, season.StartYear); err != nil {
			return err
		}

		if err := transactionSync.initPlayers(ctx, apiTeams, season); err != nil {
			return err
		}

		var err error
		change, err = transactionSync.initMatches(ctx, season)
		if err != nil {
			return err
		}

		if err := transactionSync.initGoalScorers(ctx, season, defaultGoalScorerLimit); err != nil {
			return err
		}
		return nil
	})
	return change, err
}
func (s *Sync) updateEdition(ctx context.Context, target SeasonTarget) error {
	competitionID, ok := s.seasons.CompetitionID(target.CompetitionCode)

	if !ok {
		return errors.New("not found competition ID")
	}
	season := Season{
		CompetitionID:   competitionID,
		CompetitionCode: target.CompetitionCode,
		StartYear:       target.StartYear,
	}

	apiTeams, err := s.apiClient.FetchTeams(season.CompetitionCode, season.StartYear)
	if err != nil {
		return err
	}

	change, err := updateTransaction(ctx, s, s.databaseClient, season, apiTeams)
	if err != nil {
		return err
	}
	if change {
		if err := s.storage.RefreshTable(ctx, season.CompetitionID, season.StartYear); err != nil {
			return err
		}
	}
	return nil
}
