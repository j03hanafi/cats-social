package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/id"
	"cats-social/common/logger"
	"cats-social/internal/domain"
)

type CatRepository struct {
	db *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) *CatRepository {
	return &CatRepository{
		db: db,
	}
}

func (c CatRepository) Create(ctx context.Context, dCat domain.Cat) (domain.Cat, error) {
	callerInfo := "[CatRepository.Create]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := c.db.Begin(ctx)
	if err != nil {
		l.Error("failed to begin transaction", zap.Error(err))
		return dCat, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mCat := cat{
		ID:          id.New(),
		Name:        dCat.Name,
		Race:        string(dCat.Race),
		Sex:         string(dCat.Sex),
		AgeInMonth:  dCat.AgeInMonth,
		Description: dCat.Description,
		UserID:      dCat.UserID,
		HasMatched:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	err = c.insertCat(ctx, tx, mCat)
	if err != nil {
		l.Error("failed to insert cat", zap.Error(err))
		return dCat, err
	}

	mCatImages := make([]catImages, len(dCat.ImageUrls))
	for i, imageUrl := range dCat.ImageUrls {
		mCatImages[i] = catImages{
			ID:        id.New(),
			ImageURL:  imageUrl,
			CatID:     mCat.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{
				Valid: false,
			},
		}
	}

	err = c.insertCatImages(ctx, tx, mCatImages)
	if err != nil {
		l.Error("failed to insert cat images", zap.Error(err))
		return dCat, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return dCat, err
	}

	dCat.ID = mCat.ID
	dCat.CreatedAt = mCat.CreatedAt
	return dCat, nil
}

func (c CatRepository) insertCat(ctx context.Context, tx pgx.Tx, mCat cat) error {
	callerInfo := "[CatRepository.insertCat]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	insertQuery := `INSERT INTO cats (id, name, race, sex, age_in_month, description, user_id, has_matched, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := tx.Exec(
		ctx,
		insertQuery,
		mCat.ID,
		mCat.Name,
		mCat.Race,
		mCat.Sex,
		mCat.AgeInMonth,
		mCat.Description,
		mCat.UserID,
		mCat.HasMatched,
		mCat.CreatedAt,
		mCat.UpdatedAt,
		mCat.DeletedAt,
	)
	if err != nil {
		l.Error("failed to execute insert query", zap.Error(err))
		return err
	}

	return nil
}

func (c CatRepository) insertCatImages(ctx context.Context, tx pgx.Tx, images []catImages) error {
	callerInfo := "[CatRepository.insertCatImages]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	if len(images) == 0 {
		l.Error("images is empty")
		return errors.New("images is empty")
	}

	rows := make([][]any, len(images))
	for i, image := range images {
		rows[i] = []any{
			image.ID,
			image.ImageURL,
			image.CatID,
			image.CreatedAt,
			image.UpdatedAt,
			image.DeletedAt,
		}
	}

	tableName := pgx.Identifier{images[0].tableName()}
	columns := []string{"id", "image_url", "cat_id", "created_at", "updated_at", "deleted_at"}

	_, err := tx.CopyFrom(ctx, tableName, columns, pgx.CopyFromRows(rows))
	if err != nil {
		l.Error("failed to copy from", zap.Error(err))
		return err
	}

	return nil
}

func (c CatRepository) Get(
	ctx context.Context,
	userID ulid.ULID,
	query domain.QueryParam,
	withImages bool,
) ([]domain.Cat, error) {
	callerInfo := "[CatRepository.Get]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	getCatsQuery := `SELECT id, name, race, sex, age_in_month, description, user_id, has_matched, created_at, updated_at, deleted_at FROM cats`
	getCatsQuery, params := c.getConditions(getCatsQuery, query, userID)

	rows, err := c.db.Query(ctx, getCatsQuery, params...)
	if err != nil {
		l.Error("failed to query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	cats := make([]domain.Cat, 0)
	for rows.Next() {
		var mCat cat
		err = rows.Scan(
			&mCat.ID,
			&mCat.Name,
			&mCat.Race,
			&mCat.Sex,
			&mCat.AgeInMonth,
			&mCat.Description,
			&mCat.UserID,
			&mCat.HasMatched,
			&mCat.CreatedAt,
			&mCat.UpdatedAt,
			&mCat.DeletedAt,
		)
		if err != nil {
			l.Error("failed to scan cat", zap.Error(err))
			return nil, err
		}

		cats = append(cats, domain.Cat{
			ID:          mCat.ID,
			Name:        mCat.Name,
			Race:        domain.CatRace(mCat.Race),
			Sex:         domain.CatSex(mCat.Sex),
			AgeInMonth:  mCat.AgeInMonth,
			Description: mCat.Description,
			UserID:      mCat.UserID,
			HasMatched:  mCat.HasMatched,
			CreatedAt:   mCat.CreatedAt,
		})
	}

	if err = rows.Err(); err != nil {
		l.Error("failed to scan cat", zap.Error(err))
		return nil, err
	}

	if !withImages {
		return cats, nil
	}

	err = c.getImages(ctx, cats)
	if err != nil {
		l.Error("failed to get images", zap.Error(err))
		return nil, err
	}

	return cats, nil
}

func (c CatRepository) getConditions(getQuery string, queryParam domain.QueryParam, userID ulid.ULID) (string, []any) {
	params := make([]any, 0)
	conditions := make([]string, 0)

	var emptyID ulid.ULID
	if queryParam.ID != emptyID {
		params = append(params, queryParam.ID)
		conditions = append(conditions, fmt.Sprintf("id = $%d", len(params)))
	}

	if queryParam.Race != "" {
		params = append(params, queryParam.Race)
		conditions = append(conditions, fmt.Sprintf("race = $%d", len(params)))
	}

	if queryParam.Sex != "" {
		params = append(params, queryParam.Sex)
		conditions = append(conditions, fmt.Sprintf("sex = $%d", len(params)))
	}

	if queryParam.HasMatched != "" {
		params = append(params, queryParam.HasMatched)
		conditions = append(conditions, fmt.Sprintf("has_matched = $%d", len(params)))
	}

	if queryParam.AgeInMonth != "" && len(queryParam.AgeInMonth) > 1 {
		comparison := queryParam.AgeInMonth[:1]
		if comparison != ">" && comparison != "<" && comparison != "=" {
			comparison = "="
		}

		params = append(params, queryParam.AgeInMonth[1:])
		conditions = append(conditions, fmt.Sprintf("age_in_month %s $%d", comparison, len(params)))
	}

	if queryParam.Owned != "" {
		params = append(params, userID)

		var ownedCondition string
		if queryParam.Owned == domain.TrueBool {
			ownedCondition = "="
		} else if queryParam.Owned == domain.FalseBool {
			ownedCondition = "!="
		}

		conditions = append(conditions, fmt.Sprintf("user_id %s $%d", ownedCondition, len(params)))
	}

	if queryParam.Search != "" {
		params = append(params, fmt.Sprintf("%%%s%%", queryParam.Search))
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(params)))
	}

	conditions = append(conditions, "deleted_at IS NULL")

	filter := make([]string, 0)
	filter = append(filter, "ORDER BY created_at DESC")

	if queryParam.Limit != 0 {
		params = append(params, queryParam.Limit)
		filter = append(filter, fmt.Sprintf("LIMIT $%d", len(params)))
	}

	if queryParam.Offset != 0 {
		params = append(params, queryParam.Offset)
		filter = append(filter, fmt.Sprintf("OFFSET $%d", len(params)))
	}

	if len(conditions) > 0 {
		getQuery = fmt.Sprintf("%s WHERE %s", getQuery, strings.Join(conditions, " AND "))
	}

	if len(filter) > 0 {
		getQuery = fmt.Sprintf("%s %s", getQuery, strings.Join(filter, " "))
	}

	return getQuery, params
}

func (c CatRepository) getImages(ctx context.Context, cats []domain.Cat) (err error) {
	callerInfo := "[CatRepository.getImages]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	getCatImagesQuery := `SELECT image_url FROM cat_images WHERE cat_id = $1`
	batch := &pgx.Batch{}

	for i, mCat := range cats {
		batch.Queue(getCatImagesQuery, mCat.ID).Query(func(rows pgx.Rows) error {
			defer rows.Close()

			imageUrls := make([]string, 0)
			for rows.Next() {
				var imageUrl string
				err = rows.Scan(&imageUrl)
				if err != nil {
					l.Error("failed to scan image url", zap.Error(err))
					return err
				}

				imageUrls = append(imageUrls, imageUrl)
			}

			cats[i].ImageUrls = imageUrls
			return rows.Err()
		})
	}

	err = c.db.SendBatch(ctx, batch).Close()
	if err != nil {
		l.Error("failed to send batch", zap.Error(err))
		return err
	}

	return nil
}

func (c CatRepository) Update(ctx context.Context, dCat domain.Cat) (domain.Cat, error) {
	callerInfo := "[CatRepository.Update]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := c.db.Begin(ctx)
	if err != nil {
		l.Error("failed to begin transaction", zap.Error(err))
		return dCat, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mCat := cat{
		ID:          dCat.ID,
		Name:        dCat.Name,
		Race:        string(dCat.Race),
		Sex:         string(dCat.Sex),
		AgeInMonth:  dCat.AgeInMonth,
		Description: dCat.Description,
		HasMatched:  dCat.HasMatched,
		UpdatedAt:   time.Now(),
	}

	err = c.updateCat(ctx, tx, mCat)
	if err != nil {
		l.Error("failed to update cat", zap.Error(err))
		return dCat, err
	}

	if len(dCat.ImageUrls) != 0 {
		mCatImages := make([]catImages, len(dCat.ImageUrls))
		for i, imageUrl := range dCat.ImageUrls {
			mCatImages[i] = catImages{
				ID:        id.New(),
				ImageURL:  imageUrl,
				CatID:     mCat.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: sql.NullTime{
					Valid: false,
				},
			}
		}

		err = c.updateCatImages(ctx, tx, mCatImages)
		if err != nil {
			l.Error("failed to update cat images", zap.Error(err))
			return dCat, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return dCat, err
	}

	return dCat, nil
}

func (c CatRepository) updateCat(ctx context.Context, tx pgx.Tx, mCat cat) error {
	callerInfo := "[CatRepository.updateCat]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	updateQuery := `UPDATE cats SET name = $1, race = $2, sex = $3, age_in_month = $4, description = $5, has_matched = $6, updated_at = $7 WHERE id = $8`

	_, err := tx.Exec(
		ctx,
		updateQuery,
		mCat.Name,
		mCat.Race,
		mCat.Sex,
		mCat.AgeInMonth,
		mCat.Description,
		mCat.HasMatched,
		mCat.UpdatedAt,
		mCat.ID,
	)
	if err != nil {
		l.Error("failed to execute update query", zap.Error(err))
		return err

	}

	return nil
}

func (c CatRepository) updateCatImages(ctx context.Context, tx pgx.Tx, images []catImages) error {
	callerInfo := "[CatRepository.insertCatImages]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	if len(images) == 0 {
		l.Error("images is empty")
		return errors.New("images is empty")
	}

	err := c.deleteCatImages(ctx, tx, images[0].CatID)
	if err != nil {
		l.Error("failed to delete cat images", zap.Error(err))
		return err
	}

	err = c.insertCatImages(ctx, tx, images)
	if err != nil {
		l.Error("failed to insert cat images", zap.Error(err))
		return err
	}

	return nil
}

func (c CatRepository) deleteCatImages(ctx context.Context, tx pgx.Tx, catID ulid.ULID) error {
	callerInfo := "[CatRepository.deleteCatImages]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	deleteQuery := `DELETE FROM cat_images WHERE cat_id = $1`

	_, err := tx.Exec(ctx, deleteQuery, catID)
	if err != nil {
		l.Error("failed to execute delete query", zap.Error(err))
		return err
	}

	return nil
}

func (c CatRepository) Delete(ctx context.Context, catID ulid.ULID) error {
	callerInfo := "[CatRepository.Delete]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := c.db.Begin(ctx)
	if err != nil {
		l.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deleteQuery := `UPDATE cats SET deleted_at = $1 WHERE id = $2`

	_, err = tx.Exec(ctx, deleteQuery, time.Now(), catID)
	if err != nil {
		l.Error("failed to delete cat", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

var _ CatRepositoryContract = (*CatRepository)(nil)
