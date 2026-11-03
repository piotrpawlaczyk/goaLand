package server

import (
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/state"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(db *database.Client, store *state.Store) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := handler.NewHandler(db, store)

	r.Get("/health", handler.GetHealth)

	r.Get("/competitions", handler.GetActiveEdition)
	r.Get("/competitions/{competitionID}/editions", handler.GetCompetitionEditions)
	r.Get("/competitions/{competitionID}/{startYear}/data", handler.GetCompetitionData)

	r.Get("/competitions/{competitionID}/{startYear}/matches", handler.GetEditionMatches)
	r.Get("/competitions/{competitionID}/{startYear}/goal-scorers", handler.GetEditionGoalScorers)
	r.Get("/competitions/{competitionID}/{startYear}/results", handler.GetEditionResults)
	r.Get("/competitions/{competitionID}/{startYear}/featured-matches", handler.GetFeaturedMatches)

	r.Get("/all-matches", handler.GetAllMatches)
	r.Get("/competitions/{competitionID}/{startYear}/table", handler.GetEditionTable)
	r.Get("/competitions/{competitionID}/clubs/{clubID}/players", handler.GetClubPlayers)
	r.Get("/competitions/{competitionID}/clubs/{clubID}/information", handler.GetTeamInformation)

	r.Get("/competitions/{competitionID}/clubs/{clubID}/matches", handler.GetClubMatches)
	r.Get("/competitions/{competitionID}/clubs/{clubID}/results", handler.GetClubResults)
	r.Get("/competitions/{competitionID}/clubs/{clubID}/table", handler.GetClubTable)

	return r
}
