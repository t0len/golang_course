package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"practice5/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var allowedColumns = map[string]bool{
	"id":         true,
	"name":       true,
	"email":      true,
	"gender":     true,
	"birth_date": true,
}

func (r *Repository) GetPaginatedUsers(params models.FilterParams) (models.PaginatedResponse, error) {

	orderBy := "id"
	if params.OrderBy != "" {
		col := strings.ToLower(params.OrderBy)
		if allowedColumns[col] {
			orderBy = col
		}
	}

	args := []interface{}{}
	conditions := []string{}
	argIdx := 1

	if params.ID != "" {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, params.ID)
		argIdx++
	}
	if params.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIdx))
		args = append(args, "%"+params.Name+"%")
		argIdx++
	}
	if params.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email ILIKE $%d", argIdx))
		args = append(args, "%"+params.Email+"%")
		argIdx++
	}
	if params.Gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, params.Gender)
		argIdx++
	}
	if params.BirthDate != "" {
		conditions = append(conditions, fmt.Sprintf("birth_date::date = $%d", argIdx))
		args = append(args, params.BirthDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		`SELECT id, name, email, gender, birth_date FROM users %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIdx, argIdx+1,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
func (r *Repository) GetCommonFriends(userID1, userID2 uuid.UUID) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2
	`
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
