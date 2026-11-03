package transform

import (
	"backend/internal/database"
	"backend/internal/utils"
)

type EditionResponse struct {
	CompetitionID   int    `json:"id"`
	CompetitionName string `json:"competitionName"`
	Code            string `json:"code"`

	StartYear int `json:"startYear"`

	IsCurrent bool `json:"isCurrent"`
}

func CreateEditionResponse(competition database.Competition, edition database.Edition) EditionResponse {
	return EditionResponse{
		CompetitionID: competition.CompetitionID,

		CompetitionName: competition.Name,
		Code:            competition.Code,

		StartYear: edition.StartYear,
		IsCurrent: utils.IsCurrent(edition.Status),
	}
}

func GetEditions(edition []database.Edition) []EditionResponse {
	var response []EditionResponse
	for _, e := range edition {
		response = append(response, CreateEditionResponse(e.Competition, e))
	}
	return response
}

func getAcitveOrUpcoming(competition database.Competition) (EditionResponse, bool) {
	for _, edition := range competition.Editions {
		if utils.IsCurrent(edition.Status) {
			return CreateEditionResponse(competition, edition), true
		}
	}

	return EditionResponse{}, false
}

func GetActiveEditions(competitions []database.Competition) []EditionResponse {
	response := []EditionResponse{}

	for _, competition := range competitions {
		if competitionEdition, ok := getAcitveOrUpcoming(competition); ok {
			response = append(response, competitionEdition)
		}
	}

	return response
}
