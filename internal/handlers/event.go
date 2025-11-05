package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Estriper0/EventHub/internal/config"
	pb "github.com/Estriper0/protobuf/gen/event"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Event struct {
	logger      *slog.Logger
	config      *config.Config
	eventClient pb.EventClient
}

func NewEvent(logger *slog.Logger, config *config.Config) *Event {
	// TODO Исправить
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", config.Event.Host, config.Event.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &Event{
		logger:      logger,
		config:      config,
		eventClient: pb.NewEventClient(conn),
	}
}

func (e *Event) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	req := &pb.EmptyRequest{}
	resp, err := e.eventClient.GetAll(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"events":  nil,
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
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting all events",
			"events":  resp.Events,
		},
	)
}

func (e *Event) GetAllByCreator(c *gin.Context) {
	creator, ok := c.Params.Get("creator")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "Creator field is missing",
				"events":  nil,
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	req := &pb.GetAllByCreatorRequest{Creator: creator}
	resp, err := e.eventClient.GetAllByCreator(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"events":  nil,
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.InvalidArgument:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
					"events":  nil,
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
					"events":  nil,
				},
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting all events",
			"events":  resp.Events,
		},
	)
}

func (e *Event) GetAllByStatus(c *gin.Context) {
	sts, ok := c.Params.Get("status")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "Creator field is missing",
				"events":  nil,
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	req := &pb.GetAllByStatusRequest{Status: sts}
	resp, err := e.eventClient.GetAllByStatus(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"events":  nil,
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.InvalidArgument:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
					"events":  nil,
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
					"events":  nil,
				},
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting all events",
			"events":  resp.Events,
		},
	)
}

func (e *Event) GetById(c *gin.Context) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID field is missing",
				"event":   nil,
			},
		)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID is not a number",
				"event":   nil,
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	req := &pb.GetByIdRequest{Id: int64(id)}
	event, err := e.eventClient.GetById(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"event":   nil,
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.NotFound:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Not found",
					"event":   nil,
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
					"event":   nil,
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": "Successful getting event",
			"event":   event,
		},
	)
}

func (e *Event) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	var req pb.CreateRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "JSON is incorrect",
			},
		)
		return
	}
	req.Creator = c.GetString("user_id")

	resp, err := e.eventClient.Create(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.InvalidArgument:
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusCreated,
		gin.H{
			"code":    http.StatusCreated,
			"message": fmt.Sprintf("Event with ID=%d was created", resp.Id),
		},
	)
}

func (e *Event) DeleteById(c *gin.Context) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID field is missing",
			},
		)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "ID is not a number",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	if !e.UserVerification(ctx, c, id) {
		return
	}

	req := &pb.DeleteByIdRequest{Id: int64(id)}
	_, err = e.eventClient.DeleteById(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.NotFound:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Not found",
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": fmt.Sprintf("The event with the ID=%d has been deleted.", id),
		},
	)
}

func (e *Event) Update(c *gin.Context) {
	var req pb.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"message": "JSON is incorrect",
			},
		)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), e.config.Timeout)
	defer cancel()

	if !e.UserVerification(ctx, c, int(req.Id)) {
		return
	}

	_, err := e.eventClient.Update(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.NotFound:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Not found",
				},
			)
		case codes.InvalidArgument:
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    http.StatusBadRequest,
					"message": err.Error(),
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusOK,
			"message": fmt.Sprintf("The event ID=%d has been updated.", req.Id),
		},
	)
}

func (e *Event) UserVerification(ctx context.Context, c *gin.Context, id int) bool {
	req := &pb.GetByIdRequest{Id: int64(id)}
	resp, err := e.eventClient.GetById(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"message": "Request timed out",
					"event":   nil,
				},
			)
			return false
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.NotFound:
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"code":    http.StatusNotFound,
					"message": "Not found",
					"event":   nil,
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal error",
					"event":   nil,
				},
			)
		}
		return false
	}
	if resp.Creator != c.GetString("user_id") {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"code":    http.StatusForbidden,
				"message": "The user does not have access to the requested resource.",
				"event":   nil,
			},
		)
		return false
	}
	return true
}
