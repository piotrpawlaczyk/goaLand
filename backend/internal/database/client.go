package database

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	PreloadArea        = "Area"
	PreloadEditions    = "Editions"
	PreloadMatches     = "Matches"
	PreloadCompetition = "Competition"
	PreloadPlayer      = "Player"
)

var preloadTeams = []string{"HomeTeam", "AwayTeam"}

type Client struct {
	db *gorm.DB
}

type Filter = map[string]any

func NewClient(config Config) (*Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Name, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func NewClientFromEnv() (*Client, error) {
	config, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return NewClient(config)
}
func (c *Client) Save(ctx context.Context, data any) error {
	return c.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(data).
		Error
}

func (c *Client) SaveAndCheck(ctx context.Context, data any) (bool, error) {
	result := c.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(data)

	return result.RowsAffected > 0, result.Error
}

func (c *Client) buildQuery(ctx context.Context, preloads []string) *gorm.DB {
	query := c.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query

}

func (c *Client) List(ctx context.Context, dest any, filter Filter, preloads ...string) error {
	return c.buildQuery(ctx, preloads).Where(filter).Find(dest).Error
}

func (c *Client) Get(ctx context.Context, dest any, filter Filter, preloads ...string) error {
	return c.buildQuery(ctx, preloads).Where(filter).First(dest).Error
}

func (c *Client) ListOr(ctx context.Context, dest any, filters []Filter, preloads ...string) error {
	if len(filters) == 0 {
		return errors.New("no filters provided")
	}

	query := c.buildQuery(ctx, preloads).Where(filters[0])

	for _, filter := range filters[1:] {
		query = query.Or(filter)
	}

	return query.Find(dest).Error
}

func (c *Client) DB() *gorm.DB {
	return c.db
}

func (c *Client) GetArea(ctx context.Context, id int, preloads ...string) (*Area, error) {
	var area Area
	err := c.Get(ctx, &area, Filter{"area_id": id}, preloads...)
	return &area, err
}

func (c *Client) GetCompetition(ctx context.Context, id int, preloads ...string) (*Competition, error) {
	var competition Competition
	err := c.Get(ctx, &competition, Filter{"competition_id": id}, preloads...)
	return &competition, err
}

func (c *Client) GetTeam(ctx context.Context, id int, preloads ...string) (*Team, error) {
	var team Team
	err := c.Get(ctx, &team, Filter{"team_id": id}, preloads...)
	return &team, err
}

func (c *Client) GetMatch(ctx context.Context, id int, preloads ...string) (*Match, error) {
	var match Match
	err := c.Get(ctx, &match, Filter{"match_id": id}, preloads...)
	return &match, err
}

func (c *Client) GetPlayer(ctx context.Context, id int, preloads ...string) (*Player, error) {
	var scorer Player
	err := c.Get(ctx, &scorer, Filter{"goal_scorer_id": id}, preloads...)
	return &scorer, err
}

func (c *Client) GetSeasonPlayer(ctx context.Context, id int, preloads ...string) (*SeasonPlayer, error) {
	var scorer SeasonPlayer
	err := c.Get(ctx, &scorer, Filter{"goal_scorer_id": id}, preloads...)
	return &scorer, err
}

func (c *Client) DoTransaction(ctx context.Context, fn func(client StorageDatabase) error) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		client := &Client{db: tx}
		return fn(client)
	})
}
