package repositories

import (
	"context"
	"github.com/ZaharBorisenko/jwt-auth/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, first_name, last_name, email, password, role, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query, user.Id, user.UserName, user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := models.User{}
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) (*[]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET 
                 username = $1,
                 first_name = $2,
                 last_name = $3,
                 email = $4,
                 updated_at = $5 WHERE id = $6
                 `

	_, err := r.db.ExecContext(ctx, query,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.Email,
		time.Now(),
		user.Id,
	)

	return err
}
