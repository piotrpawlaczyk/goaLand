package handler

import (
	"backend/internal/database"
	"backend/internal/sync"
	"backend/internal/transform"
	"backend/internal/utils"
	"context"
	"net/http"
	"time"
)

const (
	UpcomingStatus = "upcoming"
	LiveStatus     = "live"
	FinishedStatus = "finished"
)

const maxMatchesOnTop = 50

func (h *Handler) GetAllMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startYear := sync.CurrentSeasonStartYear(time.Now())

	var allMatches [][]transform.MatchResponse

	matches, err := h.db.GetAllMatchesWithStatus(ctx, startYear, utils.UpcomingMatchStatuses)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	allMatches = append(allMatches,
		transform.GetResponseMatchesWithLimit(matches, maxMatchesOnTop, true))

	matches, err = h.db.GetAllMatchesWithStatus(ctx, startYear, utils.LiveMatchStatuses)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	allMatches = append(allMatches, transform.GetFeaturedMatches(matches))

	matches, err = h.db.GetAllMatchesWithStatus(ctx, startYear, utils.FinishedMatchStatuses)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	allMatches = append(allMatches,
		transform.GetResponseMatchesWithLimit(matches, maxMatchesOnTop, false))

	WriteJSON(w, http.StatusOK, allMatches)
}

func (h *Handler) GetEditionResults(w http.ResponseWriter, r *http.Request) {
	matches, ok := h.loadEditionMatches(w, r, h.db.GetEditionResult)

	if ok {
		groupMatchesResponse := transform.GroupByRounds(matches)
		WriteJSON(w, http.StatusOK, groupMatchesResponse)
	}
}

func (h *Handler) GetEditionMatches(w http.ResponseWriter, r *http.Request) {
	matches, ok := h.loadEditionMatches(w, r, h.db.GetEditionMatches)

	if ok {
		groupMatchesResponse := transform.GroupByRounds(matches)
		WriteJSON(w, http.StatusOK, groupMatchesResponse)
	}
}

func (h *Handler) GetFeaturedMatches(w http.ResponseWriter, r *http.Request) {
	matches, ok := h.loadEditionMatches(w, r,
		h.db.FeaturedMatchesQuery(transform.From(), transform.To()),
	)

	if ok {
		featuredMatchesResponse := transform.GetFeaturedMatches(matches)
		WriteJSON(w, http.StatusOK, featuredMatchesResponse)
	}
}

func (h *Handler) loadEditionMatches(w http.ResponseWriter, r *http.Request,
	getEditionMatches func(context.Context, int, int) ([]database.Match, error)) (
	[]database.Match, bool) {

	ctx := r.Context()

	competitionID, startYear, err := editionParams(r)
	if err != nil {
		BadRequest(w)
		return nil, false
	}

	matches, err := getEditionMatches(ctx, competitionID, startYear)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load matches")
		return nil, false
	}

	return matches, true
}

func (h *Handler) GetClubPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clubId, err := clubParam(r)
	if err != nil {
		BadRequest(w)
		return
	}

	startYear := sync.CurrentSeasonStartYear(time.Now())

	competitionID, err := competitionParam(r)
	if err != nil {
		BadRequest(w)
		return
	}

	var players []database.SeasonPlayer
	if err := h.db.List(ctx, &players, database.Filter{
		"competition_id":    competitionID,
		"start_season_year": startYear,
		"team_id":           clubId,
	}, database.PreloadPlayer, "Player.NationalityArea"); err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load club players")
		return
	}

	playersResponse := transform.GroupByPositions(players)

	WriteJSON(w, http.StatusOK, playersResponse)
}

func (h *Handler) GetTeamInformation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clubID, err := clubParam(r)
	if err != nil {
		BadRequest(w)
		return
	}

	team, err := h.db.GetTeam(ctx, clubID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load club information")
		return
	}

	WriteJSON(w, http.StatusOK, team)
}

func (h *Handler) GetActiveEdition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var competitions []database.Competition
	if err := h.db.List(ctx, &competitions, nil, database.PreloadEditions); err != nil {
		http.Error(w, "failed to load leagues", http.StatusInternalServerError)
		return
	}

	response := transform.GetActiveEditions(competitions)

	WriteJSON(w, http.StatusOK, response)
}
func (h *Handler) GetCompetitionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	competitionID, startYear, err := editionParams(r)
	if err != nil {
		BadRequest(w)
		return
	}
	edition, err := h.db.GetEdition(ctx, competitionID, startYear, database.PreloadCompetition)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load editions")
		return
	}

	response := transform.CreateEditionResponse(edition.Competition, edition)

	WriteJSON(w, http.StatusOK, response)
}
func (h *Handler) GetCompetitionEditions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	competitionID, err := competitionParam(r)
	if err != nil {
		BadRequest(w)
		return
	}

	editions, err := h.db.GetEditions(ctx, competitionID, database.PreloadCompetition)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load editions")
		return
	}

	response := transform.GetEditions(editions)

	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetEditionGoalScorers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	competitionID, startYear, err := editionParams(r)
	if err != nil {
		BadRequest(w)
		return
	}

	players, err := h.db.GetEditionGoalScorers(ctx, competitionID, startYear)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load goal scorers")
		return
	}

	response := transform.GetEditionGoalScorers(players)

	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetEditionTable(w http.ResponseWriter, r *http.Request) {
	competitionID, startYear, err := editionParams(r)
	if err != nil {
		BadRequest(w)
		return
	}

	clubs, ok := h.store.GetTableClubs(competitionID, startYear)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "failed to load table")
		return
	}
	response := transform.GetFullInformationClubs(clubs, competitionID)

	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) GetClubTable(w http.ResponseWriter, r *http.Request) {
	clubID, competitionID, startYear, err := clubparams(r)
	if err != nil {
		BadRequest(w)
		return
	}

	clubs, ok := h.store.GetTableClubs(competitionID, startYear)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "failed to load table")
		return
	}
	response := transform.GetMiniTableForClub(clubs, competitionID, clubID)

	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) getClubMatches(w http.ResponseWriter, r *http.Request,
	getMatches func(
		ctx context.Context, competitionID int, startYear int,
	) ([]database.Match, error)) {
	ctx := r.Context()

	clubID, competitionID, startYear, err := clubparams(r)
	if err != nil {
		BadRequest(w)
		return
	}

	matches, err := getMatches(ctx, competitionID, startYear)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to load matches")
		return
	}

	teamMatchesResponse := transform.FiltTeamMatches(matches, clubID)
	WriteJSON(w, http.StatusOK, teamMatchesResponse)
}

func (h *Handler) GetClubMatches(w http.ResponseWriter, r *http.Request) {
	h.getClubMatches(w, r, h.db.GetEditionMatches)

}

func (h *Handler) GetClubResults(w http.ResponseWriter, r *http.Request) {
	h.getClubMatches(w, r, h.db.GetEditionResult)
}
