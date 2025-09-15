package repos

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepo struct {
	db *pgxpool.Pool
}

func NewProfileRepo(db *pgxpool.Pool) *ProfileRepo {
	return &ProfileRepo{
		db: db,
	}
}

func (pr *ProfileRepo) GetProfile(ctx context.Context, userID uuid.UUID) (*models.Profile, error) {
	sql := `
		SELECT user_id, firstname, lastname, phone_number, avatar, point, created_at, updated_at
		FROM profile
		WHERE user_id = $1
	`
	var profile models.Profile
	err := pr.db.QueryRow(ctx, sql, userID).Scan(
		&profile.UserID, &profile.FirstName, &profile.LastName, &profile.PhoneNumber,
		&profile.Avatar, &profile.Point, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (pr *ProfileRepo) UpdateProfile(ctx context.Context, p *models.Profile) error {
	now := time.Now()

	setParts := []string{}
	args := []any{}
	argID := 1

	if p.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("firstname = $%d", argID))
		args = append(args, *p.FirstName)
		argID++
	}
	if p.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("lastname = $%d", argID))
		args = append(args, *p.LastName)
		argID++
	}
	if p.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argID))
		args = append(args, *p.PhoneNumber)
		argID++
	}
	if p.Avatar != nil {
		setParts = append(setParts, fmt.Sprintf("avatar = $%d", argID))
		args = append(args, *p.Avatar)
		argID++
	}

	// kalau tidak ada field yang diupdate, keluar aja
	if len(setParts) == 0 {
		return nil
	}

	// updated_at wajib diupdate
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, now)
	argID++

	args = append(args, p.UserID)

	sql := fmt.Sprintf(`
		UPDATE profile
		SET %s
		WHERE user_id = $%d
	`, strings.Join(setParts, ", "), argID)

	_, err := pr.db.Exec(ctx, sql, args...)
	return err
}

func (pr *ProfileRepo) VerifyPassword(ctx context.Context, userID uuid.UUID, oldPassword string) (string, error) {
	var hashedPassword string
	sql := `SELECT password FROM users WHERE id = $1`

	err := pr.db.QueryRow(ctx, sql, userID).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}

	return hashedPassword, nil
}

func (pr *ProfileRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	sql := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err := pr.db.Exec(ctx, sql, hashedPassword, userID)
	return err
}
