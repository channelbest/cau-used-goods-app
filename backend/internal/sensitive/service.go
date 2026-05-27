package sensitive

import (
	"context"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CheckResult struct {
	Passed   bool     `json:"passed"`
	HitWords []string `json:"hitWords"`
	Message  string   `json:"message"`
}

func (s *Service) CheckText(ctx context.Context, text string) (*CheckResult, error) {
	words, err := s.repo.ListEnabledWords(ctx)
	if err != nil {
		return nil, err
	}

	text = strings.ToLower(text)
	var hits []string

	for _, w := range words {
		word := strings.ToLower(strings.TrimSpace(w.Word))
		if word == "" {
			continue
		}
		if strings.Contains(text, word) {
			hits = append(hits, w.Word)
		}
	}

	if len(hits) > 0 {
		return &CheckResult{
			Passed:   false,
			HitWords: hits,
			Message:  "内容包含敏感词，请修改后重试",
		}, nil
	}

	return &CheckResult{
		Passed:   true,
		HitWords: []string{},
		Message:  "内容检测通过",
	}, nil
}
