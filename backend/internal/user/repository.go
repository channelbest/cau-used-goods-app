package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID            uint64  `json:"id"`
	Nickname      *string `json:"nickname"`
	AvatarURL     *string `json:"avatarUrl"`
	StudentID     *string `json:"studentId"`
	RealName      *string `json:"realName"`
	College       *string `json:"college"`
	Phone         *string `json:"phone"`
	Role          string  `json:"role"`
	AuthStatus    string  `json:"authStatus"`
	AccountStatus string  `json:"accountStatus"`
}

type StudentVerification struct {
	UserID        uint64  `json:"userId,omitempty"`
	StudentID     *string `json:"studentId"`
	RealName      *string `json:"realName"`
	College       *string `json:"college"`
	AuthStatus    string  `json:"authStatus"`
	AccountStatus string  `json:"accountStatus"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, userID uint64) (*User, error) {
	const query = `
SELECT id, nickname, avatar_url, student_id, real_name, college, phone, role, auth_status, account_status
FROM users
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	var user User
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Nickname,
		&user.AvatarURL,
		&user.StudentID,
		&user.RealName,
		&user.College,
		&user.Phone,
		&user.Role,
		&user.AuthStatus,
		&user.AccountStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uint64, nickname, avatarURL, phone *string) error {
	const execSQL = `
UPDATE users
SET nickname = COALESCE(?, nickname),
    avatar_url = COALESCE(?, avatar_url),
    phone = COALESCE(?, phone),
    update_time = NOW()
WHERE id = ? AND is_deleted = 0`

	result, err := r.db.ExecContext(ctx, execSQL, nickname, avatarURL, phone, userID)
	if err != nil {
		return fmt.Errorf("update user profile: %w", err)
	}
	return checkAffected(result, "user not found")
}

func (r *Repository) SubmitStudentVerification(ctx context.Context, userID uint64, studentID, realName, college string) error {
	const execSQL = `
UPDATE users
SET student_id = ?,
    real_name = ?,
    college = ?,
    auth_status = 'PENDING',
    update_time = NOW()
WHERE id = ? AND is_deleted = 0`

	result, err := r.db.ExecContext(ctx, execSQL, studentID, realName, college, userID)
	if err != nil {
		return fmt.Errorf("submit student verification: %w", err)
	}
	return checkAffected(result, "user not found")
}

func (r *Repository) FindStudentVerification(ctx context.Context, userID uint64) (*StudentVerification, error) {
	const query = `
SELECT id, student_id, real_name, college, auth_status, account_status
FROM users
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	var verification StudentVerification
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&verification.UserID,
		&verification.StudentID,
		&verification.RealName,
		&verification.College,
		&verification.AuthStatus,
		&verification.AccountStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find student verification: %w", err)
	}
	return &verification, nil
}

func (r *Repository) ListStudentVerifications(ctx context.Context, status string) ([]StudentVerification, error) {
	const query = `
SELECT id, student_id, real_name, college, auth_status, account_status
FROM users
WHERE auth_status = ? AND is_deleted = 0
ORDER BY update_time DESC`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("list student verifications: %w", err)
	}
	defer rows.Close()

	items := make([]StudentVerification, 0)
	for rows.Next() {
		var item StudentVerification
		if err := rows.Scan(
			&item.UserID,
			&item.StudentID,
			&item.RealName,
			&item.College,
			&item.AuthStatus,
			&item.AccountStatus,
		); err != nil {
			return nil, fmt.Errorf("scan student verification: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student verifications: %w", err)
	}
	return items, nil
}

func (r *Repository) ReviewStudentVerification(ctx context.Context, adminID uint64, userID uint64, authStatus string, description string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin review student verification tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, `
UPDATE users
SET auth_status = ?, update_time = NOW()
WHERE id = ? AND auth_status = 'PENDING' AND is_deleted = 0`, authStatus, userID)
	if err != nil {
		return fmt.Errorf("review student verification: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get review affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("student verification is not pending or user not found")
	}

	operationType := "STUDENT_VERIFY_REJECT"
	if authStatus == "VERIFIED" {
		operationType = "STUDENT_VERIFY_APPROVE"
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO admin_logs (admin_id, operation_type, target_type, target_id, description, create_time)
VALUES (?, ?, 'USER', ?, ?, NOW())`, adminID, operationType, userID, description); err != nil {
		return fmt.Errorf("create student verification admin log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit review student verification tx: %w", err)
	}
	committed = true
	return nil
}

func checkAffected(result sql.Result, notFoundMsg string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf(notFoundMsg)
	}
	return nil
}
