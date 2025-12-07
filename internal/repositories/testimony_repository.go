package repositories

import (
	"errors"
	"gorm.io/gorm"
	"rekber/internal/models"
)

type TestimonyRepository interface {
	Create(testimony *models.Testimony) error
	FindByID(id uint) (*models.Testimony, error)
	FindByTransactionID(transactionID uint) (*models.Testimony, error)
	Update(testimony *models.Testimony) error
	Delete(id uint) error
	ListByUserID(userID uint, offset, limit int) ([]models.Testimony, error)
	ListAll(offset, limit int) ([]models.Testimony, error)
	CountByUserID(userID uint) (int64, error)
	CountAll() (int64, error)
}

type testimonyRepository struct {
	db *gorm.DB
}

func NewTestimonyRepository(db *gorm.DB) TestimonyRepository {
	return &testimonyRepository{db: db}
}

func (r *testimonyRepository) Create(testimony *models.Testimony) error {
	return r.db.Create(testimony).Error
}

func (r *testimonyRepository) FindByID(id uint) (*models.Testimony, error) {
	var testimony models.Testimony
	err := r.db.Preload("Transaction").Preload("User").First(&testimony, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &testimony, nil
}

func (r *testimonyRepository) FindByTransactionID(transactionID uint) (*models.Testimony, error) {
	var testimony models.Testimony
	err := r.db.Preload("Transaction").Preload("User").Where("transaction_id = ?", transactionID).First(&testimony).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &testimony, nil
}

func (r *testimonyRepository) Update(testimony *models.Testimony) error {
	return r.db.Save(testimony).Error
}

func (r *testimonyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Testimony{}, id).Error
}

func (r *testimonyRepository) ListByUserID(userID uint, offset, limit int) ([]models.Testimony, error) {
	var testimonies []models.Testimony
	err := r.db.Preload("Transaction").Preload("User").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&testimonies).Error
	return testimonies, err
}

func (r *testimonyRepository) ListAll(offset, limit int) ([]models.Testimony, error) {
	var testimonies []models.Testimony
	err := r.db.Preload("Transaction").Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&testimonies).Error
	return testimonies, err
}

func (r *testimonyRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Testimony{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *testimonyRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Testimony{}).Count(&count).Error
	return count, err
}
