package event_service

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	"github.com/Estriper0/EventHub/internal/repositories/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEventService_GetAll(t *testing.T) {
	test_cases := []struct {
		Name     string
		Response []*models.EventResponse
		Err      error
	}{
		{
			Name:     "empty",
			Response: []*models.EventResponse{},
			Err:      nil,
		},
		{
			Name: "get_all",
			Response: []*models.EventResponse{
				{
					Id:           uuid.New(),
					Title:        "Party",
					About:        "Party",
					StartDate:    time.Now(),
					Location:     "Moscow",
					Status:       models.StatusCancelled,
					MaxAttendees: 15,
				},
				{
					Id:           uuid.New(),
					Title:        "Party",
					About:        "Party",
					StartDate:    time.Now(),
					Location:     "Paris",
					Status:       models.StatusDraft,
					MaxAttendees: 20,
				},
				{
					Id:           uuid.New(),
					Title:        "Party",
					About:        "Party",
					StartDate:    time.Now(),
					Location:     "London",
					Status:       models.StatusCompleted,
					MaxAttendees: 25,
				},
			},
			Err: nil,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo := mocks.NewMockIEventRepository(ctrl)
			mockRepo.EXPECT().GetAll(ctx).Return(tc.Response, tc.Err)

			logger := slog.New(
				slog.NewJSONHandler(
					os.Stdin,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			service := New(mockRepo, logger)
			res, err := service.GetAll(ctx)
			assert.Equal(t, tc.Response, res)
			assert.Equal(t, tc.Err, err)
		})
	}
}

func TestEventService_Create(t *testing.T) {
	test_cases := []struct {
		Name     string
		Request  *models.EventCreateRequest
		Response uuid.UUID
		Err      error
	}{
		{
			Name: "create_event",
			Request: &models.EventCreateRequest{
				Title:        "Party",
				StartDate:    time.Now(),
				Location:     "Moscow",
				Status:       models.StatusCancelled,
				MaxAttendees: 15,
			},
			Response: uuid.New(),
			Err:      nil,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo := mocks.NewMockIEventRepository(ctrl)
			mockRepo.EXPECT().Create(ctx, tc.Request).Return(tc.Response, tc.Err)

			logger := slog.New(
				slog.NewJSONHandler(
					os.Stdin,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			service := New(mockRepo, logger)
			res, err := service.Create(ctx, tc.Request)
			assert.Equal(t, tc.Response, res)
			assert.Equal(t, tc.Err, err)
		})
	}
}

func TestEventService_GetById(t *testing.T) {
	test_cases := []struct {
		Name     string
		Request  uuid.UUID
		Response *models.EventResponse
		Err      error
	}{
		{
			Name:    "get_event",
			Request: uuid.New(),
			Response: &models.EventResponse{
				Title:        "Party",
				StartDate:    time.Now(),
				Location:     "Moscow",
				Status:       models.StatusCancelled,
				MaxAttendees: 15,
			},
			Err: nil,
		},
		{
			Name:     "not found",
			Request:  uuid.New(),
			Response: nil,
			Err:      repositories.ErrRecordNotFound,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo := mocks.NewMockIEventRepository(ctrl)
			mockRepo.EXPECT().GetById(ctx, tc.Request).Return(tc.Response, tc.Err)

			logger := slog.New(
				slog.NewJSONHandler(
					os.Stdin,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			service := New(mockRepo, logger)
			res, err := service.GetById(ctx, tc.Request)
			assert.Equal(t, tc.Response, res)
			assert.Equal(t, tc.Err, err)
		})
	}
}

func TestEventService_DeleteById(t *testing.T) {
	test_cases := []struct {
		Name    string
		Request uuid.UUID
		Err     error
	}{
		{
			Name:    "delete_event",
			Request: uuid.New(),
			Err:     nil,
		},
		{
			Name:    "not found",
			Request: uuid.New(),
			Err:     repositories.ErrRecordNotFound,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo := mocks.NewMockIEventRepository(ctrl)
			mockRepo.EXPECT().DeleteById(ctx, tc.Request).Return(tc.Err)

			logger := slog.New(
				slog.NewJSONHandler(
					os.Stdin,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			service := New(mockRepo, logger)
			err := service.DeleteById(ctx, tc.Request)
			assert.Equal(t, tc.Err, err)
		})
	}
}

func TestEventService_Update(t *testing.T) {
	id := uuid.New()
	title := "title"
	test_cases := []struct {
		Name    string
		Request *models.EventUpdateRequest
		Err     error
	}{
		{
			Name: "update_event",
			Request: &models.EventUpdateRequest{
				Id:    &id,
				Title: &title,
			},
			Err: nil,
		},
		{
			Name: "not_found",
			Request: &models.EventUpdateRequest{
				Id:    &id,
				Title: &title,
			},
			Err: repositories.ErrRecordNotFound,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo := mocks.NewMockIEventRepository(ctrl)
			mockRepo.EXPECT().Update(ctx, tc.Request).Return(tc.Err)

			logger := slog.New(
				slog.NewJSONHandler(
					os.Stdin,
					&slog.HandlerOptions{
						Level: slog.LevelDebug,
					},
				),
			)

			service := New(mockRepo, logger)
			err := service.Update(ctx, tc.Request)
			assert.Equal(t, tc.Err, err)
		})
	}
}
