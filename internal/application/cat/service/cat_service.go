package service

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/logger"
	catRepo "cats-social/internal/application/cat/repository"
	matchRepo "cats-social/internal/application/match/repository"
	"cats-social/internal/domain"
)

type CatService struct {
	catRepository   catRepo.CatRepositoryContract
	matchRepository matchRepo.MatchRepositoryContract
	contextTimeout  time.Duration
}

func NewCatService(
	timeout time.Duration,
	catRepository catRepo.CatRepositoryContract,
	matchRepository matchRepo.MatchRepositoryContract,
) *CatService {
	catService := &CatService{
		catRepository:   catRepository,
		matchRepository: matchRepository,
		contextTimeout:  timeout,
	}

	return catService
}

func (c CatService) AddCat(ctx context.Context, cat domain.Cat) (domain.Cat, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	callerInfo := "[CatService.AddCat]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	cat, err := c.catRepository.Create(ctx, cat)
	if err != nil {
		l.Error("error add cat", zap.Error(err))
		return cat, err
	}

	return cat, nil
}

func (c CatService) ListCats(ctx context.Context, userID ulid.ULID, query domain.QueryParam) ([]domain.Cat, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	callerInfo := "[CatService.ListCats]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	cats, err := c.catRepository.Get(ctx, userID, query, true)
	if err != nil {
		l.Error("error list cats", zap.Error(err))
		return cats, err
	}

	return cats, nil
}

func (c CatService) UpdateCat(ctx context.Context, updatedCat domain.Cat) (domain.Cat, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	callerInfo := "[CatService.UpdateCat]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// Check for ID
	cats, err := c.catRepository.Get(ctx, updatedCat.UserID, domain.QueryParam{
		ID:    updatedCat.ID,
		Owned: domain.TrueBool,
	}, false)
	if err != nil {
		l.Error("error get cat", zap.Error(err))
		return updatedCat, err
	}

	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Info("error get cat", zap.Error(err))
		return updatedCat, err
	}

	cat := cats[0]

	if cat.Sex != updatedCat.Sex {
		foundMatches := make([]domain.DetailMatch, 0)
		foundMatches, err = c.matchRepository.GetDetailMatches(ctx, updatedCat.UserID)
		if err != nil {
			l.Error("error get matches", zap.Error(err))
			return cat, err
		}
		if len(foundMatches) > 0 {
			err = domain.ErrCatAlreadyMatched
			l.Error("cat already requested to match", zap.Error(domain.ErrCatAlreadyMatched))
			return cat, err
		}
	}

	cat.Name = updatedCat.Name
	cat.Race = updatedCat.Race
	cat.Sex = updatedCat.Sex
	cat.AgeInMonth = updatedCat.AgeInMonth
	cat.Description = updatedCat.Description
	cat.ImageUrls = updatedCat.ImageUrls

	foundMatches, err := c.matchRepository.GetDetailMatches(ctx, updatedCat.UserID)
	if err != nil {
		l.Error("error get matches", zap.Error(err))
		return cat, err
	}
	if len(foundMatches) > 0 {
		err = domain.ErrCatAlreadyMatched
		l.Error("cat already requested to match", zap.Error(domain.ErrCatAlreadyMatched))
		return cat, err
	}

	cat, _, err = c.catRepository.Update(ctx, cat)
	if err != nil {
		l.Error("error update cat", zap.Error(err))
		return cat, err
	}

	return cat, nil
}

func (c CatService) DeleteCat(ctx context.Context, cat domain.Cat) error {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	callerInfo := "[CatService.DeleteCat]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// Check for ID
	cats, err := c.catRepository.Get(ctx, cat.UserID, domain.QueryParam{
		ID:    cat.ID,
		Owned: domain.TrueBool,
	}, false)
	if err != nil {
		l.Error("error get cat", zap.Error(err))
		return err
	}

	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Info("error get cat", zap.Error(err))
		return err
	}

	err = c.catRepository.Delete(ctx, cat.ID)
	if err != nil {
		l.Error("error delete cat", zap.Error(err))
		return err
	}

	return nil
}

var _ CatServiceContract = (*CatService)(nil)
