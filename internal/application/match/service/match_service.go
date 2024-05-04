package service

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/logger"
	catRepo "cats-social/internal/application/cat/repository"
	matchRepo "cats-social/internal/application/match/repository"
	userRepo "cats-social/internal/application/user/repository"
	"cats-social/internal/domain"
)

type MatchService struct {
	matchRepository matchRepo.MatchRepositoryContract
	catRepository   catRepo.CatRepositoryContract
	userRepository  userRepo.AuthRepositoryContract
	contextTimeout  time.Duration
}

func NewMatchService(
	timeout time.Duration,
	matchRepository matchRepo.MatchRepositoryContract,
	catRepository catRepo.CatRepositoryContract,
	userRepository userRepo.AuthRepositoryContract,
) *MatchService {
	matchService := &MatchService{
		matchRepository: matchRepository,
		catRepository:   catRepository,
		userRepository:  userRepository,
		contextTimeout:  timeout,
	}

	return matchService
}

func (m MatchService) NewMatch(ctx context.Context, match domain.Match, userID ulid.ULID) (domain.Match, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	callerInfo := "[MatchService.NewMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// Check matchCatId is exist
	cats, err := m.catRepository.Get(ctx, userID, domain.QueryParam{
		ID:    match.MatchCatID,
		Owned: domain.FalseBool,
	}, false)
	if err != nil {
		l.Error("error get cat", zap.Error(err))
		return match, err
	}

	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Error("error get cat", zap.Error(err))
		return match, err
	}

	matchCat := cats[0]

	// Check userCatId is existed and make sure it's owned by the user
	cats, err = m.catRepository.Get(ctx, userID, domain.QueryParam{
		ID:    match.UserCatID,
		Owned: domain.TrueBool,
	}, false)
	if err != nil {
		l.Error("error get cat", zap.Error(err))
		return match, err
	}

	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Error("error get cat", zap.Error(err))
		return match, err
	}

	userCat := cats[0]

	// Compare gender
	if matchCat.Sex == userCat.Sex {
		err = domain.ErrCatGenderNotMatch
		l.Error("error compare sex of cats", zap.Error(err))
		return match, err
	}

	// Check if both cats already matched
	if matchCat.HasMatched || userCat.HasMatched {
		err = domain.ErrCatAlreadyMatched
		l.Error("error check hasMatched status", zap.Error(err))
		return match, err
	}

	hasMatched, err := m.matchRepository.HasMatched(ctx, match)
	if err != nil {
		l.Error("error find match", zap.Error(err))
		return match, err
	}

	if hasMatched {
		err = domain.ErrCatAlreadyMatched
		l.Error("error check hasMatched status", zap.Error(err))
		return match, err
	}

	// Make sure not the same owner
	if matchCat.UserID == userCat.UserID {
		err = domain.ErrCatSameOwner
		l.Error("error check same owner", zap.Error(err))
		return match, err
	}

	match, err = m.matchRepository.NewMatch(ctx, match)
	if err != nil {
		l.Error("error create match", zap.Error(err))
		return match, err
	}

	return match, nil
}

func (m MatchService) GetMatch(ctx context.Context, userID ulid.ULID) ([]domain.DetailMatch, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	callerInfo := "[MatchService.GetMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	detailMatches := make([]domain.DetailMatch, 0)

	// get matches with user_id as issuer and receiver
	detailMatches, err := m.matchRepository.GetDetailMatches(ctx, userID)
	if err != nil {
		l.Error("error get matches", zap.Error(err))
		return detailMatches, err
	}

	for i, detailMatch := range detailMatches {
		// get user data based on user_id from issuer
		var user domain.User
		user, err = m.userRepository.Get(ctx, detailMatch.Issuer.ID)
		if err != nil {
			l.Error("error get issuer", zap.Error(err))
			return detailMatches, err
		}
		detailMatches[i].Issuer = user

		// get matchCatDetail
		var cats []domain.Cat
		cats, err = m.catRepository.Get(ctx, detailMatch.Receiver.ID, domain.QueryParam{
			ID: detailMatch.MatchCatID,
		}, true)
		if err != nil {
			l.Error("error get match cat", zap.Error(err))
			return detailMatches, err
		}
		if len(cats) != 1 {
			err = domain.ErrCatNotFound
			l.Error("error get match cat", zap.Error(err))
			return detailMatches, err
		}
		detailMatches[i].MatchCat = cats[0]

		// get userCatDetail
		cats, err = m.catRepository.Get(ctx, detailMatch.Issuer.ID, domain.QueryParam{
			ID: detailMatch.UserCatID,
		}, true)
		if err != nil {
			l.Error("error get user cat", zap.Error(err))
			return detailMatches, err
		}
		if len(cats) != 1 {
			err = domain.ErrCatNotFound
			l.Error("error get user cat", zap.Error(err))
			return detailMatches, err
		}
		detailMatches[i].UserCat = cats[0]
	}

	return detailMatches, nil
}

func (m MatchService) ApproveMatch(ctx context.Context, matchID, userID ulid.ULID) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	callerInfo := "[MatchService.ApproveMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// check matchID is existed
	detailMatch, err := m.matchRepository.Get(ctx, matchID)
	if err != nil {
		return err
	}

	// check if receiver is the user and matchID valid
	if detailMatch.Receiver.ID != userID {
		err = domain.ErrMatchNotFound
		l.Error("error check match", zap.Error(err))
		return err
	}

	if !detailMatch.Match.DeletedAt.IsZero() {
		err = domain.ErrMatchNotValid
		l.Error("error check match", zap.Error(err))
		return err
	}

	// make sure both cats are not matched
	var cats []domain.Cat
	cats, err = m.catRepository.Get(ctx, detailMatch.Receiver.ID, domain.QueryParam{
		ID: detailMatch.MatchCatID,
	}, false)
	if err != nil {
		l.Error("error get match cat", zap.Error(err))
		return err
	}
	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Error("error get match cat", zap.Error(err))
		return err
	}
	receiverCat := cats[0]

	cats, err = m.catRepository.Get(ctx, detailMatch.Issuer.ID, domain.QueryParam{
		ID: detailMatch.UserCatID,
	}, false)
	if err != nil {
		l.Error("error get user cat", zap.Error(err))
		return err
	}
	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Error("error get user cat", zap.Error(err))
		return err
	}
	issuerCat := cats[0]

	if receiverCat.HasMatched || issuerCat.HasMatched {
		err = domain.ErrCatAlreadyMatched
		l.Error("error check hasMatched status", zap.Error(err))
		return err
	}

	tx, err := m.matchRepository.TxBegin(ctx)
	if err != nil {
		l.Error("error begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Update matchCatId and userCatId hasMatched status
	receiverCat.HasMatched = true
	receiverCat, tx, err = m.catRepository.Update(ctx, receiverCat, tx)
	if err != nil {
		l.Error("error update match cat", zap.Error(err))
		return err
	}

	issuerCat.HasMatched = true
	issuerCat, tx, err = m.catRepository.Update(ctx, issuerCat, tx)
	if err != nil {
		l.Error("error update user cat", zap.Error(err))
		return err
	}

	// Update match.deleted_at where userID is the issuer and the receiver except the matchID
	tx, err = m.matchRepository.DeleteExceptApproved(ctx, userID, matchID, tx)
	if err != nil {
		l.Error("error delete match", zap.Error(err))
		return err
	}

	err = m.matchRepository.TxCommit(ctx, tx)
	if err != nil {
		l.Error("error commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (m MatchService) RejectMatch(ctx context.Context, matchID, userID ulid.ULID) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	callerInfo := "[MatchService.RejectMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// check matchID is existed
	detailMatch, err := m.matchRepository.Get(ctx, matchID)
	if err != nil {
		return err
	}

	// check if receiver is the user and matchID valid
	if detailMatch.Receiver.ID != userID {
		err = domain.ErrMatchNotFound
		l.Error("error check match", zap.Error(err))
		return err
	}

	if !detailMatch.Match.DeletedAt.IsZero() {
		err = domain.ErrMatchNotValid
		l.Error("error check match", zap.Error(err))
		return err
	}

	err = m.matchRepository.Delete(ctx, matchID)
	if err != nil {
		l.Error("error delete match", zap.Error(err))
		return err
	}

	return nil
}

func (m MatchService) DeleteMatch(ctx context.Context, matchID, userID ulid.ULID) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	callerInfo := "[MatchService.DeleteMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	// check matchID is existed
	detailMatch, err := m.matchRepository.Get(ctx, matchID)
	if err != nil {
		return err
	}

	// check if issuer is the user and matchID valid
	if detailMatch.Issuer.ID != userID {
		err = domain.ErrMatchNotFound
		l.Error("error check match", zap.Error(err))
		return err
	}

	if !detailMatch.Match.DeletedAt.IsZero() {
		err = domain.ErrMatchNotValid
		l.Error("error check match", zap.Error(err))
		return err
	}

	var cats []domain.Cat
	cats, err = m.catRepository.Get(ctx, detailMatch.Issuer.ID, domain.QueryParam{
		ID: detailMatch.UserCatID,
	}, false)
	if err != nil {
		l.Error("error get user cat", zap.Error(err))
		return err
	}
	if len(cats) != 1 {
		err = domain.ErrCatNotFound
		l.Error("error get user cat", zap.Error(err))
		return err
	}
	issuerCat := cats[0]

	if issuerCat.HasMatched {
		err = domain.ErrCatAlreadyMatched
		l.Error("error check hasMatched status", zap.Error(err))
		return err
	}

	err = m.matchRepository.Delete(ctx, matchID)
	if err != nil {
		l.Error("error delete match", zap.Error(err))
		return err
	}

	return nil
}

var _ MatchServiceContract = (*MatchService)(nil)
