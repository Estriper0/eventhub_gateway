package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/models"
	"github.com/Estriper0/EventHub/internal/repositories"
	"github.com/Estriper0/EventHub/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Event struct {
	logger  *slog.Logger
	config  *config.Config
	service service.IEventService
}

func NewEvents(logger *slog.Logger, config *config.Config, service service.IEventService) *Event {
	return &Event{
		logger:  logger,
		config:  config,
		service: service,
	}
}

func (e *Event) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	events, err := e.service.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"events":  []*models.EventResponse{},
					"error":   err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
				"events":  nil,
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting all events",
			"events":  events,
			"error":   "",
		},
	)
}

func (e *Event) GetById(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID field is missing",
				"event":   nil,
				"error":   ErrNoId.Error(),
			},
		)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "Incorrect id",
				"event":   nil,
				"error":   err.Error(),
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	event, err := e.service.GetById(ctx, uuid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"event":   nil,
					"error":   err.Error(),
				},
			)
			return
		}
		if errors.Is(err, repositories.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Event not found",
					"event":   nil,
					"error":   err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
				"event":   nil,
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting event",
			"event":   event,
			"error":   "",
		},
	)
}

func (e *Event) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	var req models.EventCreateRequest

	err := c.ShouldBindJSON(&req)
	if err == nil {
		err = req.Validate()
	}
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "JSON is incorrect",
				"error":   ErrValidateNotPass.Error(),
			},
		)
		return
	}

	id, err := e.service.Create(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"error":   err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusCreated,
		gin.H{
			"code":    http.StatusCreated,
			"message": fmt.Sprintf("Event with ID=%s was created", id.String()),
			"error":   "",
		},
	)
}

func (e *Event) DeleteById(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID field is missing",
				"error":   ErrNoId.Error(),
			},
		)
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "Incorrect id",
				"error":   err.Error(),
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	err = e.service.DeleteById(ctx, uuid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"error":   err.Error(),
				},
			)
			return
		}
		if errors.Is(err, repositories.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Event not found",
					"error":   err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": fmt.Sprintf("The event with the ID=%s has been deleted.", id),
			"error":   "",
		},
	)
}

func (e *Event) Update(c *gin.Context) {
	var req models.EventUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "JSON is incorrect",
				"error":   err.Error(),
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	err := e.service.Update(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"error":   err.Error(),
				},
			)
			return
		}
		if errors.Is(err, repositories.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Event not found",
					"error":   err.Error(),
				},
			)
			return
		}
		if errors.Is(err, repositories.ErrMissingData) {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    http.StatusBadRequest,
					"message": "Missing data",
					"error":   err.Error(),
				},
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": fmt.Sprintf("The event ID=%s has been updated.", (*req.Id).String()),
			"error":   "",
		},
	)
}
