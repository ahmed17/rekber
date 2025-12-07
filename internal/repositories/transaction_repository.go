package repositories

import (
	"errors"
	"gorm.io/gorm"
	"rekber/internal/models"
	"time"
	"log" // Tambahkan ini
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByCode(code string) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uint) error
	ListByUserID(userID uint, role string, offset, limit int) ([]models.Transaction, error)
	ListByStatus(status models.TransactionStatus, offset, limit int) ([]models.Transaction, error)
	ListAll(offset, limit int) ([]models.Transaction, error)
	CountByUserID(userID uint, role string) (int64, error)
	CountByStatus(status models.TransactionStatus) (int64, error)
	CountAll() (int64, error)
	// Untuk scheduler: transaksi yang sudah shipped dan auto_complete_at sudah lewat
	FindShippedWithExpiredAutoComplete() ([]models.Transaction, error)
	// Update status dan set auto_complete_at
	UpdateStatusAndAutoComplete(id uint, status models.TransactionStatus, autoCompleteAt *time.Time) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
    // Debug log
    log.Printf("Creating transaction with BuyerID: %d, SellerID: %d", 
        transaction.BuyerID, transaction.SellerID)
    
    // Validasi foreign key exist
    var buyerCount, sellerCount int64
    r.db.Model(&models.User{}).Where("id = ?", transaction.BuyerID).Count(&buyerCount)
    r.db.Model(&models.User{}).Where("id = ?", transaction.SellerID).Count(&sellerCount)
    
    if buyerCount == 0 {
        return errors.New("buyer not found")
    }
    if sellerCount == 0 {
        return errors.New("seller not found")
    }
    
    return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Admin").First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByCode(code string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Admin").Where("transaction_code = ?", code).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

func (r *transactionRepository) ListByUserID(userID uint, role string, offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Preload("Buyer").Preload("Seller").Preload("Admin").Offset(offset).Limit(limit).Order("created_at DESC")

	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	} else {
		// Jika role admin, tampilkan semua atau berdasarkan userID?
		// Untuk sekarang, kita kembalikan transaksi yang berhubungan dengan user baik sebagai buyer atau seller
		query = query.Where("buyer_id = ? OR seller_id = ?", userID, userID)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) ListByStatus(status models.TransactionStatus, offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Admin").
		Where("status = ?", status).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) ListAll(offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Buyer").Preload("Seller").Preload("Admin").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) CountByUserID(userID uint, role string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Transaction{})

	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "seller" {
		query = query.Where("seller_id = ?", userID)
	} else {
		query = query.Where("buyer_id = ? OR seller_id = ?", userID, userID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *transactionRepository) CountByStatus(status models.TransactionStatus) (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *transactionRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Count(&count).Error
	return count, err
}

func (r *transactionRepository) FindShippedWithExpiredAutoComplete() ([]models.Transaction, error) {
	var transactions []models.Transaction
	now := time.Now()
	err := r.db.Where("status = ? AND auto_complete_at <= ?", models.StatusShipped, now).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) UpdateStatusAndAutoComplete(id uint, status models.TransactionStatus, autoCompleteAt *time.Time) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           status,
		"auto_complete_at": autoCompleteAt,
		"updated_at":       time.Now(),
	}).Error
}
