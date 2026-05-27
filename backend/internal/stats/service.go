package stats

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ProductOverview(ctx context.Context) (*ProductOverview, error) {
	return s.repo.ProductOverview(ctx)
}

func (s *Service) CategoryDistribution(ctx context.Context) ([]CategoryDistributionItem, error) {
	return s.repo.CategoryDistribution(ctx)
}

func (s *Service) StatusDistribution(ctx context.Context) ([]StatusDistributionItem, error) {
	return s.repo.StatusDistribution(ctx)
}

func (s *Service) ProductTrend(ctx context.Context, days int) ([]ProductTrendItem, error) {
	return s.repo.ProductTrend(ctx, days)
}
