package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cau-used-goods-app/backend/internal/config"
	jwtutil "cau-used-goods-app/backend/pkg/jwt"
)

type Service struct {
	repo        *Repository
	jwt         config.JWTConfig
	wechat      config.WechatConfig
	reactivator AccountReactivator
}

type LoginResult struct {
	Token                string `json:"token,omitempty"`
	ReactivationToken    string `json:"reactivationToken,omitempty"`
	RequiresReactivation bool   `json:"requiresReactivation,omitempty"`
	User                 *User  `json:"user"`
}

type AccountReactivator interface {
	ReactivateAccount(ctx context.Context, userID uint64) error
}

type DevLoginInput struct {
	OpenID string
	Role   string
}

func NewService(repo *Repository, jwtCfg config.JWTConfig, wechatCfg config.WechatConfig, reactivator AccountReactivator) *Service {
	return &Service{repo: repo, jwt: jwtCfg, wechat: wechatCfg, reactivator: reactivator}
}

func (s *Service) DevLogin(ctx context.Context, input DevLoginInput) (*LoginResult, error) {
	input.OpenID = strings.TrimSpace(input.OpenID)
	input.Role = strings.TrimSpace(input.Role)
	if input.OpenID == "" {
		input.OpenID = "dev_openid_001"
	}
	if len(input.OpenID) > 64 {
		return nil, fmt.Errorf("openid length cannot exceed 64 characters")
	}
	if input.Role != "" && input.Role != "USER" && input.Role != "ADMIN" && input.Role != "SUPER_ADMIN" {
		return nil, fmt.Errorf("role can only be USER, ADMIN or SUPER_ADMIN")
	}

	result, err := s.loginByOpenID(ctx, input.OpenID)
	if err != nil {
		return nil, err
	}
	if result.RequiresReactivation {
		return result, nil
	}
	if input.Role != "" && result.User.Role != input.Role {
		if err := s.repo.UpdateRole(ctx, result.User.ID, input.Role); err != nil {
			return nil, err
		}
		result.User.Role = input.Role
		result.Token, err = jwtutil.Generate(s.jwt.Secret, s.jwt.ExpireHours, result.User.ID, result.User.Role, result.User.TokenVersion)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Service) WechatLogin(ctx context.Context, code string) (*LoginResult, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, fmt.Errorf("wechat login code cannot be empty")
	}
	if len(code) > 128 {
		return nil, fmt.Errorf("code length cannot exceed 128 characters")
	}

	openid, err := s.fetchWechatOpenID(ctx, code)
	if err != nil {
		return nil, err
	}
	return s.loginByOpenID(ctx, openid)
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

func (s *Service) loginByOpenID(ctx context.Context, openid string) (*LoginResult, error) {
	user, err := s.repo.FindByOpenID(ctx, openid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = s.repo.CreateByOpenID(ctx, openid)
		if err != nil {
			return nil, err
		}
	} else if user.AccountStatus == "CANCELED" || user.IsDeleted {
		token, err := jwtutil.GenerateReactivation(s.jwt.Secret, user.ID, user.TokenVersion)
		if err != nil {
			return nil, err
		}
		return &LoginResult{
			ReactivationToken:    token,
			RequiresReactivation: true,
			User:                 user,
		}, nil
	} else if err := s.repo.UpdateLastLoginTime(ctx, user.ID); err != nil {
		return nil, err
	}

	token, err := jwtutil.Generate(s.jwt.Secret, s.jwt.ExpireHours, user.ID, user.Role, user.TokenVersion)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, User: user}, nil
}

func (s *Service) Reactivate(ctx context.Context, token string, confirmed bool) (*LoginResult, error) {
	if !confirmed {
		return nil, fmt.Errorf("reactivation confirmation is required")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("reactivationToken is required")
	}
	claims, err := jwtutil.Parse(s.jwt.Secret, token)
	if err != nil || claims.TokenType != jwtutil.TokenTypeReactivate {
		return nil, fmt.Errorf("invalid or expired reactivation token")
	}
	current, err := s.repo.FindByIDIncludingDeleted(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if current == nil || current.TokenVersion != claims.TokenVersion {
		return nil, fmt.Errorf("invalid or expired reactivation token")
	}
	if s.reactivator == nil {
		return nil, fmt.Errorf("reactivation service is unavailable")
	}
	if err := s.reactivator.ReactivateAccount(ctx, claims.UserID); err != nil {
		return nil, err
	}
	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	accessToken, err := jwtutil.Generate(s.jwt.Secret, s.jwt.ExpireHours, user.ID, user.Role, user.TokenVersion)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: accessToken, User: user}, nil
}

type wechatSessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (s *Service) fetchWechatOpenID(ctx context.Context, code string) (string, error) {
	if s.wechat.AppID == "" || s.wechat.AppSecret == "" {
		return "", fmt.Errorf("wechat appid or secret is not configured")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weixin.qq.com/sns/jscode2session", nil)
	if err != nil {
		return "", fmt.Errorf("create wechat request: %w", err)
	}

	q := req.URL.Query()
	q.Set("appid", s.wechat.AppID)
	q.Set("secret", s.wechat.AppSecret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request wechat code2session: %w", err)
	}
	defer resp.Body.Close()

	var result wechatSessionResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode wechat response: %w", err)
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat login failed: %s", result.ErrMsg)
	}
	if result.OpenID == "" {
		return "", fmt.Errorf("wechat openid is empty")
	}
	return result.OpenID, nil
}
