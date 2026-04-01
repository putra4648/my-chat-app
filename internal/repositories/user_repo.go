package repositories

import (
	"context"
	"errors"
	"putra4648/my-chat-app/internal/models"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx, "SELECT id, username, email, password_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, pgErr
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersWithoutUserLogin(ctx context.Context, id string) ([]*models.User, error) {
	var users []*models.User
	rows, err := r.db.Query(ctx, "SELECT id, username, email FROM users WHERE id != $1", id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error(pgErr.Error())
			return nil, pgErr
		}
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				return nil, pgErr
			}
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	var createdUser models.User
	// check if exist
	var isUserExist bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&isUserExist)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, pgErr
		}
		return nil, err
	}
	if isUserExist {
		return nil, errors.New("user already exists")
	}
	err = r.db.QueryRow(ctx, "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Email, user.PasswordHash).Scan(&createdUser.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, pgErr
		}
		return nil, err
	}
	return &createdUser, nil
}
