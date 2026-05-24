package user

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo *Repository
}

type UpdateProfileInput struct {
	Nickname  *string
	AvatarURL *string
	Phone     *string
}

type SubmitStudentVerificationInput struct {
	StudentID string
	RealName  string
	College   string
}

type ReviewStudentVerificationInput struct {
	UserID      uint64
	AuthStatus  string
	Description string
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Me(ctx context.Context, userID uint64) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uint64, input UpdateProfileInput) (*User, error) {
	if input.Nickname == nil && input.AvatarURL == nil && input.Phone == nil {
		return nil, fmt.Errorf("nothing to update")
	}

	if input.Nickname != nil {
		trimmed := strings.TrimSpace(*input.Nickname)
		input.Nickname = &trimmed
	}
	if input.AvatarURL != nil {
		trimmed := strings.TrimSpace(*input.AvatarURL)
		input.AvatarURL = &trimmed
	}
	if input.Phone != nil {
		trimmed := strings.TrimSpace(*input.Phone)
		input.Phone = &trimmed
	}

	if err := s.repo.UpdateProfile(ctx, userID, input.Nickname, input.AvatarURL, input.Phone); err != nil {
		return nil, err
	}
	return s.Me(ctx, userID)
}

func (s *Service) SubmitStudentVerification(ctx context.Context, userID uint64, input SubmitStudentVerificationInput) (*StudentVerification, error) {
	input.StudentID = strings.TrimSpace(input.StudentID)
	input.RealName = strings.TrimSpace(input.RealName)
	input.College = strings.TrimSpace(input.College)

	if input.StudentID == "" {
		return nil, fmt.Errorf("studentId is required")
	}
	if input.RealName == "" {
		return nil, fmt.Errorf("realName is required")
	}
	if input.College == "" {
		return nil, fmt.Errorf("college is required")
	}

	current, err := s.repo.FindStudentVerification(ctx, userID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("user not found")
	}
	if current.AccountStatus != "NORMAL" {
		return nil, fmt.Errorf("account status is not normal")
	}
	if current.AuthStatus == "VERIFIED" {
		return nil, fmt.Errorf("student verification already approved")
	}

	if err := s.repo.SubmitStudentVerification(ctx, userID, input.StudentID, input.RealName, input.College); err != nil {
		return nil, err
	}
	return s.StudentVerification(ctx, userID)
}

func (s *Service) StudentVerification(ctx context.Context, userID uint64) (*StudentVerification, error) {
	verification, err := s.repo.FindStudentVerification(ctx, userID)
	if err != nil {
		return nil, err
	}
	if verification == nil {
		return nil, fmt.Errorf("user not found")
	}
	return verification, nil
}

func (s *Service) ListStudentVerifications(ctx context.Context, status string) ([]StudentVerification, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		status = "PENDING"
	}
	switch status {
	case "PENDING", "VERIFIED", "REJECTED", "UNVERIFIED":
	default:
		return nil, fmt.Errorf("invalid authStatus")
	}
	return s.repo.ListStudentVerifications(ctx, status)
}

func (s *Service) ReviewStudentVerification(ctx context.Context, adminID uint64, input ReviewStudentVerificationInput) (*StudentVerification, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("userId is required")
	}
	input.AuthStatus = strings.TrimSpace(input.AuthStatus)
	input.Description = strings.TrimSpace(input.Description)

	if input.AuthStatus != "VERIFIED" && input.AuthStatus != "REJECTED" {
		return nil, fmt.Errorf("authStatus must be VERIFIED or REJECTED")
	}
	if input.Description == "" {
		if input.AuthStatus == "VERIFIED" {
			input.Description = "学生认证审核通过"
		} else {
			input.Description = "学生认证审核驳回"
		}
	}

	if err := s.repo.ReviewStudentVerification(ctx, adminID, input.UserID, input.AuthStatus, input.Description); err != nil {
		return nil, err
	}
	return s.StudentVerification(ctx, input.UserID)
}
