package service

import (
	"context"
	"fmt"
	"menghitung-tekanan/internal/domain"
)

type TekananRepository interface {
	Save(ctx context.Context, data domain.Data) error
	GetAll(ctx context.Context) ([]domain.Data, error)
	GetLatest(ctx context.Context) (domain.Data, error)
}

type TekananService struct {
	repo TekananRepository
}

func NewTekananService(repo TekananRepository) *TekananService {
	return &TekananService{
		repo: repo,
	}
}

// Menyimpan data tekanan
func (s *TekananService) Save(ctx context.Context, data domain.Data) error {
	err := s.repo.Save(ctx, data)
	if err != nil {
		return fmt.Errorf("service save tekanan: %w", err)
	}

	return nil
}

// Mengambil seluruh history tekanan
func (s *TekananService) GetAll(ctx context.Context) ([]domain.Data, error) {
	result, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service get history tekanan: %w", err)
	}

	return result, nil
}

// Mengambil tekanan terbaru
func (s *TekananService) GetLatest(ctx context.Context) (domain.Data, error) {
	result, err := s.repo.GetLatest(ctx)
	if err != nil {
		return domain.Data{}, fmt.Errorf("service get latest tekanan: %w", err)
	}

	return result, nil
}
