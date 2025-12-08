package services

import (
    "errors"
    "time"

    "rekber/internal/models"
    "rekber/internal/repositories"
)

// Tipe lokal untuk testimony service
type CreateTestimonyRequest struct {
    TransactionID uint   `json:"transaction_id" binding:"required"`
    Rating        int    `json:"rating" binding:"required,min=1,max=5"`
    Comment       string `json:"comment" binding:"required,min=1,max=500"`
}

type UpdateTestimonyRequest struct {
    Rating  int    `json:"rating" binding:"omitempty,min=1,max=5"`
    Comment string `json:"comment" binding:"omitempty,min=1,max=500"`
}

type TestimonyResponse struct {
    ID            uint            `json:"id"`
    TransactionID uint            `json:"transaction_id"`
    UserID        uint            `json:"user_id"`
    User          UserResponse    `json:"user"`
    Rating        int             `json:"rating"`
    Comment       string          `json:"comment"`
    AdminNotes    string          `json:"admin_notes,omitempty"`
    CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
}

// TestimonyService interface
type TestimonyService interface {
    CreateTestimony(req CreateTestimonyRequest, userID uint) (*TestimonyResponse, error)
    GetTestimonyByID(testimonyID uint) (*TestimonyResponse, error)
    GetTestimonyByTransactionID(transactionID uint) (*TestimonyResponse, error)
    UpdateTestimony(testimonyID uint, req UpdateTestimonyRequest, userID uint) (*TestimonyResponse, error)
    DeleteTestimony(testimonyID uint, userID uint) error
    ListTestimonies(page, limit int) ([]TestimonyResponse, int64, error)
    ListUserTestimonies(userID uint, page, limit int) ([]TestimonyResponse, int64, error)
}

type testimonyService struct {
    testimonyRepo   repositories.TestimonyRepository
    transactionRepo repositories.TransactionRepository
    userRepo        repositories.UserRepository
}

func NewTestimonyService(testimonyRepo repositories.TestimonyRepository, transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository) TestimonyService {
    return &testimonyService{
        testimonyRepo:   testimonyRepo,
        transactionRepo: transactionRepo,
        userRepo:        userRepo,
    }
}

func (s *testimonyService) CreateTestimony(req CreateTestimonyRequest, userID uint) (*TestimonyResponse, error) {
    // Check if transaction exists and is completed
    transaction, err := s.transactionRepo.FindByID(req.TransactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Check if transaction is completed
    if transaction.Status != models.StatusCompleted {
        return nil, errors.New("can only create testimony for completed transactions")
    }

    // Check if user is buyer or seller of the transaction
    if transaction.BuyerID != userID && transaction.SellerID != userID {
        return nil, errors.New("you are not part of this transaction")
    }

    // Check if testimony already exists for this transaction by this user
    existingTestimony, err := s.testimonyRepo.FindByTransactionID(req.TransactionID)
    if err != nil {
        return nil, err
    }
    if existingTestimony != nil && existingTestimony.UserID == userID {
        return nil, errors.New("you have already created testimony for this transaction")
    }

    // Validate rating
    if req.Rating < 1 || req.Rating > 5 {
        return nil, errors.New("rating must be between 1 and 5")
    }

    // Create testimony
    testimony := &models.Testimony{
        TransactionID: req.TransactionID,
        UserID:        userID,
        Rating:        req.Rating,
        Comment:       req.Comment,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    if err := s.testimonyRepo.Create(testimony); err != nil {
        return nil, err
    }

    // Get user for response
    user, err := s.userRepo.FindByID(userID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    // Convert to response
    response := &TestimonyResponse{
        ID:            testimony.ID,
        TransactionID: testimony.TransactionID,
        UserID:        testimony.UserID,
        User: UserResponse{
            ID:            user.ID,
            Username:      user.Username,
            Email:         user.Email,
            FullName:      user.FullName,
            Phone:         user.Phone,
            AvatarURL:     user.AvatarURL,
            Role:          user.Role,
            IsVerified:    user.IsVerified,
            IsActive:      user.IsActive,
            BankName:      user.BankName,
            AccountNumber: user.AccountNumber,
            AccountHolder: user.AccountHolder,
            CreatedAt:     user.CreatedAt,
            UpdatedAt:     user.UpdatedAt,
        },
        Rating:     testimony.Rating,
        Comment:    testimony.Comment,
        AdminNotes: testimony.AdminNotes,
        CreatedAt:  testimony.CreatedAt,
        UpdatedAt:  testimony.UpdatedAt,
    }

    return response, nil
}

func (s *testimonyService) GetTestimonyByID(testimonyID uint) (*TestimonyResponse, error) {
    testimony, err := s.testimonyRepo.FindByID(testimonyID)
    if err != nil {
        return nil, err
    }
    if testimony == nil {
        return nil, errors.New("testimony not found")
    }

    // Get user
    user, err := s.userRepo.FindByID(testimony.UserID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    response := &TestimonyResponse{
        ID:            testimony.ID,
        TransactionID: testimony.TransactionID,
        UserID:        testimony.UserID,
        User: UserResponse{
            ID:            user.ID,
            Username:      user.Username,
            Email:         user.Email,
            FullName:      user.FullName,
            Phone:         user.Phone,
            AvatarURL:     user.AvatarURL,
            Role:          user.Role,
            IsVerified:    user.IsVerified,
            IsActive:      user.IsActive,
            BankName:      user.BankName,
            AccountNumber: user.AccountNumber,
            AccountHolder: user.AccountHolder,
            CreatedAt:     user.CreatedAt,
            UpdatedAt:     user.UpdatedAt,
        },
        Rating:     testimony.Rating,
        Comment:    testimony.Comment,
        AdminNotes: testimony.AdminNotes,
        CreatedAt:  testimony.CreatedAt,
        UpdatedAt:  testimony.UpdatedAt,
    }

    return response, nil
}

func (s *testimonyService) GetTestimonyByTransactionID(transactionID uint) (*TestimonyResponse, error) {
    testimony, err := s.testimonyRepo.FindByTransactionID(transactionID)
    if err != nil {
        return nil, err
    }
    if testimony == nil {
        return nil, errors.New("testimony not found for this transaction")
    }

    // Get user
    user, err := s.userRepo.FindByID(testimony.UserID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    response := &TestimonyResponse{
        ID:            testimony.ID,
        TransactionID: testimony.TransactionID,
        UserID:        testimony.UserID,
        User: UserResponse{
            ID:            user.ID,
            Username:      user.Username,
            Email:         user.Email,
            FullName:      user.FullName,
            Phone:         user.Phone,
            AvatarURL:     user.AvatarURL,
            Role:          user.Role,
            IsVerified:    user.IsVerified,
            IsActive:      user.IsActive,
            BankName:      user.BankName,
            AccountNumber: user.AccountNumber,
            AccountHolder: user.AccountHolder,
            CreatedAt:     user.CreatedAt,
            UpdatedAt:     user.UpdatedAt,
        },
        Rating:     testimony.Rating,
        Comment:    testimony.Comment,
        AdminNotes: testimony.AdminNotes,
        CreatedAt:  testimony.CreatedAt,
        UpdatedAt:  testimony.UpdatedAt,
    }

    return response, nil
}

func (s *testimonyService) UpdateTestimony(testimonyID uint, req UpdateTestimonyRequest, userID uint) (*TestimonyResponse, error) {
    testimony, err := s.testimonyRepo.FindByID(testimonyID)
    if err != nil {
        return nil, err
    }
    if testimony == nil {
        return nil, errors.New("testimony not found")
    }

    // Only the user who created the testimony can update it
    if testimony.UserID != userID {
        return nil, errors.New("you can only update your own testimony")
    }

    // Validate rating if provided
    if req.Rating > 0 && (req.Rating < 1 || req.Rating > 5) {
        return nil, errors.New("rating must be between 1 and 5")
    }

    // Update fields if provided
    if req.Rating > 0 {
        testimony.Rating = req.Rating
    }
    if req.Comment != "" {
        testimony.Comment = req.Comment
    }

    testimony.UpdatedAt = time.Now()

    if err := s.testimonyRepo.Update(testimony); err != nil {
        return nil, err
    }

    // Get user
    user, err := s.userRepo.FindByID(userID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    response := &TestimonyResponse{
        ID:            testimony.ID,
        TransactionID: testimony.TransactionID,
        UserID:        testimony.UserID,
        User: UserResponse{
            ID:            user.ID,
            Username:      user.Username,
            Email:         user.Email,
            FullName:      user.FullName,
            Phone:         user.Phone,
            AvatarURL:     user.AvatarURL,
            Role:          user.Role,
            IsVerified:    user.IsVerified,
            IsActive:      user.IsActive,
            BankName:      user.BankName,
            AccountNumber: user.AccountNumber,
            AccountHolder: user.AccountHolder,
            CreatedAt:     user.CreatedAt,
            UpdatedAt:     user.UpdatedAt,
        },
        Rating:     testimony.Rating,
        Comment:    testimony.Comment,
        AdminNotes: testimony.AdminNotes,
        CreatedAt:  testimony.CreatedAt,
        UpdatedAt:  testimony.UpdatedAt,
    }

    return response, nil
}

func (s *testimonyService) DeleteTestimony(testimonyID uint, userID uint) error {
    testimony, err := s.testimonyRepo.FindByID(testimonyID)
    if err != nil {
        return err
    }
    if testimony == nil {
        return errors.New("testimony not found")
    }

    // Only the user who created the testimony can delete it
    if testimony.UserID != userID {
        return errors.New("you can only delete your own testimony")
    }

    return s.testimonyRepo.Delete(testimonyID)
}

func (s *testimonyService) ListTestimonies(page, limit int) ([]TestimonyResponse, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    testimonies, err := s.testimonyRepo.ListAll(offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.testimonyRepo.CountAll()
    if err != nil {
        return nil, 0, err
    }

    responses := make([]TestimonyResponse, len(testimonies))
    for i, testimony := range testimonies {
        // Get user for each testimony
        user, err := s.userRepo.FindByID(testimony.UserID)
        if err != nil || user == nil {
            continue // Skip if user not found
        }

        responses[i] = TestimonyResponse{
            ID:            testimony.ID,
            TransactionID: testimony.TransactionID,
            UserID:        testimony.UserID,
            User: UserResponse{
                ID:            user.ID,
                Username:      user.Username,
                Email:         user.Email,
                FullName:      user.FullName,
                Phone:         user.Phone,
                AvatarURL:     user.AvatarURL,
                Role:          user.Role,
                IsVerified:    user.IsVerified,
                IsActive:      user.IsActive,
                BankName:      user.BankName,
                AccountNumber: user.AccountNumber,
                AccountHolder: user.AccountHolder,
                CreatedAt:     user.CreatedAt,
                UpdatedAt:     user.UpdatedAt,
            },
            Rating:     testimony.Rating,
            Comment:    testimony.Comment,
            AdminNotes: testimony.AdminNotes,
            CreatedAt:  testimony.CreatedAt,
            UpdatedAt:  testimony.UpdatedAt,
        }
    }

    return responses, total, nil
}

func (s *testimonyService) ListUserTestimonies(userID uint, page, limit int) ([]TestimonyResponse, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    testimonies, err := s.testimonyRepo.ListByUserID(userID, offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.testimonyRepo.CountByUserID(userID)
    if err != nil {
        return nil, 0, err
    }

    responses := make([]TestimonyResponse, len(testimonies))
    for i, testimony := range testimonies {
        // Get user
        user, err := s.userRepo.FindByID(testimony.UserID)
        if err != nil || user == nil {
            continue
        }

        responses[i] = TestimonyResponse{
            ID:            testimony.ID,
            TransactionID: testimony.TransactionID,
            UserID:        testimony.UserID,
            User: UserResponse{
                ID:            user.ID,
                Username:      user.Username,
                Email:         user.Email,
                FullName:      user.FullName,
                Phone:         user.Phone,
                AvatarURL:     user.AvatarURL,
                Role:          user.Role,
                IsVerified:    user.IsVerified,
                IsActive:      user.IsActive,
                BankName:      user.BankName,
                AccountNumber: user.AccountNumber,
                AccountHolder: user.AccountHolder,
                CreatedAt:     user.CreatedAt,
                UpdatedAt:     user.UpdatedAt,
            },
            Rating:     testimony.Rating,
            Comment:    testimony.Comment,
            AdminNotes: testimony.AdminNotes,
            CreatedAt:  testimony.CreatedAt,
            UpdatedAt:  testimony.UpdatedAt,
        }
    }

    return responses, total, nil
}