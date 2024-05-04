package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cats-social/common/logger"
	"cats-social/internal/application/user/service"
	"cats-social/internal/domain"
)

type authHandler struct {
	authService service.AuthServiceContract
}

func NewAuthHandler(router fiber.Router, authService service.AuthServiceContract) {
	handler := authHandler{
		authService: authService,
	}

	authRouter := router.Group("/user")

	authRouter.Post("/register", handler.Register)
	authRouter.Post("/login", handler.Login)
}

func (h authHandler) Register(c *fiber.Ctx) error {
	callerInfo := "[authHandler.User]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req, res := &registerRequest{}, baseResponse{}
	if err := c.BodyParser(req); err != nil {
		l.Error("error binding data",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if err := req.validate(); err != nil {
		l.Error("error validate data",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	registerData := domain.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	user, err := h.authService.Register(userCtx, registerData)
	if err != nil {
		if errors.Is(err, domain.DuplicateEmailError) {
			l.Error("email already exists",
				zap.Error(err),
			)
			res = baseResponse{
				Message: duplicateEmailErrorMessage,
				Data: fiber.Map{
					"error": err.Error(),
				},
			}
			return c.Status(http.StatusConflict).JSON(res)
		}
		l.Error("error register user",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InternalServerErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	token, err := h.authService.GenerateToken(userCtx, user)
	if err != nil {
		l.Error("error generate token",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InternalServerErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	res = baseResponse{
		Message: successRegisterMessage,
		Data: authResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: token,
		},
	}
	return c.Status(http.StatusCreated).JSON(res)
}

func (h authHandler) Login(c *fiber.Ctx) error {
	callerInfo := "[authHandler.Login]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	req, res := &loginRequest{}, baseResponse{}
	if err := c.BodyParser(req); err != nil {
		l.Error("error binding data",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if err := req.validate(); err != nil {
		l.Error("error validate data",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	loginData := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.authService.Login(userCtx, loginData)
	if err != nil {
		switch {
		case errors.Is(err, domain.UserNotFoundError):
			l.Error("user not found",
				zap.Error(err),
			)
			res = baseResponse{
				Message: userNotFoundErrorMessage,
				Data: fiber.Map{
					"error": err.Error(),
				},
			}

			return c.Status(http.StatusNotFound).JSON(res)

		case errors.Is(err, domain.InvalidPassword):
			l.Error("invalid password",
				zap.Error(err),
			)
			res = baseResponse{
				Message: invalidPasswordMessage,
				Data: fiber.Map{
					"error": err.Error(),
				},
			}

			return c.Status(http.StatusBadRequest).JSON(res)

		default:
			l.Error("error login user",
				zap.Error(err),
			)
			res = baseResponse{
				Message: domain.InternalServerErrorMessage,
				Data: fiber.Map{
					"error": err.Error(),
				},
			}

			return c.Status(http.StatusInternalServerError).JSON(res)
		}
	}

	token, err := h.authService.GenerateToken(userCtx, user)
	if err != nil {
		l.Error("error generate token",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InternalServerErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	res = baseResponse{
		Message: successLoginMessage,
		Data: authResponse{
			Email:       user.Email,
			Name:        user.Name,
			AccessToken: token,
		},
	}

	return c.Status(http.StatusOK).JSON(res)
}
