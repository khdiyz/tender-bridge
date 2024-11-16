package repository

import (
	"database/sql"
	"errors"
	"strings"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	errNoRowsAffected = errors.New("no rows affected")
)

type userRepo struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewUserRepo(db *sqlx.DB, logger *logger.Logger) *userRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

func (r *userRepo) Create(request models.CreateUser) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INSERT INTO users (
		id,
		role,
		username,
		email,
		password
	) VALUES ($1, $2, $3, $4, $5);`

	if _, err := r.db.Exec(query,
		id,
		request.Role,
		request.Username,
		request.Email,
		request.Password,
	); err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepo) GetList(filter models.UserFilter) ([]models.User, int, error) {
	baseQuery := `
	SELECT 
		id, 
		role, 
		username, 
		email, 
		password 
	FROM users WHERE TRUE `

	countQuery := `SELECT COUNT(*) FROM users WHERE TRUE `

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.Search != "" {
		conditions = append(conditions, "(username || email) ILIKE :search")
		params["search"] = "%" + filter.Search + "%"
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Add pagination
	baseQuery += " LIMIT :limit OFFSET :offset"

	// Execute the main query
	users := []models.User{}
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Id,
			&user.Role,
			&user.Username,
			&user.Email,
			&user.Password,
		); err != nil {
			r.logger.Error(err)
			return nil, 0, err
		}
		users = append(users, user)
	}

	// Execute the count query
	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}
	countQuery = r.db.Rebind(countQuery)

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) GetById(id uuid.UUID) (models.User, error) {
	var user models.User

	query := `
	SELECT 
		id, 
		role, 
		username, 
		email, 
		password 
	FROM users 
	WHERE id = $1;`

	if err := r.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.Password,
	); err != nil {
		r.logger.Error(err)
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) Update(request models.UpdateUser) error {
	query := `
	UPDATE users
	SET
		role = $2,
		username = $3,
		email = $4,
		password = $5 
	WHERE
		id = $1;`

	// Execute the query
	row, err := r.db.Exec(query,
		request.Id,
		request.Role,
		request.Username,
		request.Email,
		request.Password,
	)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		r.logger.Error(err)
		return err
	}

	if rowAffected == 0 {
		return errNoRowsAffected
	}

	return nil
}

func (r *userRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1;`

	// Execute the query
	row, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	rowAffected, err := row.RowsAffected()
	if err != nil {
		r.logger.Error(err)
		return err
	}

	if rowAffected == 0 {
		return errNoRowsAffected
	}

	return nil
}

func (r *userRepo) GetByUsername(username string) (models.User, error) {
	var user models.User

	query := `
	SELECT 
		id, 
		role, 
		username, 
		email, 
		password 
	FROM users 
	WHERE username = $1;`

	if err := r.db.QueryRow(query, username).Scan(
		&user.Id,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		}
		r.logger.Error(err)
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepo) GetByEmail(email string) (models.User, error) {
	var user models.User

	query := `
	SELECT 
		id, 
		role, 
		username, 
		email, 
		password 
	FROM users 
	WHERE email = $1;`

	if err := r.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		}
		r.logger.Error(err)
		return models.User{}, err
	}

	return user, nil
}
