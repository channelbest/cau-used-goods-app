package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID            uint64     `json:"id"`
	OpenID        string     `json:"-"`
	Nickname      *string    `json:"nickname"`
	AvatarURL     *string    `json:"avatarUrl"`
	Role          string     `json:"role"`
	AuthStatus    string     `json:"authStatus"`
	AccountStatus string     `json:"accountStatus"`
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByOpenID(ctx context.Context, openid string) (*User, error) {
	const query = `
SELECT id, openid, nickname, avatar_url, role, auth_status, account_status, last_login_time
FROM users
WHERE openid = ? AND is_deleted = 0
LIMIT 1`

	var user User
	if err := r.db.QueryRowContext(ctx, query, openid).Scan(
		&user.ID,
		&user.OpenID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Role,
		&user.AuthStatus,
		&user.AccountStatus,
		&user.LastLoginTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by openid: %w", err)
	}
	return &user, nil
}

func (r *Repository) CreateByOpenID(ctx context.Context, openid string) (*User, error) {
	const execSQL = `
INSERT INTO users (openid, role, auth_status, account_status, last_login_time, create_time, update_time, is_deleted)
VALUES (?, 'USER', 'UNVERIFIED', 'NORMAL', NOW(), NOW(), NOW(), 0)`

	result, err := r.db.ExecContext(ctx, execSQL, openid)
	if err != nil {
		return nil, fmt.Errorf("create user by openid: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get new user id: %w", err)
	}

	return &User{
		ID:            uint64(id),
		OpenID:        openid,
		Role:          "USER",
		AuthStatus:    "UNVERIFIED",
		AccountStatus: "NORMAL",
		LastLoginTime: ptrTime(time.Now()),
	}, nil
}

func (r *Repository) UpdateLastLoginTime(ctx context.Context, userID uint64) error {
	const execSQL = `UPDATE users SET last_login_time = NOW(), update_time = NOW() WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, execSQL, userID); err != nil {
		return fmt.Errorf("update last login time: %w", err)
	}
	return nil
}

func (r *Repository) UpdateRole(ctx context.Context, userID uint64, role string) error {
	const execSQL = `UPDATE users SET role = ?, update_time = NOW() WHERE id = ? AND is_deleted = 0`
	result, err := r.db.ExecContext(ctx, execSQL, role, userID)
	if err != nil {
		return fmt.Errorf("update user role: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get update role affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, userID uint64) (*User, error) {
	const query = `
SELECT id, openid, nickname, avatar_url, role, auth_status, account_status, last_login_time
FROM users
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	var user User
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.OpenID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Role,
		&user.AuthStatus,
		&user.AccountStatus,
		&user.LastLoginTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
