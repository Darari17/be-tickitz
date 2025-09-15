package repos

import (
	"context"

	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "select id, email, password, role from users where email = $1"
	user := models.User{}
	err := ar.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ar *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := ar.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	user.ID = uuid.New()

	queryInsertUsers := "insert into users (id, email, password, role, created_at) values ($1, $2, $3, 'user', now())"
	_, err = tx.Exec(ctx, queryInsertUsers, user.ID, user.Email, user.Password)
	if err != nil {
		return err
	}

	queryInsertProfile := "insert into profile (user_id, created_at) values ($1, now())"
	_, err = tx.Exec(ctx, queryInsertProfile, user.ID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
