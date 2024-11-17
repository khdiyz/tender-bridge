package repository

import (
	"database/sql"
	"errors"
	"strings"
	"tender-bridge/internal/models"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type tenderRepo struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewTenderRepo(db *sqlx.DB, logger *logger.Logger) *tenderRepo {
	return &tenderRepo{
		db:     db,
		logger: logger,
	}
}

func (r *tenderRepo) Create(request models.CreateTender) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INSERT INTO tenders (
		id,
		client_id,
		title,
		description,
		deadline,
		budget,
		file,
		status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	if _, err := r.db.Exec(query,
		id,
		request.ClientId,
		request.Title,
		request.Description,
		request.Deadline,
		request.Budget,
		request.File,
		request.Status,
	); err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *tenderRepo) GetList(filter models.TenderFilter) ([]models.Tender, int, error) {
	baseQuery := `
	SELECT 
		id, 
		client_id,
		title,
		description,
		deadline,
		budget,
		file,
		status 
	FROM tenders WHERE TRUE `

	countQuery := `SELECT COUNT(*) FROM tenders WHERE TRUE `

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.Search != "" {
		conditions = append(conditions, "(title || description) ILIKE :search")
		params["search"] = "%" + filter.Search + "%"
	}

	if filter.ClientId != uuid.Nil {
		conditions = append(conditions, "client_id = :client_id")
		params["client_id"] = filter.ClientId
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
	tenders := []models.Tender{}
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err := rows.Scan(
			&tender.Id,
			&tender.ClientId,
			&tender.Title,
			&tender.Description,
			&tender.Deadline,
			&tender.Budget,
			&tender.File,
			&tender.Status,
		); err != nil {
			r.logger.Error(err)
			return nil, 0, err
		}
		tenders = append(tenders, tender)
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

	return tenders, total, nil
}

func (r *tenderRepo) GetById(id uuid.UUID) (models.Tender, error) {
	var tender models.Tender

	query := `
	SELECT 
		id, 
		client_id,
		title,
		description,
		deadline,
		budget,
		file,
		status 
	FROM tenders WHERE id = $1;`

	if err := r.db.QueryRow(query, id).Scan(
		&tender.Id,
		&tender.ClientId,
		&tender.Title,
		&tender.Description,
		&tender.Deadline,
		&tender.Budget,
		&tender.File,
		&tender.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Tender{}, err
		}
		r.logger.Error(err)
		return models.Tender{}, err
	}

	return tender, nil
}

func (r *tenderRepo) Update(request models.UpdateTender) error {
	query := `
	UPDATE tenders
	SET
		client_id = $2,
		title = $3,
		description = $4,
		deadline = $5,
		budget = $6,
		file = $7,
		status = $8
	WHERE id = $1;`

	// Execute the query
	row, err := r.db.Exec(query,
		request.Id,
		request.ClientId,
		request.Title,
		request.Description,
		request.Deadline,
		request.Budget,
		request.File,
		request.Status,
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

func (r *tenderRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM tenders WHERE id = $1;`

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

func (r *tenderRepo) GetByIds(ids []uuid.UUID) ([]models.Tender, error) {
	tenders := []models.Tender{}

	query := `
	SELECT 
		id, 
		client_id,
		title,
		description,
		deadline,
		budget,
		file,
		status 
	FROM tenders WHERE id = ANY($1);`

	rows, err := r.db.Query(query, pq.Array(ids))
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender models.Tender
		if err = rows.Scan(
			&tender.Id,
			&tender.ClientId,
			&tender.Title,
			&tender.Description,
			&tender.Deadline,
			&tender.Budget,
			&tender.File,
			&tender.Status,
		); err != nil {
			r.logger.Error(err)
			return nil, err
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil
}
