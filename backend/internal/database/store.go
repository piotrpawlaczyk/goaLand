package database

import (
	"context"
)

type StorageDatabase interface {
	List(ctx context.Context, dest any, filter Filter, preloads ...string) error
	Save(context.Context, any) error
	SaveAndCheck(context.Context, any) (bool, error)
	ClearSeason(ctx context.Context, competitonID, startYear int) error
	DeleteSeasonPlayers(ctx context.Context, competitonID, startYear int) error
	DeleteSeasonMatches(ctx context.Context, competitonID, startYear int) error
	DoTransaction(ctx context.Context, fn func(client StorageDatabase) error) error
}
