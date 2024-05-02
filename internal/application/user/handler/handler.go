package handler

import (
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
		res = invalidRequestBody
		res.Data = err
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	if err := req.validate(); err != nil {
		l.Error("error validate data",
			zap.Error(err),
		)
		res = invalidRequestBody
		res.Data = fiber.Map{
			"error": err.Error(),
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
		l.Error("error register user",
			zap.Error(err),
		)
		res = invalidRequestBody
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	return c.JSON(fiber.Map{
		"message": user,
	})
}
