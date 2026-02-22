package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang/internal/repository/_postgres"
	"golang/pkg/modules"

	"github.com/lib/pq"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers(ctx context.Context) ([]modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var users []modules.User
	err := r.db.DB.SelectContext(ctx, &users,
		`SELECT id, name, email, age, created_at
		 FROM users
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var u modules.User
	err := r.db.DB.GetContext(ctx, &u,
		`SELECT id, name, email, age, created_at
		 FROM users
		 WHERE id=$1`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, modules.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateUser(ctx context.Context, u modules.User) (*modules.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	var created modules.User
	created.Name = u.Name
	created.Email = u.Email
	created.Age = u.Age

	row := r.db.DB.QueryRowxContext(ctx,
		`INSERT INTO users (name, email, age)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		u.Name, u.Email, u.Age,
	)

	if err := row.Scan(&created.ID, &created.CreatedAt); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, fmt.Errorf("%w: email already exists", modules.ErrConflict)
		}
		return nil, err
	}

	return &created, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, u modules.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx,
		`UPDATE users
		 SET name=$1, email=$2, age=$3
		 WHERE id=$4`,
		u.Name, u.Email, u.Age, id,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%w: email already exists", modules.ErrConflict)
		}
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return modules.ErrUserNotFound
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.executionTimeout)
	defer cancel()

	res, err := r.db.DB.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return 0, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if ra == 0 {
		return 0, modules.ErrUserNotFound
	}
	return ra, nil
}
