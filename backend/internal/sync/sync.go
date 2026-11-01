package sync

import (
	"backend/internal/api"
	"backend/internal/database"
	"backend/internal/state"
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	stdSync "sync"
)

type SeasonTarget struct {
	CompetitionCode string
	StartYear       int
}

type SeasonTargets []SeasonTarget

type Season struct {
	CompetitionID   int
	CompetitionCode string
	StartYear       int
}
type Seasons []Season

// TODO: load seasons and areas from DB on start working if DB is not empty

//go:generate mockgen -destination=mocks/sync_mock.go -package=mocks . API,Database
type API interface {
	FetchAreas() ([]api.Area, error)
	FetchCompetition(code string) (api.Competition, error)
	FetchCompetitions() ([]api.Competition, error)

	FetchEdition(code string, startYear int) (api.Edition, error)
	FetchTeams(code string, startYear int) ([]api.Team, error)
	FetchMatches(code string, startYear int) ([]api.Match, error)
	FetchGoalScorers(code string, startYear, limit int) ([]api.GoalScorer, error)
}

type Database = database.StorageDatabase

type Sync struct {
	apiClient      API
	databaseClient Database
	seasons        Seasons
	areasByName    map[string]int
	storage        *state.Store

	cron *cron.Cron
}

type SeasonJob struct {
	mu     stdSync.Mutex
	season SeasonTarget
	sync   *Sync
}

func (s *SeasonJob) Run() {
	ctx := context.Background()
	now := time.Now()

	fmt.Printf("update competition %s", s.season.CompetitionCode)
	if startYear := CurrentSeasonStartYear(now); startYear == s.season.StartYear {
		err := s.sync.updateEdition(ctx, s.season)
		if err != nil {
			fmt.Print(err.Error())
		}
	} else {
		s.season.StartYear = startYear
		err := s.sync.AddSeason(ctx, now, s.season)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

var nationalityAliases = map[string]string{
	"Ireland":            "Republic of Ireland",
	"DR Congo":           "Congo DR",
	"Bosnia-Herzegovina": "Bosnia and Herzegovina",
	"Cote d'Ivoire":      "Ivory Coast",
	"North Macedonia":    "FYR Macedonia",
	"Korea, South":       "South Korea",
}

func New(apiClient API, databaseClient Database, storage *state.Store) *Sync {
	return &Sync{
		apiClient:      apiClient,
		databaseClient: databaseClient,
		areasByName:    make(map[string]int),
		storage:        storage,
		cron: cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		)),
	}
}

func (s *Sync) replaceDatabase(databaseClient Database) *Sync {
	return &Sync{
		apiClient:      s.apiClient,
		databaseClient: databaseClient,
		seasons:        s.seasons,
		areasByName:    s.areasByName,
		storage:        s.storage,
	}
}

func (s *Sync) Start() {
	s.cron.Start()
}

func (s *Sync) getAreaByNameWithError(nationality string, playerId int) (int, error) {
	if _, exist := nationalityAliases[nationality]; exist {
		nationality = nationalityAliases[nationality]
	}

	areasID, exist := s.areasByName[nationality]
	if !exist {
		return 0, fmt.Errorf("Not found Area %s, players %d", nationality, playerId)
	}
	return areasID, nil
}
func (s *Sync) getAreaByName(nationality string) int {
	areasID, exist := s.areasByName[nationality]
	if !exist {
		return UnknownAreaID
	}
	return areasID
}

const UnknownAreaID = -1

var unknownArea = database.Area{
	AreaID:    UnknownAreaID,
	Name:      "Unknown",
	Code:      "UNK",
	IsCountry: false,
}

func NewFromEnv() (*Sync, error) {
	var databaseClient Database
	databaseClient, err := database.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	return &Sync{
		apiClient:      api.NewClientFromEnv(),
		databaseClient: databaseClient,
		seasons:        nil,
		areasByName:    make(map[string]int),
	}, nil
}

func (targets SeasonTargets) competitionCodes() map[string]struct{} {
	set := make(map[string]struct{})

	for _, target := range targets {
		set[target.CompetitionCode] = struct{}{}
	}

	return set
}

func (seasons Seasons) CompetitionID(code string) (int, bool) {
	for _, season := range seasons {
		if season.CompetitionCode == code {
			return season.CompetitionID, true
		}
	}

	return 0, false
}

func CurrentSeasonStartYear(now time.Time) int {
	if now.Month() <= time.June {
		return now.Year() - 1
	}
	return now.Year()
}

func isCurrentSeason(startYear int, now time.Time) bool {
	return startYear == CurrentSeasonStartYear(now)
}

func createExpression(number int) []string {
	var expressions []string

	for i := range number {
		expressions = append(expressions, fmt.Sprintf("%d-59/5 * * * *", i))
	}
	return expressions
}

func (s *Sync) InitAreasFromDB(ctx context.Context) error {
	var areas []database.Area
	err := s.databaseClient.List(ctx, &areas, database.Filter{})
	if err != nil {
		return err
	}

	for _, area := range areas {
		s.areasByName[area.Name] = area.AreaID
		s.areasByName[area.Code] = area.AreaID

	}

	var editions []database.Edition
	err = s.databaseClient.List(ctx, &editions, database.Filter{}, database.PreloadCompetition)
	if err != nil {
		return err
	}

	for _, edition := range editions {
		s.seasons = append(s.seasons, Season{
			CompetitionID:   edition.CompetitionID,
			CompetitionCode: edition.Competition.Code,
			StartYear:       edition.StartYear,
		})
	}
	return nil
}
