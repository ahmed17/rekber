package services

import "rekber/internal/repositories"

// Service container untuk semua service
type Service struct {
    User        UserService
    Auth        AuthService
    Transaction TransactionService
    Testimony   TestimonyService
}

var service *Service

// InitService menginisialisasi semua service dengan repository
func InitService(repo *repositories.Repository) {
    authService := NewAuthService()
    
    service = &Service{
        User:        NewUserService(repo.User, authService),
        Auth:        authService,
        Transaction: NewTransactionService(repo.Transaction, repo.User),
        Testimony:   NewTestimonyService(repo.Testimony, repo.Transaction, repo.User),
    }
}

// GetService mengembalikan singleton service instance
func GetService() *Service {
    return service
}