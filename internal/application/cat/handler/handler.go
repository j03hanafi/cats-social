package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/logger"
	"cats-social/internal/application/cat/service"
	"cats-social/internal/domain"
)

const (
	catIDFromParam = "catID"
)

type catHandler struct {
	catService service.CatServiceContract
}

func NewCatHandler(router fiber.Router, jwtMiddleware fiber.Handler, catService service.CatServiceContract) {
	handler := catHandler{
		catService: catService,
	}

	catRouter := router.Group("/cat")

	catRouter.Use(jwtMiddleware)
	catRouter.Get("", handler.ListCats)
	catRouter.Post("", handler.AddCat)
	catRouter.Put("/:"+catIDFromParam, handler.UpdateCat)
	catRouter.Delete("/:"+catIDFromParam, handler.DeleteCat)
}

func (h catHandler) ListCats(c *fiber.Ctx) error {
	callerInfo := "[catHandler.ListCats]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	query, res := &domain.QueryParam{}, baseResponse{}
	if err := c.QueryParser(query); err != nil {
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

	if err := query.Validate(); err != nil {
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

	cats, err := h.catService.ListCats(userCtx, userData.ID, *query)
	if err != nil {
		l.Error("error listing cats",
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

	catsRes := make([]listCatResponse, len(cats))
	for i, cat := range cats {
		catsRes[i] = listCatResponse{
			ID:          cat.ID.String(),
			Name:        cat.Name,
			Race:        cat.Race,
			Sex:         cat.Sex,
			AgeInMonth:  cat.AgeInMonth,
			Description: cat.Description,
			ImageUrls:   cat.ImageUrls,
			HasMatched:  cat.HasMatched,
			CreatedAt:   cat.CreatedAt.Format(time.DateOnly),
		}
	}

	res = baseResponse{
		Message: successListCatMessage,
		Data:    catsRes,
	}
	return c.JSON(res)
}

func (h catHandler) AddCat(c *fiber.Ctx) error {
	callerInfo := "[catHandler.AddCat]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	req, res := &catRequest{}, baseResponse{}
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

	catData := domain.Cat{
		Name:        req.Name,
		Race:        req.Race,
		Sex:         req.Sex,
		AgeInMonth:  req.AgeInMonth,
		Description: req.Description,
		ImageUrls:   req.ImageUrls,
		UserID:      userData.ID,
	}

	cat, err := h.catService.AddCat(userCtx, catData)
	if err != nil {
		l.Error("error adding cat",
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
		Message: successAddCatMessage,
		Data: addCatResponse{
			ID:        cat.ID.String(),
			CreatedAt: cat.CreatedAt.Format(time.DateOnly),
		},
	}

	return c.Status(http.StatusCreated).JSON(res)
}

func (h catHandler) UpdateCat(c *fiber.Ctx) error {
	callerInfo := "[catHandler.UpdateCat]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	catIDParam := c.Params(catIDFromParam)
	if catIDParam == "" {
		l.Error("empty catID")
		res := baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": "catID is required",
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	req, res := &catRequest{}, baseResponse{}
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

	catID, err := ulid.Parse(catIDParam)
	if err != nil {
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

	if err = req.validate(); err != nil {
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

	catData := domain.Cat{
		ID:          catID,
		Name:        req.Name,
		Race:        req.Race,
		Sex:         req.Sex,
		AgeInMonth:  req.AgeInMonth,
		Description: req.Description,
		ImageUrls:   req.ImageUrls,
		UserID:      userData.ID,
	}

	_, err = h.catService.UpdateCat(userCtx, catData)
	switch {
	case errors.Is(err, domain.ErrCatNotFound):
		l.Info("cat not found",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case err != nil:
		l.Error("error updating cat",
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
		Message: successUpdateCatMessage,
	}

	return c.JSON(res)
}

func (h catHandler) DeleteCat(c *fiber.Ctx) error {
	callerInfo := "[catHandler.DeleteCat]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	catIDParam := c.Params(catIDFromParam)
	if catIDParam == "" {
		l.Error("empty catID")
		res := baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": "catID is required",
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	catID, err := ulid.Parse(catIDParam)
	if err != nil {
		l.Error("error validate data",
			zap.Error(err),
		)
		res := baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	catData := domain.Cat{
		ID:     catID,
		UserID: userData.ID,
	}

	err = h.catService.DeleteCat(userCtx, catData)
	switch {
	case errors.Is(err, domain.ErrCatNotFound):
		l.Info("cat not found",
			zap.Error(err),
		)
		res := baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case err != nil:
		l.Error("error deleting cat",
			zap.Error(err),
		)
		res := baseResponse{
			Message: domain.InternalServerErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	res := baseResponse{
		Message: successDeleteCatMessage,
	}

	return c.JSON(res)
}
