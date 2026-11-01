package sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/api"
	"backend/internal/database"
	"backend/internal/sync/mocks"
	"backend/internal/testutils"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func strPtr(s string) *string { return &s }

func testSync(t *testing.T) (*gomock.Controller, *mocks.MockAPI, *mocks.MockDatabase, *Sync) {
	t.Helper()

	ctrl := gomock.NewController(t)

	apiMock := mocks.NewMockAPI(ctrl)
	dbMock := mocks.NewMockDatabase(ctrl)

	return ctrl, apiMock, dbMock, New(apiMock, dbMock, nil)
}

func TestInitAreas(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, dbMock, s := testSync(t)
	defer ctrl.Finish()

	england := testutils.England()
	poland := testutils.Poland()

	apiMock.EXPECT().FetchAreas().Return([]api.Area{
		{ID: testutils.PolandAreaID, Name: "Poland", CountryCode: testutils.PremierLeagueCode, ParentArea: strPtr("Europe")},
		{ID: testutils.EnglandAreaID, Name: testutils.EnglandName, CountryCode: "EN", ParentArea: strPtr("Europe")},
	}, nil)

	dbMock.EXPECT().Save(ctx, []database.Area{
		poland,
		england,
	}).Return(nil)
	dbMock.EXPECT().Save(ctx, &unknownArea).Return(nil)

	require.NoError(t, s.initAreas(ctx))

	require.Equal(t, s.areasByName["Poland"], testutils.PolandAreaID)
	require.Equal(t, s.areasByName[testutils.EnglandName], testutils.EnglandAreaID)
}

func TestInitCompetition(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, dbMock, s := testSync(t)
	defer ctrl.Finish()

	apiMock.EXPECT().FetchCompetition(testutils.PremierLeagueCode).Return(api.Competition{
		ID:   2021,
		Name: "Premier League",
		Code: testutils.PremierLeagueCode,
		Type: "LEAGUE",
		Area: struct {
			ID int `json:"id"`
		}{ID: testutils.EnglandAreaID},
	}, nil)

	s.areasByName[testutils.EnglandName] = testutils.EnglandAreaID

	competition := testutils.NewCompetition()
	expectedCompetition := competition
	expectedCompetition.Area = database.Area{}

	dbMock.EXPECT().Save(ctx, []database.Competition{expectedCompetition}).Return(nil)

	id, err := s.initCompetition(ctx, competition.Code)
	require.NoError(t, err)

	require.Equal(t, id, competition.CompetitionID)
}

func TestInitCompetitionFetchError(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, _, s := testSync(t)
	defer ctrl.Finish()

	expectedErr := errors.New("api error")

	apiMock.EXPECT().
		FetchCompetition(testutils.PremierLeagueCode).
		Return(api.Competition{}, expectedErr)

	_, err := s.initCompetition(ctx, testutils.PremierLeagueCode)

	require.ErrorIs(t, err, expectedErr)
}

func TestInitCompetitions(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, dbMock, s := testSync(t)
	defer ctrl.Finish()

	premierLeague := testutils.PremierLeague()
	ekstraklasa := testutils.Ekstraklasa()
	expectedPremierLeague := premierLeague
	expectedPremierLeague.Area = database.Area{}

	apiMock.EXPECT().FetchCompetitions().Return([]api.Competition{
		{
			ID:   premierLeague.CompetitionID,
			Name: premierLeague.Name,
			Code: premierLeague.Code,
			Type: premierLeague.CompetitionType,
			Area: struct {
				ID int `json:"id"`
			}{
				ID: premierLeague.AreaID,
			},
		}, {
			ID:   ekstraklasa.CompetitionID,
			Name: ekstraklasa.Name,
			Code: ekstraklasa.Code,
			Type: ekstraklasa.CompetitionType,
			Area: struct {
				ID int `json:"id"`
			}{
				ID: ekstraklasa.AreaID,
			},
		},
	}, nil)

	dbMock.EXPECT().Save(ctx, []database.Competition{expectedPremierLeague}).Return(nil)

	err := s.initCompetitions(ctx, map[string]struct{}{
		premierLeague.Code: {},
	})
	require.NoError(t, err)
}

func TestInitCompetitionsMissingCodes(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, _, s := testSync(t)
	defer ctrl.Finish()

	apiMock.EXPECT().FetchCompetitions().Return([]api.Competition{}, nil)

	err := s.initCompetitions(ctx, map[string]struct{}{
		testutils.PremierLeagueCode: {},
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "competitions not found")
}

func TestInitMatches(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, dbMock, s := testSync(t)
	defer ctrl.Finish()

	match := testutils.NewMatch()
	expectedMatch := match
	expectedMatch.HomeTeam = database.Team{}
	expectedMatch.AwayTeam = database.Team{}
	expectedMatch.Edition = database.Edition{}

	apiMatch := api.Match{
		ID:       match.MatchID,
		UtcDate:  match.StartTime.Format(time.RFC3339),
		Status:   match.Status,
		Matchday: match.Matchday,
		Stage:    match.Stage,
	}
	apiMatch.HomeTeam.ID = match.HomeTeamID
	apiMatch.AwayTeam.ID = match.AwayTeamID

	apiMatch.Score.FullTime.Home = match.HomeGoals
	apiMatch.Score.FullTime.Away = match.AwayGoals

	apiMock.EXPECT().FetchMatches(testutils.PremierLeagueCode, testutils.Year).Return([]api.Match{apiMatch}, nil)

	dbMock.EXPECT().SaveAndCheck(ctx, []database.Match{expectedMatch}).Return(true, nil)

	_, err := s.initMatches(ctx, Season{
		CompetitionID:   testutils.PremierLeagueID,
		CompetitionCode: testutils.PremierLeagueCode,
		StartYear:       testutils.Year,
	})

	require.NoError(t, err)
}

func TestInitGoalScorers(t *testing.T) {
	ctx := t.Context()

	ctrl, apiMock, dbMock, s := testSync(t)
	defer ctrl.Finish()

	player := testutils.NewPlayer()
	seasonPlayer := testutils.NewSeasonPlayer()

	apiScorer := api.GoalScorer{
		Player: api.GoalScorerPlayer{
			ID:          player.PlayerID,
			Name:        player.Name,
			Nationality: testutils.EnglandName,
			Section:     player.Position,
		},
		Goals:     seasonPlayer.Goals,
		Assists:   seasonPlayer.Assists,
		Penalties: seasonPlayer.GoalsFromPenalty,
	}

	apiScorer.Team.ID = seasonPlayer.TeamID

	s.areasByName[testutils.EnglandName] = testutils.EnglandAreaID

	expectedPlayer := database.Player{
		PlayerID:          player.PlayerID,
		Name:              player.Name,
		Position:          player.Position,
		NationalityAreaID: player.NationalityAreaID,
	}
	expectedSeasonPlayer := database.SeasonPlayer{
		PlayerID:         seasonPlayer.PlayerID,
		TeamID:           seasonPlayer.TeamID,
		CompetitionID:    seasonPlayer.CompetitionID,
		StartSeasonYear:  seasonPlayer.StartSeasonYear,
		Goals:            seasonPlayer.Goals,
		Assists:          seasonPlayer.Assists,
		GoalsFromPenalty: seasonPlayer.GoalsFromPenalty,
	}

	apiMock.EXPECT().
		FetchGoalScorers(testutils.PremierLeagueCode, testutils.Year, defaultGoalScorerLimit).
		Return([]api.GoalScorer{apiScorer}, nil)

	dbMock.EXPECT().Save(ctx, []database.Player{expectedPlayer}).
		Return(nil)
	dbMock.EXPECT().Save(ctx, []database.SeasonPlayer{expectedSeasonPlayer}).
		Return(nil)

	err := s.initGoalScorers(ctx, Season{
		CompetitionID:   testutils.PremierLeagueID,
		CompetitionCode: testutils.PremierLeagueCode,
		StartYear:       testutils.Year,
	}, defaultGoalScorerLimit)

	require.NoError(t, err)
}

func TestUpdateTeamFetchError(t *testing.T) {
	target := SeasonTarget{
		CompetitionCode: testutils.PremierLeagueCode,
		StartYear:       testutils.Year,
	}
	season := Season{
		CompetitionID:   testutils.PremierLeagueID,
		CompetitionCode: target.CompetitionCode,
		StartYear:       target.StartYear,
	}

	ctx := t.Context()
	_, apiMock, dbMock, s := testSync(t)
	s.seasons = Seasons{season}
	expectedErr := errors.New("failed to fetch teams")

	apiMock.EXPECT().FetchTeams(target.CompetitionCode, target.StartYear).
		Return(nil, expectedErr)
	dbMock.EXPECT().DoTransaction(gomock.Any(), gomock.Any()).Times(0)
	dbMock.EXPECT().ClearSeason(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	dbMock.EXPECT().SaveAndCheck(gomock.Any(), gomock.Any()).Times(0)

	require.ErrorIs(t, s.updateEdition(ctx, target), expectedErr)
}

func TestUpdateMatchFetchError(t *testing.T) {
	target := SeasonTarget{
		CompetitionCode: testutils.PremierLeagueCode,
		StartYear:       testutils.Year,
	}
	season := Season{
		CompetitionID:   testutils.PremierLeagueID,
		CompetitionCode: target.CompetitionCode,
		StartYear:       target.StartYear,
	}

	ctx := t.Context()
	ctrl, apiMock, dbMock, s := testSync(t)
	txMock := mocks.NewMockDatabase(ctrl)
	s.seasons = Seasons{season}
	s.areasByName[testutils.EnglandName] = testutils.EnglandAreaID
	expectedErr := errors.New("failed to fetch matches")

	player := testutils.NewPlayer()
	apiMock.EXPECT().FetchTeams(target.CompetitionCode, target.StartYear).
		Return([]api.Team{{
			ID: testutils.ArsenalID,
			Players: []api.Player{{
				ID:          player.PlayerID,
				Name:        player.Name,
				Nationality: testutils.EnglandName,
				Position:    player.Position,
			}},
		}}, nil)

	dbMock.EXPECT().DoTransaction(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, fn func(Database) error) error {
			err := fn(txMock)
			require.ErrorIs(t, err, expectedErr)
			return err
		})
	txMock.EXPECT().ClearSeason(ctx, season.CompetitionID, season.StartYear).Return(nil)
	txMock.EXPECT().Save(ctx, []database.Player{{
		PlayerID:          player.PlayerID,
		Name:              player.Name,
		NationalityAreaID: testutils.EnglandAreaID,
		Position:          player.Position,
	}}).Return(nil)
	txMock.EXPECT().Save(ctx, []database.SeasonPlayer{{
		PlayerID:        player.PlayerID,
		TeamID:          testutils.ArsenalID,
		CompetitionID:   season.CompetitionID,
		StartSeasonYear: season.StartYear,
	}}).Return(nil)
	apiMock.EXPECT().FetchMatches(target.CompetitionCode, target.StartYear).
		Return(nil, expectedErr)
	txMock.EXPECT().SaveAndCheck(gomock.Any(), gomock.Any()).Times(0)
	dbMock.EXPECT().SaveAndCheck(gomock.Any(), gomock.Any()).Times(0)
	apiMock.EXPECT().FetchGoalScorers(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	require.ErrorIs(t, s.updateEdition(ctx, target), expectedErr)
}
