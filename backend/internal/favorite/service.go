package favorite

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(ctx context.Context, userID, productID uint64) error {
	return s.repo.Add(ctx, userID, productID)
}

func (s *Service) Remove(ctx context.Context, userID, productID uint64) error {
	return s.repo.Remove(ctx, userID, productID)
}

func (s *Service) IsFavorited(ctx context.Context, userID, productID uint64) (bool, error) {
	return s.repo.IsFavorited(ctx, userID, productID)
}

func (s *Service) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]FavoriteDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListByUser(ctx, userID, page, pageSize)
}
