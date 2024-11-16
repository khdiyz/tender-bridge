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

type bidRepo struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewBidRepo(db *sqlx.DB, logger *logger.Logger) *bidRepo {
	return &bidRepo{
		db:     db,
		logger: logger,
	}
}

func (r *bidRepo) Create(request models.CreateBid) (uuid.UUID, error) {
	id := uuid.New()

	query := `
	INSERT INTO bids (
		id,
		contractor_id,
		tender_id,
		price,
		delivery_time,
		comment,
		status
	) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	if _, err := r.db.Exec(query,
		id,
		request.ContractorId,
		request.TenderId,
		request.Price,
		request.DeliveryTime,
		request.Comment,
		request.Status,
	); err != nil {
		r.logger.Error(err)
		return uuid.Nil, err
	}

	return id, nil
}

func (r *bidRepo) GetList(filter models.BidFilter) ([]models.Bid, int, error) {
	baseQuery := `
	SELECT 
		id, 
		contractor_id,
		tender_id,
		price,
		delivery_time,
		comment,
		status
	FROM bids WHERE TRUE `

	countQuery := `SELECT COUNT(*) FROM bids WHERE TRUE `

	conditions := []string{}

	params := map[string]any{
		"limit":  filter.Limit,
		"offset": filter.Offset,
	}

	// Add search condition
	if filter.Search != "" {
		conditions = append(conditions, "comment ILIKE :search")
		params["search"] = "%" + filter.Search + "%"
	}

	if filter.FromPrice > 0 {
		conditions = append(conditions, "price >= :from_price")
		params["from_price"] = filter.FromPrice
	}

	if filter.ToPrice > 0 {
		conditions = append(conditions, "price <= :to_price")
		params["to_price"] = filter.FromPrice
	}

	if filter.TenderId != uuid.Nil {
		conditions = append(conditions, "tender_id = :tender_id")
		params["tender_id"] = filter.TenderId
	}

	if filter.ContractorId != uuid.Nil {
		conditions = append(conditions, "contractor_id = :contractor_id")
		params["contractor_id"] = filter.ContractorId
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
	bids := []models.Bid{}
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		r.logger.Error(err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid models.Bid
		if err := rows.Scan(
			&bid.Id,
			&bid.ContractorId,
			&bid.TenderId,
			&bid.Price,
			&bid.DeliveryTime,
			&bid.Comment,
			&bid.Status,
		); err != nil {
			r.logger.Error(err)
			return nil, 0, err
		}
		bids = append(bids, bid)
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

	return bids, total, nil
}

func (r *bidRepo) GetById(id uuid.UUID) (models.Bid, error) {
	var bid models.Bid

	query := `
	SELECT 
		id, 
		contractor_id,
		tender_id,
		price,
		delivery_time,
		comment,
		status
	FROM bids WHERE id = $1;`

	if err := r.db.QueryRow(query, id).Scan(
		&bid.Id,
		&bid.ContractorId,
		&bid.TenderId,
		&bid.Price,
		&bid.DeliveryTime,
		&bid.Comment,
		&bid.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Bid{}, err
		}
		r.logger.Error(err)
		return models.Bid{}, err
	}

	return bid, nil
}

func (r *bidRepo) Update(request models.UpdateBid) error {
	query := `
	UPDATE bids
	SET
		contractor_id = $2,
		tender_id = $3,
		price = $4,
		delivery_time = $5,
		comment = $6,
		status = $7
	WHERE id = $1;`

	// Execute the query
	row, err := r.db.Exec(query,
		request.Id,
		request.ContractorId,
		request.TenderId,
		request.Price,
		request.DeliveryTime,
		request.Comment,
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

func (r *bidRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM bids WHERE id = $1;`

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
