package model

import (
	"context"

	"github.com/TestRpc/model/result"
)

type Model interface {
	Save(ctx context.Context) result.DbResult
	Update(ctx context.Context) result.DbResult
	Get(ctx context.Context) result.DbResult
}
