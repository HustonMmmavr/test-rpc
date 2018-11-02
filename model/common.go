package model

import (
	"context"

	"github.com/jackc/pgx"
)

func getDb(ctx context.Context) *pgx.ConnPool {
	switch ctx.Value("db").(type) {
	case *pgx.ConnPool:
		return ctx.Value("db").(*pgx.ConnPool)
	default:
		return nil
	}
}
