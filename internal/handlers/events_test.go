package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	"github.com/Estriper0/EventHub/internal/service/mocks"
	"github.com/Estriper0/EventHub/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEventHandlers_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("ENV", "test")
	config := config.New()
	logger := logger.GetLogger(config.Env)

	test_cases := []struct {
		Name            string
		ServiceResponse []*models.EventResponse
		ServiceErr      error
		Code            int
		Message         string
	}{
		{
			Name: "get_all",
			ServiceResponse: []*models.EventResponse{
				{
					Id:           uuid.New(),
					Title:        "Tech Conference 2025",
					About:        "Annual technology conference",
					StartDate:    time.Now(),
					Location:     "San Francisco, CA",
					Status:       "draft",
					MaxAttendees: 50,
				},
			},
			ServiceErr: nil,
			Code:       200,
			Message:    "Successful getting all events",
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			response := struct {
				Code    int                     `json:"code"`
				Error   string                  `json:"error"`
				Events  []*models.EventResponse `json:"events"`
				Message string                  `json:"message"`
			}{
				Code:    tc.Code,
				Events:  tc.ServiceResponse,
				Message: tc.Message,
			}
			if tc.ServiceErr != nil {
				response.Error = tc.ServiceErr.Error()
			}

			res, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("GET", "/events", nil)
			if err != nil {
				t.Fatal(err)
			}

			c.Request = req

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockIEventService(ctrl)
			mockService.EXPECT().GetAll(gomock.Any()).Return(tc.ServiceResponse, tc.ServiceErr)

			handler := NewEvents(logger, config, mockService)
			handler.GetAll(c)
			assert.Equal(t, rr.Body.String(), string(res))
		})
	}
}

func TestEventHandlers_GetById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("ENV", "test")
	config := config.New()
	logger := logger.GetLogger(config.Env)

	id := uuid.New()

	test_cases := []struct {
		Name            string
		Id              uuid.UUID
		ServiceResponse *models.EventResponse
		ServiceErr      error
		Code            int
		Message         string
	}{
		{
			Name: "get_all",
			Id:   id,
			ServiceResponse: &models.EventResponse{
				Id:           id,
				Title:        "Tech Conference 2025",
				About:        "Annual technology conference",
				StartDate:    time.Now(),
				Location:     "San Francisco, CA",
				Status:       "draft",
				MaxAttendees: 50,
			},
			ServiceErr: nil,
			Code:       200,
			Message:    "Successful getting event",
		},
		{
			Name:            "not_found",
			Id:              id,
			ServiceResponse: nil,
			ServiceErr:      repositories.ErrRecordNotFound,
			Code:            404,
			Message:         "Event not found",
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			response := struct {
				Code    int                   `json:"code"`
				Error   string                `json:"error"`
				Event   *models.EventResponse `json:"event"`
				Message string                `json:"message"`
			}{
				Code:    tc.Code,
				Event:   tc.ServiceResponse,
				Message: tc.Message,
			}
			if tc.ServiceErr != nil {
				response.Error = tc.ServiceErr.Error()
			}

			res, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("GET", fmt.Sprintf("/events/%s", tc.Id.String()), nil)
			if err != nil {
				t.Fatal(err)
			}

			c.Request = req
			c.Params = []gin.Param{
				{Key: "id", Value: tc.Id.String()},
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockIEventService(ctrl)
			mockService.EXPECT().GetById(gomock.Any(), tc.Id).Return(tc.ServiceResponse, tc.ServiceErr)
			handler := NewEvents(logger, config, mockService)
			handler.GetById(c)
			assert.Equal(t, rr.Body.String(), string(res))
		})
	}
}

func TestEventHandlers_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("ENV", "test")
	config := config.New()
	logger := logger.GetLogger(config.Env)

	test_cases := []struct {
		Name        string
		Id          uuid.UUID
		RequestBody *models.EventCreateRequest
		ServiceErr  error
		Code        int
		Message     string
	}{
		{
			Name: "create",
			Id:   uuid.New(),
			RequestBody: &models.EventCreateRequest{
				Title:        "Tech Conference 2025",
				About:        "Annual technology conference",
				StartDate:    time.Now().UTC(),
				Location:     "San Francisco, CA",
				Status:       "draft",
				MaxAttendees: 50,
			},
			ServiceErr: nil,
			Code:       201,
			Message:    "",
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ServiceErr == nil {
				tc.Message = fmt.Sprintf("Event with ID=%s was created", tc.Id.String())
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			response := struct {
				Code    int    `json:"code"`
				Error   string `json:"error"`
				Message string `json:"message"`
			}{
				Code:    tc.Code,
				Message: tc.Message,
			}
			if tc.ServiceErr != nil {
				response.Error = tc.ServiceErr.Error()
			}

			res, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}

			body, err := json.Marshal(tc.RequestBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			c.Request = req
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockIEventService(ctrl)
			mockService.EXPECT().Create(gomock.Any(), tc.RequestBody).Return(tc.Id, tc.ServiceErr)
			handler := NewEvents(logger, config, mockService)
			handler.Create(c)
			assert.Equal(t, rr.Body.String(), string(res))
		})
	}
}

func TestEventHandlers_DeleteById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("ENV", "test")
	config := config.New()
	logger := logger.GetLogger(config.Env)

	test_cases := []struct {
		Name       string
		Id         uuid.UUID
		ServiceErr error
		Code       int
		Message    string
	}{
		{
			Name:       "delete",
			Id:         uuid.New(),
			ServiceErr: nil,
			Code:       200,
		},
		{
			Name:       "not_found",
			Id:         uuid.New(),
			ServiceErr: repositories.ErrRecordNotFound,
			Code:       404,
			Message:    "Event not found",
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ServiceErr == nil {
				tc.Message = fmt.Sprintf("The event with the ID=%s has been deleted.", tc.Id.String())
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			response := struct {
				Code    int    `json:"code"`
				Error   string `json:"error"`
				Message string `json:"message"`
			}{
				Code:    tc.Code,
				Error:   "",
				Message: tc.Message,
			}
			if tc.ServiceErr != nil {
				response.Error = tc.ServiceErr.Error()
			}

			res, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/events/%s", tc.Id.String()), nil)
			if err != nil {
				t.Fatal(err)
			}

			c.Request = req
			c.Params = []gin.Param{
				{Key: "id", Value: tc.Id.String()},
			}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockIEventService(ctrl)
			mockService.EXPECT().DeleteById(gomock.Any(), tc.Id).Return(tc.ServiceErr)
			handler := NewEvents(logger, config, mockService)
			handler.DeleteById(c)
			assert.Equal(t, rr.Body.String(), string(res))
		})
	}
}

func TestEventHandlers_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("ENV", "test")
	config := config.New()
	logger := logger.GetLogger(config.Env)

	Id := uuid.New()
	Title := "Tech Conference 2025"
	About := "Annual technology conference"
	StartDate := time.Now().UTC()
	Location := "San Francisco, CA"
	Status := models.StatusDraft
	MaxAttendees := 50

	test_cases := []struct {
		Name        string
		RequestBody *models.EventUpdateRequest
		ServiceErr  error
		Code        int
		Message     string
	}{
		{
			Name: "update",
			RequestBody: &models.EventUpdateRequest{
				Id:           &Id,
				Title:        &Title,
				About:        &About,
				StartDate:    &StartDate,
				Location:     &Location,
				Status:       &Status,
				MaxAttendees: &MaxAttendees,
			},
			ServiceErr: nil,
			Code:       200,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ServiceErr == nil {
				tc.Message = fmt.Sprintf("The event ID=%s has been updated.", (*tc.RequestBody.Id).String())
			}
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)

			response := struct {
				Code    int    `json:"code"`
				Error   string `json:"error"`
				Message string `json:"message"`
			}{
				Code:    tc.Code,
				Message: tc.Message,
			}
			if tc.ServiceErr != nil {
				response.Error = tc.ServiceErr.Error()
			}

			res, err := json.Marshal(response)
			if err != nil {
				t.Fatal(err)
			}

			body, err := json.Marshal(tc.RequestBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			c.Request = req
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockIEventService(ctrl)
			mockService.EXPECT().Update(gomock.Any(), tc.RequestBody).Return(tc.ServiceErr)
			handler := NewEvents(logger, config, mockService)
			handler.Update(c)
			assert.Equal(t, rr.Body.String(), string(res))
		})
	}
}
