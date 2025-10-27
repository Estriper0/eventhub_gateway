package repositories

import (
	"context"

	"github.com/Estriper0/EventHub/internal/models"
	"github.com/google/uuid"
)

type IEventRepository interface {
	GetById(
		ctx context.Context,
		id uuid.UUID,
	) (*models.EventResponse, error)
	Create(
		ctx context.Context,
		event *models.EventCreateRequest,
	) (uuid.UUID, error)
	GetAll(
		ctx context.Context,
	) ([]*models.EventResponse, error)
	DeleteById(
		ctx context.Context,
		id uuid.UUID,
	) error
	Update(
		ctx context.Context,
		event *models.EventUpdateRequest,
	) error
}
