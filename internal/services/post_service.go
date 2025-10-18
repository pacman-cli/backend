// Package services
package services

import (
	"context"
	"errors"
	"strings"

	"github.com/puspo/basicnewproject/internal/models"
	"github.com/puspo/basicnewproject/internal/repositories"
)

// PostService holds business logic and validation for posts.
type PostService struct {
	repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// Validate basic invariants for a post.
func (s *PostService) validate(p *models.Post) error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("content is required")
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}
	if p.Metadata == nil {
		p.Metadata = map[string]any{}
	}
	return nil
}

func (s *PostService) Create(ctx context.Context, p *models.Post) (int64, error) {
	if err := s.validate(p); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, p)
}

func (s *PostService) Get(ctx context.Context, id int64) (*models.Post, error) {
	return s.repo.Get(ctx, id)
}

func (s *PostService) List(ctx context.Context, offset, limit int) ([]models.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, offset, limit)
}

func (s *PostService) Update(ctx context.Context, id int64, p *models.Post) (bool, error) {
	if err := s.validate(p); err != nil {
		return false, err
	}
	return s.repo.Update(ctx, id, p)
}

func (s *PostService) Delete(ctx context.Context, id int64) (bool, error) {
	return s.repo.Delete(ctx, id)
}
