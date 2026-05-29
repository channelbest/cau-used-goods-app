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

func (s *Service) OrderOverview(ctx context.Context) (*OrderOverview, error) {
	return s.repo.OrderOverview(ctx)
}

func (s *Service) UserOverview(ctx context.Context) (*UserOverview, error) {
	return s.repo.UserOverview(ctx)
}

func (s *Service) ReportOverview(ctx context.Context) (*ReportOverview, error) {
	return s.repo.ReportOverview(ctx)
}
