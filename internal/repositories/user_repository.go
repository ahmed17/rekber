package repositories

import (
    "errors"
    "gorm.io/gorm"
    "rekber/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    List(offset, limit int) ([]models.User, error)
    Count() (int64, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// ... [implementasi method tetap sama seperti sebelumnya]
// Pastikan semua method di bawah ini ada:
func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User
    err := r.db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(offset, limit int) ([]models.User, error) {
    var users []models.User
    err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
    return users, err
}

func (r *userRepository) Count() (int64, error) {
    var count int64
    err := r.db.Model(&models.User{}).Count(&count).Error
    return count, err
}