package event_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*models.EventResponse, error) {
	query := "SELECT * FROM event WHERE id = $1"
	event := &models.EventResponse{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.Id,
		&event.Title,
		&event.About,
		&event.StartDate,
		&event.Location,
		&event.Status,
		&event.MaxAttendees,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositories.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *EventRepository) GetAll(
	ctx context.Context,
) ([]*models.EventResponse, error) {
	query := "SELECT * FROM event ORDER BY title"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []*models.EventResponse{}

	for rows.Next() {
		event := &models.EventResponse{}
		err := rows.Scan(
			&event.Id,
			&event.Title,
			&event.About,
			&event.StartDate,
			&event.Location,
			&event.Status,
			&event.MaxAttendees,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *EventRepository) Create(
	ctx context.Context,
	event *models.EventCreateRequest,
) (uuid.UUID, error) {
	query := "INSERT INTO event (id, title, about, start_date, location, status, max_attendees) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	id := uuid.New()
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		event.Title,
		event.About,
		event.StartDate,
		event.Location,
		event.Status,
		event.MaxAttendees,
	).Scan(&id)

	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *EventRepository) DeleteById(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := "DELETE FROM event WHERE id = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}
	return nil
}

func (r *EventRepository) Update(
	ctx context.Context,
	event *models.EventUpdateRequest,
) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	statement := psql.Update("event").Where(sq.Eq{"id": *event.Id})
	var c int
	if event.Title != nil {
		statement = statement.Set("title", *event.Title)
		c++
	}
	if event.About != nil {
		statement = statement.Set("about", *event.About)
		c++
	}
	if event.StartDate != nil {
		statement = statement.Set("start_date", *event.StartDate)
		c++
	}
	if event.Location != nil {
		statement = statement.Set("location", *event.Location)
		c++
	}
	if event.Status != nil {
		statement = statement.Set("status", *event.Status)
		c++
	}
	if event.MaxAttendees != nil {
		statement = statement.Set("max_attendees", *event.MaxAttendees)
		c++
	}

	if c == 0 {
		return repositories.ErrMissingData
	}

	sql, args, err := statement.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	i, _ := res.RowsAffected()
	if i == 0 {
		return repositories.ErrRecordNotFound
	}

	return nil
}
