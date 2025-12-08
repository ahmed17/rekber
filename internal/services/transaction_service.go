package services

import (
    "errors"
    "time"

    "rekber/internal/models"
    "rekber/internal/repositories"
)

// Tipe lokal untuk transaction service
type CreateTransactionRequest struct {
    Title       string  `json:"title" binding:"required,min=5,max=200"`
    Description string  `json:"description" binding:"omitempty,max=1000"`
    Amount      float64 `json:"amount" binding:"required,gt=0"`
    Fee         float64 `json:"fee" binding:"gte=0"`
    SellerID    uint    `json:"seller_id" binding:"required"`
}

type UpdateTransactionRequest struct {
    Title       string  `json:"title" binding:"omitempty,min=5,max=200"`
    Description string  `json:"description" binding:"omitempty,max=1000"`
}

type UploadPaymentProofRequest struct {
    PaymentProofURL string `json:"payment_proof_url" binding:"required,url"`
}

type UpdateShippingDataRequest struct {
    ShippingData string `json:"shipping_data" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=validated processing shipped completed cancelled"`
}

type TransactionResponse struct {
    ID              uint      `json:"id"`
    TransactionCode string    `json:"transaction_code"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    Amount          float64   `json:"amount"`
    Fee             float64   `json:"fee"`
    TotalAmount     float64   `json:"total_amount"`
    Status          string    `json:"status"`
    
    BuyerID  uint         `json:"buyer_id"`
    Buyer    UserResponse `json:"buyer"`
    SellerID uint         `json:"seller_id"`
    Seller   UserResponse `json:"seller"`
    
    PaymentProofURL  string     `json:"payment_proof_url,omitempty"`
    PaymentVerifiedAt *time.Time `json:"payment_verified_at,omitempty"`
    
    ShippingData   string     `json:"shipping_data,omitempty"`
    ShippingSentAt *time.Time `json:"shipping_sent_at,omitempty"`
    
    PaidAt       *time.Time `json:"paid_at,omitempty"`
    ValidatedAt  *time.Time `json:"validated_at,omitempty"`
    ProcessingAt *time.Time `json:"processing_at,omitempty"`
    ShippedAt    *time.Time `json:"shipped_at,omitempty"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
    CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
    
    AutoCompleteAt *time.Time `json:"auto_complete_at,omitempty"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}

// TransactionService interface
type TransactionService interface {
    CreateTransaction(req CreateTransactionRequest, buyerID uint) (*TransactionResponse, error)
    GetTransactionByID(transactionID uint, userID uint, role string) (*TransactionResponse, error)
    GetTransactionByCode(code string, userID uint, role string) (*TransactionResponse, error)
    ListUserTransactions(userID uint, role string, page, limit int) ([]TransactionResponse, int64, error)
    ListAllTransactions(page, limit int) ([]TransactionResponse, int64, error)
    UpdateTransaction(transactionID uint, req UpdateTransactionRequest, userID uint) (*TransactionResponse, error)
    UploadPaymentProof(transactionID uint, req UploadPaymentProofRequest, userID uint) (*TransactionResponse, error)
    UpdateShippingData(transactionID uint, req UpdateShippingDataRequest, sellerID uint) (*TransactionResponse, error)
    UpdateStatus(transactionID uint, req UpdateStatusRequest, adminID uint) (*TransactionResponse, error)
    CancelTransaction(transactionID uint, userID uint) (*TransactionResponse, error)
    AutoCompleteTransactions() error
    ValidatePayment(transactionID uint, adminID uint) (*TransactionResponse, error)
    CompleteTransaction(transactionID uint, userID uint) (*TransactionResponse, error)
}

type transactionService struct {
    transactionRepo repositories.TransactionRepository
    userRepo        repositories.UserRepository
}

func NewTransactionService(transactionRepo repositories.TransactionRepository, userRepo repositories.UserRepository) TransactionService {
    return &transactionService{
        transactionRepo: transactionRepo,
        userRepo:        userRepo,
    }
}

func (s *transactionService) CreateTransaction(req CreateTransactionRequest, buyerID uint) (*TransactionResponse, error) {
    // Validate seller exists and is a seller
    seller, err := s.userRepo.FindByID(req.SellerID)
    if err != nil {
        return nil, err
    }
    if seller == nil {
        return nil, errors.New("seller not found")
    }
    if seller.Role != "seller" {
        return nil, errors.New("user is not a seller")
    }

    // Check if buyer and seller are the same
    if buyerID == req.SellerID {
        return nil, errors.New("buyer and seller cannot be the same")
    }

    // Validate amount
    if req.Amount <= 0 {
        return nil, errors.New("amount must be greater than 0")
    }
    if req.Fee < 0 {
        return nil, errors.New("fee cannot be negative")
    }

    // Create transaction
    transaction := &models.Transaction{
        TransactionCode: generateTransactionCode(),
        Title:           req.Title,
        Description:     req.Description,
        Amount:          req.Amount,
        Fee:             req.Fee,
        TotalAmount:     req.Amount + req.Fee,
        BuyerID:         buyerID,
        SellerID:        req.SellerID,
        Status:          models.StatusWaitingPayment,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }

    if err := s.transactionRepo.Create(transaction); err != nil {
        return nil, err
    }

    // Reload to get relations
    newTransaction, err := s.transactionRepo.FindByID(transaction.ID)
    if err != nil {
        return nil, err
    }

    // Convert to response
    response := s.transactionToResponse(newTransaction)
    return response, nil
}

func (s *transactionService) GetTransactionByID(transactionID uint, userID uint, role string) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Check authorization: user must be buyer, seller, or admin
    if transaction.BuyerID != userID && transaction.SellerID != userID && role != "admin" {
        return nil, errors.New("unauthorized to view this transaction")
    }

    response := s.transactionToResponse(transaction)
    return response, nil
}

func (s *transactionService) GetTransactionByCode(code string, userID uint, role string) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByCode(code)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Check authorization
    if transaction.BuyerID != userID && transaction.SellerID != userID && role != "admin" {
        return nil, errors.New("unauthorized to view this transaction")
    }

    response := s.transactionToResponse(transaction)
    return response, nil
}

func (s *transactionService) ListUserTransactions(userID uint, role string, page, limit int) ([]TransactionResponse, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    transactions, err := s.transactionRepo.ListByUserID(userID, role, offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.transactionRepo.CountByUserID(userID, role)
    if err != nil {
        return nil, 0, err
    }

    responses := make([]TransactionResponse, len(transactions))
    for i, transaction := range transactions {
        responses[i] = *s.transactionToResponse(&transaction)
    }

    return responses, total, nil
}

func (s *transactionService) ListAllTransactions(page, limit int) ([]TransactionResponse, int64, error) {
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    offset := (page - 1) * limit

    transactions, err := s.transactionRepo.ListAll(offset, limit)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.transactionRepo.CountAll()
    if err != nil {
        return nil, 0, err
    }

    responses := make([]TransactionResponse, len(transactions))
    for i, transaction := range transactions {
        responses[i] = *s.transactionToResponse(&transaction)
    }

    return responses, total, nil
}

func (s *transactionService) UpdateTransaction(transactionID uint, req UpdateTransactionRequest, userID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only buyer can update transaction details (only when status is waiting_payment)
    if transaction.BuyerID != userID {
        return nil, errors.New("only buyer can update transaction details")
    }

    if transaction.Status != models.StatusWaitingPayment {
        return nil, errors.New("can only update transaction in waiting_payment status")
    }

    // Update fields
    if req.Title != "" {
        transaction.Title = req.Title
    }
    if req.Description != "" {
        transaction.Description = req.Description
    }

    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) UploadPaymentProof(transactionID uint, req UploadPaymentProofRequest, userID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only buyer can upload payment proof
    if transaction.BuyerID != userID {
        return nil, errors.New("only buyer can upload payment proof")
    }

    if transaction.Status != models.StatusWaitingPayment {
        return nil, errors.New("can only upload payment proof in waiting_payment status")
    }

    // Update payment proof and status
    transaction.PaymentProofURL = req.PaymentProofURL
    transaction.Status = models.StatusPending
    transaction.PaidAt = timePtr(time.Now())
    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) ValidatePayment(transactionID uint, adminID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    if transaction.Status != models.StatusPending {
        return nil, errors.New("can only validate payment in pending status")
    }

    // Validate payment (admin action)
    transaction.Status = models.StatusValidated
    transaction.ValidatedAt = timePtr(time.Now())
    transaction.PaymentVerifiedAt = timePtr(time.Now())
    transaction.PaymentVerifiedBy = &adminID
    transaction.AdminID = &adminID
    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) UpdateStatus(transactionID uint, req UpdateStatusRequest, adminID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only admin can update status
    // (we assume the adminID is passed and validated at handler level)

    // Validate status transition
    currentStatus := transaction.Status
    newStatus := models.TransactionStatus(req.Status)

    // Valid status transitions (admin can change to certain statuses)
    validTransitions := map[models.TransactionStatus][]models.TransactionStatus{
        models.StatusPending:    {models.StatusValidated, models.StatusCancelled},
        models.StatusValidated:  {models.StatusProcessing, models.StatusCancelled},
        models.StatusProcessing: {models.StatusShipped, models.StatusCancelled},
        models.StatusShipped:    {models.StatusCompleted, models.StatusCancelled},
    }

    // Check if transition is valid
    if allowedTransitions, ok := validTransitions[currentStatus]; ok {
        valid := false
        for _, allowed := range allowedTransitions {
            if newStatus == allowed {
                valid = true
                break
            }
        }
        if !valid {
            return nil, errors.New("invalid status transition")
        }
    } else {
        // Cannot change from other statuses
        return nil, errors.New("cannot change status from current status")
    }

    // Update status and timestamps
    transaction.Status = newStatus
    transaction.AdminID = &adminID
    transaction.UpdatedAt = time.Now()

    // Set timestamp based on new status
    now := time.Now()
    switch newStatus {
    case models.StatusValidated:
        transaction.ValidatedAt = &now
    case models.StatusProcessing:
        transaction.ProcessingAt = &now
    case models.StatusShipped:
        transaction.ShippedAt = &now
        // Set auto complete deadline (2x24 jam dari sekarang)
        autoCompleteAt := now.Add(48 * time.Hour)
        transaction.AutoCompleteAt = &autoCompleteAt
    case models.StatusCompleted:
        transaction.CompletedAt = &now
    case models.StatusCancelled:
        transaction.CancelledAt = &now
    }

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) UpdateShippingData(transactionID uint, req UpdateShippingDataRequest, sellerID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only seller can update shipping data
    if transaction.SellerID != sellerID {
        return nil, errors.New("only seller can update shipping data")
    }

    // Only validated transactions can have shipping data updated
    if transaction.Status != models.StatusValidated {
        return nil, errors.New("can only update shipping data in validated status")
    }

    // Update shipping data and status
    transaction.ShippingData = req.ShippingData
    transaction.Status = models.StatusShipped
    transaction.ShippingSentAt = timePtr(time.Now())
    
    // Set auto complete deadline (2x24 jam dari sekarang)
    autoCompleteAt := time.Now().Add(48 * time.Hour)
    transaction.AutoCompleteAt = &autoCompleteAt
    
    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) CompleteTransaction(transactionID uint, userID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only buyer can complete transaction
    if transaction.BuyerID != userID {
        return nil, errors.New("only buyer can complete transaction")
    }

    if transaction.Status != models.StatusShipped {
        return nil, errors.New("can only complete shipped transactions")
    }

    // Complete transaction
    transaction.Status = models.StatusCompleted
    transaction.CompletedAt = timePtr(time.Now())
    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) CancelTransaction(transactionID uint, userID uint) (*TransactionResponse, error) {
    transaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }
    if transaction == nil {
        return nil, errors.New("transaction not found")
    }

    // Only buyer can cancel, and only in waiting_payment status
    if transaction.BuyerID != userID {
        return nil, errors.New("only buyer can cancel transaction")
    }

    if transaction.Status != models.StatusWaitingPayment {
        return nil, errors.New("can only cancel transaction in waiting_payment status")
    }

    // Update status
    transaction.Status = models.StatusCancelled
    transaction.CancelledAt = timePtr(time.Now())
    transaction.UpdatedAt = time.Now()

    if err := s.transactionRepo.Update(transaction); err != nil {
        return nil, err
    }

    // Reload
    updatedTransaction, err := s.transactionRepo.FindByID(transactionID)
    if err != nil {
        return nil, err
    }

    response := s.transactionToResponse(updatedTransaction)
    return response, nil
}

func (s *transactionService) AutoCompleteTransactions() error {
    // Find shipped transactions with expired auto_complete_at
    transactions, err := s.transactionRepo.FindShippedWithExpiredAutoComplete()
    if err != nil {
        return err
    }

    for _, transaction := range transactions {
        // Update status to completed
        transaction.Status = models.StatusCompleted
        now := time.Now()
        transaction.CompletedAt = &now
        transaction.UpdatedAt = now

        if err := s.transactionRepo.Update(&transaction); err != nil {
            // Log error but continue with other transactions
            continue
        }
    }

    return nil
}

// Helper function untuk convert models.Transaction ke TransactionResponse
func (s *transactionService) transactionToResponse(t *models.Transaction) *TransactionResponse {
    if t == nil {
        return nil
    }

    // Get user responses
    buyerResponse := UserResponse{
        ID:            t.Buyer.ID,
        Username:      t.Buyer.Username,
        Email:         t.Buyer.Email,
        FullName:      t.Buyer.FullName,
        Phone:         t.Buyer.Phone,
        AvatarURL:     t.Buyer.AvatarURL,
        Role:          t.Buyer.Role,
        IsVerified:    t.Buyer.IsVerified,
        IsActive:      t.Buyer.IsActive,
        BankName:      t.Buyer.BankName,
        AccountNumber: t.Buyer.AccountNumber,
        AccountHolder: t.Buyer.AccountHolder,
        CreatedAt:     t.Buyer.CreatedAt,
        UpdatedAt:     t.Buyer.UpdatedAt,
    }

    sellerResponse := UserResponse{
        ID:            t.Seller.ID,
        Username:      t.Seller.Username,
        Email:         t.Seller.Email,
        FullName:      t.Seller.FullName,
        Phone:         t.Seller.Phone,
        AvatarURL:     t.Seller.AvatarURL,
        Role:          t.Seller.Role,
        IsVerified:    t.Seller.IsVerified,
        IsActive:      t.Seller.IsActive,
        BankName:      t.Seller.BankName,
        AccountNumber: t.Seller.AccountNumber,
        AccountHolder: t.Seller.AccountHolder,
        CreatedAt:     t.Seller.CreatedAt,
        UpdatedAt:     t.Seller.UpdatedAt,
    }

    return &TransactionResponse{
        ID:              t.ID,
        TransactionCode: t.TransactionCode,
        Title:           t.Title,
        Description:     t.Description,
        Amount:          t.Amount,
        Fee:             t.Fee,
        TotalAmount:     t.TotalAmount,
        Status:          string(t.Status),
        BuyerID:         t.BuyerID,
        Buyer:           buyerResponse,
        SellerID:        t.SellerID,
        Seller:          sellerResponse,
        PaymentProofURL: t.PaymentProofURL,
        PaymentVerifiedAt: t.PaymentVerifiedAt,
        ShippingData:    t.ShippingData,
        ShippingSentAt:  t.ShippingSentAt,
        PaidAt:          t.PaidAt,
        ValidatedAt:     t.ValidatedAt,
        ProcessingAt:    t.ProcessingAt,
        ShippedAt:       t.ShippedAt,
        CompletedAt:     t.CompletedAt,
        CancelledAt:     t.CancelledAt,
        AutoCompleteAt:  t.AutoCompleteAt,
        CreatedAt:       t.CreatedAt,
        UpdatedAt:       t.UpdatedAt,
    }
}

func generateTransactionCode() string {
    timestamp := time.Now().UnixNano() / 1000000
    return "REKBER-" + string(timestamp)
}

func timePtr(t time.Time) *time.Time {
    return &t
}