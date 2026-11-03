package transform

import (
	"backend/internal/database"
	"backend/internal/utils"
	"cmp"
	"slices"
)

type GoalScorerResponse struct {
	PlayerID   int    `json:"playerID"`
	PlayerName string `json:"playerName"`
	TeamName   string `json:"teamName"`

	Position int `json:"position"`
	Goals    int `json:"goals"`
	Assists  int `json:"assists"`
}

func createGoalScorerResponse(player database.SeasonPlayer, position int) GoalScorerResponse {
	return GoalScorerResponse{
		PlayerID:   player.PlayerID,
		PlayerName: player.Player.Name,
		TeamName:   player.Team.FullName,
		Position:   position,
		Goals:      *player.Goals,
		Assists:    utils.IntOrZero(player.Assists),
	}
}
func GetEditionGoalScorers(players []database.SeasonPlayer) []GoalScorerResponse {
	goalScorers := []GoalScorerResponse{}

	if len(players) == 0 {
		return goalScorers
	}

	slices.SortFunc(players, func(a, b database.SeasonPlayer) int {
		if result := cmp.Compare(*b.Goals, *a.Goals); result != 0 {
			return result
		}
		return cmp.Compare(a.PlayerID, b.PlayerID)
	})

	position := 1

	goalScorers = append(goalScorers, createGoalScorerResponse(players[0], position))

	for i := range len(players) - 1 {
		if *players[i].Goals != *players[i+1].Goals {
			position += 1
		}
		goalScorers = append(goalScorers, createGoalScorerResponse(players[i+1], position))

	}

	return goalScorers

}
