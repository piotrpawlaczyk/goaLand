package main

import (
	"backend/internal/api"
	"backend/internal/state"
	"backend/internal/sync"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"backend/internal/database"
	"backend/internal/server"
)

var trackedCompetitionCodes = []string{"PL", "PD", "BL1", "SA", "FL1"}

func main() {
	initData := flag.Bool("init-data", false, "Create missing tables and fetch initial data, then exit")
	track := flag.Bool("track", true, "Update competition data periodically")
	flag.Parse()

	ctx := context.Background()

	db, err := database.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	store := state.NewStore(ctx, db)

	apiClient := api.NewClientFromEnv()

	worker := sync.New(apiClient, db, store)
	if *initData {
		if err := db.CreateDatabase(); err != nil {
			log.Fatal(err)
		}
		if err := fetchData(ctx, worker); err != nil {
			log.Fatal(err)
		}
		log.Println("Data initialization complete")
		return
	}

	if err := worker.InitAreasFromDB(ctx); err != nil {
		log.Fatalf("failed to load areas and seasons from database: %v", err)
	}

	if err = store.InitTables(ctx); err != nil {
		log.Fatal(err)
	}

	if *track {
		if err := worker.TrackCompetitions(trackedCompetitionCodes); err != nil {
			log.Fatal(err)
		}
	}

	server := server.NewServer(db, store)

	log.Println("Starting server on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func fetchData(ctx context.Context, worker *sync.Sync) error {
	now := time.Now()

	oldSeasonTargets := []sync.SeasonTarget{{
		CompetitionCode: "PL",
		StartYear:       2025,
	}}

	err := worker.InitializeData(ctx, oldSeasonTargets)
	if err != nil {
		return err
	}

	for _, code := range trackedCompetitionCodes {
		target := sync.SeasonTarget{
			CompetitionCode: code,
			StartYear:       sync.CurrentSeasonStartYear(now),
		}
		if err := worker.AddSeason(ctx, now, target); err != nil {
			return fmt.Errorf("failed to load competition %s for season %d: %w", code, target.StartYear, err)
		}
	}
	return nil
}
