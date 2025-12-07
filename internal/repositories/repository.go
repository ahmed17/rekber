package repositories

import "gorm.io/gorm"

type Repository struct {
	User        UserRepository
	Transaction TransactionRepository
	Testimony   TestimonyRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:        NewUserRepository(db),
		Transaction: NewTransactionRepository(db),
		Testimony:   NewTestimonyRepository(db),
	}
}
