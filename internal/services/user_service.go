package services

import (
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"
    "rekber/internal/models"
    "rekber/internal/repositories"
)

// Tipe lokal untuk user service
type CreateUserRequest struct {
    Username    string `json:"username" binding:"required,min=3,max=50"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    FullName    string `json:"full_name" binding:"required,min=2,max=100"`
    Phone       string `json:"phone" binding:"omitempty"`
    Role        string `json:"role" binding:"omitempty,oneof=buyer seller"`
    
    // Optional for seller
    BankName      string `json:"bank_name" binding:"omitempty"`
    AccountNumber string `json:"account_number" binding:"omitempty"`
    AccountHolder string `json:"account_holder" binding:"omitempty"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
    FullName     string `json:"full_name" binding:"omitempty,min=2,max=100"`
    Phone        string `json:"phone" binding:"omitempty"`
    AvatarURL    string `json:"avatar_url" binding:"omitempty,url"`
    BankName     string `json:"bank_name" binding:"omitempty"`
    AccountNumber string `json:"account_number" binding:"omitempty"`
    AccountHolder string `json:"account_holder" binding:"omitempty"`
}

type UserResponse struct {
    ID           uint      `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    FullName     string    `json:"full_name"`
    Phone        string    `json:"phone"`
    AvatarURL    string    `json:"avatar_url"`
    Role         string    `json:"role"`
    IsVerified   bool      `json:"is_verified"`
    IsActive     bool      `json:"is_active"`
    BankName     string    `json:"bank_name,omitempty"`
    AccountNumber string   `json:"account_number,omitempty"`
    AccountHolder string   `json:"account_holder,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// UserService interface
type UserService interface {
    Register(req CreateUserRequest) (*UserResponse, error)
    Login(req LoginRequest) (string, *UserResponse, error)
    GetProfile(userID uint) (*UserResponse, error)
    UpdateProfile(userID uint, req UpdateUserRequest) (*UserResponse, error)
    DeleteAccount(userID uint) error
    ListUsers(page, limit int) ([]UserResponse, int64, error)
    GetUserByID(userID uint) (*UserResponse, error)
    ChangePassword(userID uint, oldPassword, newPassword string) error
}

type userService struct {
    userRepo    repositories.UserRepository
    authService AuthService
}

func NewUserService(userRepo repositories.UserRepository, authService AuthService) UserService {
    return &userService{
        userRepo:    userRepo,
        authService: authService,
    }
}

func (s *userService) Register(req CreateUserRequest) (*UserResponse, error) {
    // Check if user already exists
    existingUser, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    existingUserByUsername, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        return nil, err
    }
    if existingUserByUsername != nil {
        return nil, errors.New("username already taken")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
        Phone:        req.Phone,
        Role:         req.Role,
        IsVerified:   false,
        IsActive:     true,
        BankName:      req.BankName,
        AccountNumber: req.AccountNumber,
        AccountHolder: req.AccountHolder,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if user.Role == "" {
        user.Role = "buyer"
    }

    // Validate seller data
    if user.Role == "seller" && (user.BankName == "" || user.AccountNumber == "" || user.AccountHolder == "") {
        return nil, errors.New("sellers must provide bank account information")
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Convert to response
    response := &UserResponse{
        ID:            user.ID,
        Username:      user.Username,
        Email:         user.Email,
        FullName:      user.FullName,
        Phone:         user.Phone,
        Role:          user.Role,
        IsVerified:    user.IsVerified,
        IsActive:      user.IsActive,
        BankName:      user.BankName,
        AccountNumber: user.AccountNumber,
        AccountHolder: user.AccountHolder,
        CreatedAt:     user.CreatedAt,
        UpdatedAt:     user.UpdatedAt,
    }

    return response, nil
}

func (s *userService) Login(req LoginRequest) (string, *UserResponse, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        return "", nil, err
    }
    if user == nil {
        return "", nil, errors.New("invalid credentials")
    }

    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    // Check if user is active
    if !user.IsActive {
        return "", nil, errors.New("account is not active")
    }

    // Generate JWT token
    token, err := s.authService.GenerateToken(user.ID, user.Role)
    if err != nil {
        return "", nil, err
    }

    // Convert to response
    response := &UserResponse{
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
    }

    return token, response, nil
}

func (s *userService) GetProfile(userID uint) (*UserResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    response := &UserResponse{
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
    }

    return response, nil
}

func (s *userService) UpdateProfile(userID uint, req UpdateUserRequest) (*UserResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    // Update fields if provided
    if req.FullName != "" {
        user.FullName = req.FullName
    }
    if req.Phone != "" {
        user.Phone = req.Phone
    }
    if req.AvatarURL != "" {
        user.AvatarURL = req.AvatarURL
    }
    if req.BankName != "" {
        user.BankName = req.BankName
    }
    if req.AccountNumber != "" {
        user.AccountNumber = req.AccountNumber
    }
    if req.AccountHolder != "" {
        user.AccountHolder = req.AccountHolder
    }

    user.UpdatedAt = time.Now()

    if err := s.userRepo.Update(user); err != nil {
        return nil, err
    }

    response := &UserResponse{
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
    }

    return response, nil
}

func (s *userService) DeleteAccount(userID uint) error {
    return s.userRepo.Delete(userID)
}

func (s *userService) ListUsers(page, limit int) ([]UserResponse, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    users, err := s.userRepo.List(offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.userRepo.Count()
    if err != nil {
        return nil, 0, err
    }

    responses := make([]UserResponse, len(users))
    for i, user := range users {
        responses[i] = UserResponse{
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
        }
    }

    return responses, total, nil
}

func (s *userService) GetUserByID(userID uint) (*UserResponse, error) {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("user not found")
    }

    response := &UserResponse{
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
    }

    return response, nil
}

func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    if user == nil {
        return errors.New("user not found")
    }

    // Verify old password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
        return errors.New("incorrect old password")
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.PasswordHash = string(hashedPassword)
    user.UpdatedAt = time.Now()

    return s.userRepo.Update(user)
}