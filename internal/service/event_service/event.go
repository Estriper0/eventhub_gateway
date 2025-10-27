package event_service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	"github.com/google/uuid"
)

type EventService struct {
	eventRepo repositories.IEventRepository
	logger    *slog.Logger
}

func New(repo repositories.IEventRepository, logger *slog.Logger) *EventService {
	return &EventService{
		eventRepo: repo,
		logger:    logger,
	}
}

func (s *EventService) GetAll(ctx context.Context) ([]*models.EventResponse, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error(
			"Error getting all events",
			slog.String("err", err.Error()),
		)
		return nil, err
	}
	s.logger.Info(
		"Successful getting all events",
	)
	return events, nil
}

func (s *EventService) Create(ctx context.Context, event *models.EventCreateRequest) (uuid.UUID, error) {
	id, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		s.logger.Error(
			"Error create event",
			slog.String("err", err.Error()),
		)
		return uuid.UUID{}, err
	}
	s.logger.Info(
		"Successful create event",
		slog.String("id", id.String()),
	)
	return id, nil
}

func (s *EventService) GetById(ctx context.Context, id uuid.UUID) (*models.EventResponse, error) {
	event, err := s.eventRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Warn(
				"Event not found",
				slog.String("id", id.String()),
				slog.String("err", err.Error()),
			)
			return nil, err
		}
		s.logger.Error(
			"Error getting event",
			slog.String("id", id.String()),
			slog.String("err", err.Error()),
		)
		return nil, err
	}
	s.logger.Info(
		"Successful getting event",
		slog.String("id", id.String()),
	)
	return event, nil
}

func (s *EventService) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := s.eventRepo.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Warn(
				"Event not found",
				slog.String("id", id.String()),
				slog.String("err", err.Error()),
			)
			return err
		}
		s.logger.Error(
			"Error delete event",
			slog.String("id", id.String()),
			slog.String("err", err.Error()),
		)
		return err
	}
	s.logger.Info(
		"Successful delete event",
		slog.String("id", id.String()),
	)
	return nil
}

func (s *EventService) Update(ctx context.Context, event *models.EventUpdateRequest) error {
	err := s.eventRepo.Update(ctx, event)
	if err != nil {
		if errors.Is(err, repositories.ErrMissingData) {
			s.logger.Warn(
				"Missing data",
				slog.String("id", (*event.Id).String()),
				slog.String("err", err.Error()),
			)
			return err
		}
		if errors.Is(err, repositories.ErrRecordNotFound) {
			s.logger.Warn(
				"Event not found",
				slog.String("id", (*event.Id).String()),
				slog.String("err", err.Error()),
			)
			return err
		}
		s.logger.Error(
			"Error update event",
			slog.String("id", (*event.Id).String()),
			slog.String("err", err.Error()),
		)
		return err
	}
	s.logger.Info(
		"Successful update event",
		slog.String("id", (*event.Id).String()),
	)
	return nil
}
