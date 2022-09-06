package refsessions

import (
	"context"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

// Create создаёт новую сессию.
func (r *Repository) Create(
	ctx context.Context,
	uid models.UserId,
	refToken string,
	accToken string,
	fp string,
	expAt time.Time,
) (*models.RefreshSession, error) {
	sess := &models.RefreshSession{
		UserId:       uid,
		RefreshToken: refToken,
		AccessToken:  accToken,
		Fingerprint:  fp,
		ExpiresAt:    expAt,
		CreatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(sess).Error; err != nil {
		return nil, err
	}
	return sess, nil
}

// DeleteByIds удаляет записи по их идентификаторам.
func (r *Repository) DeleteByIds(ctx context.Context, ids []models.RefreshSessionId) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	res := r.
		db.
		WithContext(ctx).
		Where("id IN ?", ids).
		Delete(models.RefreshSession{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// DeleteByUserId удаляет записи по идентификатору пользователя.
func (r *Repository) DeleteByUserId(ctx context.Context, uid models.UserId) (int64, error) {
	res := r.db.WithContext(ctx).Delete(models.RefreshSession{}, "user_id = ?", uid)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// FindByRefreshToken возвращает сессию по её токену обновления.
func (r *Repository) FindByRefreshToken(ctx context.Context, token string) (*models.RefreshSession, error) {
	var res []models.RefreshSession

	if err := r.
		db.
		WithContext(ctx).
		Where("refresh_token = ?", token).
		Limit(1).
		Find(&res).
		Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

// FindByUserId возвращает все сессии указанного пользователя.
func (r *Repository) FindByUserId(ctx context.Context, uid models.UserId) ([]*models.RefreshSession, error) {
	var res []*models.RefreshSession

	if err := r.
		db.
		WithContext(ctx).
		Where("user_id = ?", uid).
		Find(&res).
		Error; err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
}

// RefreshById обновляет информацию о сессии по её идентификатору.
func (r *Repository) RefreshById(
	ctx context.Context,
	id models.RefreshSessionId,
	refToken string,
	refExpiresAt time.Time,
	accToken string,
) error {
	return r.
		db.
		WithContext(ctx).
		Where("id = ?", id).
		Select("RefreshToken", "AccessToken", "ExpiresAt", "CreatedAt").
		Updates(models.RefreshSession{
			RefreshToken: refToken,
			AccessToken:  accToken,
			ExpiresAt:    refExpiresAt,
			CreatedAt:    time.Now(),
		}).
		Error
}

// NewRepository создает новый экземпляр Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
