package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/logger"
	"cats-social/internal/application/match/service"
	"cats-social/internal/domain"
)

const (
	matchIDFromParam = "matchID"
)

type matchHandler struct {
	matchService service.MatchServiceContract
}

func NewMatchHandler(router fiber.Router, jwtMiddleware fiber.Handler, matchService service.MatchServiceContract) {
	handler := matchHandler{
		matchService: matchService,
	}

	matchRouter := router.Group("/cat/match")

	matchRouter.Use(jwtMiddleware)
	matchRouter.Post("", handler.NewMatch)
	matchRouter.Get("", handler.GetMatch)
	matchRouter.Post("/approve", handler.ApproveMatch)
	matchRouter.Post("/reject", handler.RejectMatch)
	matchRouter.Delete("/:"+matchIDFromParam, handler.DeleteMatch)
}

func (h matchHandler) NewMatch(c *fiber.Ctx) error {
	callerInfo := "[matchHandler.NewMatch]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	req, res := &matchRequest{}, baseResponse{}
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
		l.Error("error validating request",
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

	matchData := domain.Match{
		MatchCatID: req.MatchCatId,
		UserCatID:  req.UserCatId,
		Message:    req.Message,
	}

	_, err := h.matchService.NewMatch(userCtx, matchData, userData.ID)
	switch {
	case errors.Is(err, domain.ErrCatNotFound):
		l.Error("cat not found",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case errors.Is(err, domain.ErrCatGenderNotMatch),
		errors.Is(err, domain.ErrCatAlreadyMatched),
		errors.Is(err, domain.ErrCatSameOwner):
		l.Error(err.Error(),
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)

	case err != nil:
		l.Error("error creating new match",
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
		Message: successMatchMessage,
	}

	return c.Status(http.StatusCreated).JSON(res)
}

func (h matchHandler) GetMatch(c *fiber.Ctx) error {
	callerInfo := "[matchHandler.GetMatch]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	detailMatches, err := h.matchService.GetMatch(userCtx, userData.ID)
	if err != nil {
		l.Error("error getting match",
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

	detailMatchesRes := make([]detailMatchResponse, len(detailMatches))
	for i, detailMatch := range detailMatches {
		detailMatchesRes[i] = detailMatchResponse{
			ID: detailMatch.ID,
			IssuedBy: issuedBy{
				Name:      detailMatch.Issuer.Name,
				Email:     detailMatch.Issuer.Email,
				CreatedAt: detailMatch.Issuer.CreatedAt.Format(time.DateOnly),
			},
			MatchCatDetail: catDetail{
				ID:          detailMatch.MatchCat.ID.String(),
				Name:        detailMatch.MatchCat.Name,
				Race:        detailMatch.MatchCat.Race,
				Sex:         detailMatch.MatchCat.Sex,
				AgeInMonth:  detailMatch.MatchCat.AgeInMonth,
				Description: detailMatch.MatchCat.Description,
				ImageUrls:   detailMatch.MatchCat.ImageUrls,
				HasMatched:  detailMatch.MatchCat.HasMatched,
				CreatedAt:   detailMatch.MatchCat.CreatedAt.Format(time.DateOnly),
			},
			UserCatDetail: catDetail{
				ID:          detailMatch.UserCat.ID.String(),
				Name:        detailMatch.UserCat.Name,
				Race:        detailMatch.UserCat.Race,
				Sex:         detailMatch.UserCat.Sex,
				AgeInMonth:  detailMatch.UserCat.AgeInMonth,
				Description: detailMatch.UserCat.Description,
				ImageUrls:   detailMatch.UserCat.ImageUrls,
				HasMatched:  detailMatch.UserCat.HasMatched,
				CreatedAt:   detailMatch.UserCat.CreatedAt.Format(time.DateOnly),
			},
			Message:   detailMatch.Message,
			CreatedAt: detailMatch.CreatedAt.Format(time.DateOnly),
		}
	}

	res := baseResponse{
		Message: successGetMatchMessage,
		Data:    detailMatchesRes,
	}

	return c.JSON(res)
}

func (h matchHandler) ApproveMatch(c *fiber.Ctx) error {
	callerInfo := "[matchHandler.ApproveMatch]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	req, res := &approvalMatchRequest{}, baseResponse{}
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
		l.Error("error validating request",
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

	err := h.matchService.ApproveMatch(userCtx, req.MatchID, userData.ID)
	switch {
	case errors.Is(err, domain.ErrMatchNotFound):
		l.Error("match not found",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case errors.Is(err, domain.ErrMatchNotValid), errors.Is(err, domain.ErrCatAlreadyMatched):
		l.Error("match not valid",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)

	case err != nil:
		l.Error("error approving match",
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
		Message: successApproveMatchMessage,
	}

	return c.JSON(res)
}

func (h matchHandler) RejectMatch(c *fiber.Ctx) error {
	callerInfo := "[matchHandler.RejectMatch]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	req, res := &approvalMatchRequest{}, baseResponse{}
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
		l.Error("error validating request",
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

	err := h.matchService.RejectMatch(userCtx, req.MatchID, userData.ID)
	switch {
	case errors.Is(err, domain.ErrMatchNotFound):
		l.Error("match not found",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case errors.Is(err, domain.ErrMatchNotValid):
		l.Error("match not valid",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)

	case err != nil:
		l.Error("error approving match",
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
		Message: successRejectMatchMessage,
	}

	return c.JSON(res)
}

func (h matchHandler) DeleteMatch(c *fiber.Ctx) error {
	callerInfo := "[matchHandler.DeleteMatch]"

	userCtx := c.UserContext()
	l := logger.FromCtx(userCtx).With(zap.String("caller", callerInfo))

	userData := c.Locals(domain.UserFromToken).(domain.User)

	matchIDParam := c.Params(matchIDFromParam)
	if matchIDParam == "" {
		l.Error("empty matchID")
		res := baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": "matchID is required",
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)
	}

	matchID, err := ulid.Parse(matchIDParam)
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

	res := baseResponse{}
	err = h.matchService.DeleteMatch(userCtx, matchID, userData.ID)
	switch {
	case errors.Is(err, domain.ErrMatchNotFound):
		l.Error("match not found",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.NotFoundErrorMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusNotFound).JSON(res)

	case errors.Is(err, domain.ErrMatchNotValid), errors.Is(err, domain.ErrCatAlreadyMatched):
		l.Error("match not valid",
			zap.Error(err),
		)
		res = baseResponse{
			Message: domain.InvalidRequestBodyMessage,
			Data: fiber.Map{
				"error": err.Error(),
			},
		}
		return c.Status(http.StatusBadRequest).JSON(res)

	case err != nil:
		l.Error("error approving match",
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
		Message: successDeleteMatchMessage,
	}

	return c.JSON(res)
}
