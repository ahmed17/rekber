package repositories

import "gorm.io/gorm"

// Repository container untuk semua repository
type Repository struct {
    User        UserRepository
    Transaction TransactionRepository
    Testimony   TestimonyRepository
}

var repo *Repository

// InitRepository menginisialisasi semua repository dengan database connection
func InitRepository(db *gorm.DB) {
    repo = &Repository{
        User:        NewUserRepository(db),
        Transaction: NewTransactionRepository(db),
        Testimony:   NewTestimonyRepository(db),
    }
}

// GetRepository mengembalikan singleton repository instance
func GetRepository() *Repository {
    return repo
}