package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Estriper0/EventHub/internal/config"
	pb "github.com/Estriper0/protobuf/gen/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Auth struct {
	logger     *slog.Logger
	config     *config.Config
	authClient pb.AuthClient
}

func NewAuth(logger *slog.Logger, config *config.Config) *Auth {
	// TODO Исправить
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", config.Auth.Host, config.Auth.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &Auth{
		logger:     logger,
		config:     config,
		authClient: pb.NewAuthClient(conn),
	}
}

func (a *Auth) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout)
	defer cancel()

	var req pb.RegisterRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"user_id": nil,
				"message": "JSON is incorrect",
			},
		)
		return
	}

	a.logger.Info("data", slog.Any("data", req.Email), slog.Any("data", req.Password))
	resp, err := a.authClient.Register(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"user_id": nil,
					"message": "Request timed out",
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.AlreadyExists:
			c.JSON(
				http.StatusConflict,
				gin.H{
					"code":    http.StatusConflict,
					"user_id": nil,
					"message": "User exists",
				},
			)
		case codes.InvalidArgument:
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"code":    http.StatusBadRequest,
					"user_id": nil,
					"message": err.Error(),
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"user_id": nil,
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
			"user_id": resp.UserUuid,
			"message": fmt.Sprintf("User with ID=%s was registered", resp.UserUuid),
		},
	)
}

func (a *Auth) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout)
	defer cancel()

	var req pb.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":          http.StatusBadRequest,
				"access_token":  nil,
				"refresh_token": nil,
				"message":       "JSON is incorrect",
			},
		)
		return
	}

	resp, err := a.authClient.Login(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":          http.StatusGatewayTimeout,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Request timed out",
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
					"code":          http.StatusBadRequest,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Invalid credentials",
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":          http.StatusInternalServerError,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":          http.StatusCreated,
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
			"message":       "Successful login user",
		},
	)
}

func (a *Auth) IsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout)
	defer cancel()

	var req pb.IsAdminRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    http.StatusBadRequest,
				"isAdmin": nil,
				"message": "JSON is incorrect",
			},
		)
		return
	}

	resp, err := a.authClient.IsAdmin(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":    http.StatusGatewayTimeout,
					"isAdmin": nil,
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
					"isAdmin": nil,
					"message": "Not found",
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":    http.StatusInternalServerError,
					"isAdmin": nil,
					"message": "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    http.StatusCreated,
			"isAdmin": resp.IsAdmin,
			"message": "Successful user verification for admin",
		},
	)
}

func (a *Auth) Refresh(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout)
	defer cancel()

	var req pb.RefreshRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":          http.StatusBadRequest,
				"access_token":  nil,
				"refresh_token": nil,
				"message":       "JSON is incorrect",
			},
		)
		return
	}

	resp, err := a.authClient.Refresh(ctx, &req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(
				http.StatusGatewayTimeout,
				gin.H{
					"code":          http.StatusGatewayTimeout,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Request timed out",
				},
			)
			return
		}
		st, _ := status.FromError(err)
		code := st.Code()
		switch code {
		case codes.InvalidArgument:
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"code":          http.StatusUnauthorized,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Invalid refresh token",
				},
			)
		case codes.Internal:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"code":          http.StatusInternalServerError,
					"access_token":  nil,
					"refresh_token": nil,
					"message":       "Internal error",
				},
			)
		}
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":          http.StatusOK,
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
			"message":       "Successfully refresh tokens",
		},
	)
}

func (a *Auth) Logout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Timeout)
	defer cancel()

	var req pb.LogoutRequest

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
	_, err = a.authClient.Logout(ctx, &req)
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
				http.StatusUnauthorized,
				gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid refresh token",
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
			"message": "Successfully logout user",
		},
	)
}
