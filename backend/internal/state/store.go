package state

import (
	"backend/internal/database"
	"backend/internal/state/table"
	"context"
	"sync"
)

type tableKey struct {
	competitionID int
	startYear     int
}

type Store struct {
	tables map[tableKey]*table.Table
	db     *database.Client
	mu     sync.RWMutex
}

func NewStore(ctx context.Context, db *database.Client) *Store {
	return &Store{
		tables: make(map[tableKey]*table.Table),
		db:     db,
	}
}

func (c *Store) InitTables(ctx context.Context) error {
	var editions []database.Edition

	if err := c.db.List(ctx, &editions, database.Filter{}); err != nil {
		return err
	}

	for _, edition := range editions {
		table, err := table.CreateTable(ctx, c.db, edition.CompetitionID, edition.StartYear)
		if err != nil {
			return err
		}
		c.tables[GetKey(table)] = table
	}
	return nil
}
func (s *Store) RefreshTable(ctx context.Context, competitionID, startYear int) error {
	key := tableKey{
		competitionID: competitionID,
		startYear:     startYear,
	}
	newTable, err := table.CreateTable(ctx, s.db, competitionID, startYear)

	if err != nil {
		return err
	}
	s.mu.Lock()
	s.tables[key] = newTable
	s.mu.Unlock()
	return nil
}

func GetKey(table *table.Table) tableKey {
	return tableKey{
		competitionID: table.CompetitionID,
		startYear:     table.StartYear,
	}
}

func (s *Store) GetTable(competitionID int, startYear int) (*table.Table, bool) {
	key := tableKey{
		competitionID: competitionID,
		startYear:     startYear,
	}

	s.mu.RLock()
	table, exists := s.tables[key]
	s.mu.RUnlock()

	if !exists {
		return nil, false
	}
	return table, true
}

func (s *Store) GetTableClubs(competitionID int, startYear int) ([]*table.Club, bool) {
	table, ok := s.GetTable(competitionID, startYear)
	if ok {
		return table.Clubs, ok
	}
	return nil, ok
}
