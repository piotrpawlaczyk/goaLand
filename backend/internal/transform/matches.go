package transform

import (
	"backend/internal/database"
	"backend/internal/utils"
	"slices"
)

type ClubResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	CrestUrl string `json:"crestUrl"`
}
type MatchResponse struct {
	MatchID       int `json:"id"`
	CompetitionID int `json:"competitionID"`

	StartDate string `json:"date"`
	StartTime string `json:"time"`

	Status string `json:"status"`

	HomeTeam ClubResponse `json:"homeTeam"`
	AwayTeam ClubResponse `json:"awayTeam"`

	HomeScore *int `json:"homeScore,omitempty"`
	AwayScore *int `json:"awayScore,omitempty"`
}

type RoundResponse struct {
	Stage    string `json:"stage"`
	Matchday *int   `json:"matchday"`
}

type RoundMatchesResponse struct {
	Round   RoundResponse   `json:"round"`
	Matches []MatchResponse `json:"matches"`
}

type roundKey struct {
	stage    string
	matchday int
}

func matchStatus(status string) string {
	if slices.Contains(utils.UpcomingMatchStatuses, status) {
		return "scheduled"
	} else if slices.Contains(utils.LiveMatchStatuses, status) {
		return "live"
	} else if slices.Contains(utils.FinishedMatchStatuses, status) {
		return "finished"
	} else {
		return status
	}
}

func createMatchResponse(match database.Match) MatchResponse {
	return MatchResponse{
		MatchID:       match.MatchID,
		CompetitionID: match.CompetitionID,

		StartTime: match.StartTime.Format("15:04"),
		StartDate: match.StartTime.Format("02.01.2006"),

		Status: matchStatus(match.Status),

		HomeTeam: ClubResponse{
			ID:       match.HomeTeamID,
			Name:     match.HomeTeam.FullName,
			Code:     match.HomeTeam.Code,
			CrestUrl: match.HomeTeam.CrestUrl,
		},
		AwayTeam: ClubResponse{
			ID:       match.AwayTeamID,
			Name:     match.AwayTeam.FullName,
			Code:     match.AwayTeam.Code,
			CrestUrl: match.AwayTeam.CrestUrl,
		},

		HomeScore: match.HomeGoals,
		AwayScore: match.AwayGoals,
	}
}

func GetResponseMatchesWithLimit(matches []database.Match, limit int, firsts bool) []MatchResponse {
	if len(matches) > limit {
		if firsts {
			matches = matches[:limit]
		} else {
			matches = matches[len(matches)-limit:]
		}
	}

	return GetResponseMatches(matches)
}

func GetResponseMatches(matches []database.Match) []MatchResponse {
	response := make([]MatchResponse, 0, len(matches))
	for _, match := range matches {
		response = append(response, createMatchResponse(match))
	}
	return response
}

func GroupByRounds(matches []database.Match) []RoundMatchesResponse {
	groupedMatches := make([]RoundMatchesResponse, 0)
	var lastRound roundKey

	for _, match := range matches {
		matchday := 0
		if match.Matchday != nil {
			matchday = *match.Matchday
		}

		round := roundKey{
			stage:    match.Stage,
			matchday: matchday,
		}
		response := createMatchResponse(match)

		if len(groupedMatches) == 0 || lastRound != round {
			groupedMatches = append(groupedMatches, RoundMatchesResponse{
				Round: RoundResponse{
					Stage:    match.Stage,
					Matchday: match.Matchday,
				},
				Matches: []MatchResponse{response},
			})
			lastRound = round
			continue
		}

		lastGroup := &groupedMatches[len(groupedMatches)-1]
		lastGroup.Matches = append(lastGroup.Matches, response)
	}

	return groupedMatches
}

func FiltTeamMatches(matches []database.Match, clubID int) []MatchResponse {
	response := make([]MatchResponse, 0, len(matches))
	for _, match := range matches {
		if match.HomeTeamID == clubID || match.AwayTeamID == clubID {
			response = append(response, createMatchResponse(match))
		}
	}
	return response
}
