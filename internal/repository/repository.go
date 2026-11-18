package repository

import (
	"context"
	"menghitung-tekanan/internal/domain"
)

type Tekanan interface {
	Save(ctx context.Context, data domain.Data) error
	GetAll(ctx context.Context) ([]domain.Data, error)
	GetLatest(ctx context.Context) (domain.Data, error)
}
