package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cau-used-goods-app/backend/internal/config"
	jwtutil "cau-used-goods-app/backend/pkg/jwt"
)

type Service struct {
	repo   *Repository
	jwt    config.JWTConfig
	wechat config.WechatConfig
}

type LoginResult struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type DevLoginInput struct {
	OpenID string
	Role   string
}

func NewService(repo *Repository, jwtCfg config.JWTConfig, wechatCfg config.WechatConfig) *Service {
	return &Service{repo: repo, jwt: jwtCfg, wechat: wechatCfg}
}

func (s *Service) DevLogin(ctx context.Context, input DevLoginInput) (*LoginResult, error) {
	if input.OpenID == "" {
		input.OpenID = "dev_openid_001"
	}
	if input.Role != "" && input.Role != "USER" && input.Role != "ADMIN" {
		return nil, fmt.Errorf("role must be USER or ADMIN")
	}

	result, err := s.loginByOpenID(ctx, input.OpenID)
	if err != nil {
		return nil, err
	}
	if input.Role != "" && result.User.Role != input.Role {
		if err := s.repo.UpdateRole(ctx, result.User.ID, input.Role); err != nil {
			return nil, err
		}
		result.User.Role = input.Role
		result.Token, err = jwtutil.Generate(s.jwt.Secret, s.jwt.ExpireHours, result.User.ID, result.User.Role)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Service) WechatLogin(ctx context.Context, code string) (*LoginResult, error) {
	if code == "" {
		return nil, fmt.Errorf("code is required")
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
	} else if err := s.repo.UpdateLastLoginTime(ctx, user.ID); err != nil {
		return nil, err
	}

	token, err := jwtutil.Generate(s.jwt.Secret, s.jwt.ExpireHours, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{Token: token, User: user}, nil
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
		return "", fmt.Errorf("wechat app_id/app_secret is not configured")
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
		return "", fmt.Errorf("wechat code2session failed: %d %s", result.ErrCode, result.ErrMsg)
	}
	if result.OpenID == "" {
		return "", fmt.Errorf("wechat response missing openid")
	}
	return result.OpenID, nil
}
