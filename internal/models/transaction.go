package models

import (
    "time"
    "gorm.io/gorm"
)

type TransactionStatus string

const (
    StatusWaitingPayment TransactionStatus = "waiting_payment"
    StatusPending        TransactionStatus = "pending"
    StatusValidated      TransactionStatus = "validated"
    StatusProcessing     TransactionStatus = "processing"
    StatusShipped        TransactionStatus = "shipped"
    StatusCompleted      TransactionStatus = "completed"
    StatusCancelled      TransactionStatus = "cancelled"
    StatusDisputed       TransactionStatus = "disputed"
)

type Transaction struct {
    ID              uint              `gorm:"primaryKey" json:"id"`
    TransactionCode string            `gorm:"size:50;uniqueIndex;not null" json:"transaction_code"`
    Title           string            `gorm:"size:200;not null" json:"title"`
    Description     string            `json:"description"`
    
    Amount      float64 `gorm:"type:decimal(12,2);not null" json:"amount"`
    Fee         float64 `gorm:"type:decimal(12,2);default:0" json:"fee"`
    TotalAmount float64 `gorm:"type:decimal(12,2);not null" json:"total_amount"`
    
    Status TransactionStatus `gorm:"size:30;default:'waiting_payment';index" json:"status"`
    
    // Foreign keys
    BuyerID  uint `gorm:"not null;index" json:"buyer_id"`
    Buyer    User `gorm:"foreignKey:BuyerID" json:"buyer"`
    
    SellerID uint `gorm:"not null;index" json:"seller_id"`
    Seller   User `gorm:"foreignKey:SellerID" json:"seller"`
    
    AdminID *uint `gorm:"index" json:"admin_id,omitempty"`
    Admin   *User `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
    
    PaymentProofURL  string     `gorm:"size:255" json:"payment_proof_url,omitempty"`
    PaymentVerifiedBy *uint     `gorm:"index" json:"payment_verified_by,omitempty"`
    PaymentVerifiedAt *time.Time `json:"payment_verified_at,omitempty"`
    
    ShippingData  string     `json:"shipping_data,omitempty"`
    ShippingSentAt *time.Time `json:"shipping_sent_at,omitempty"`
    
    PaidAt       *time.Time `json:"paid_at,omitempty"`
    ValidatedAt  *time.Time `json:"validated_at,omitempty"`
    ProcessingAt *time.Time `json:"processing_at,omitempty"`
    ShippedAt    *time.Time `json:"shipped_at,omitempty"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
    CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
    
    AutoCompleteAt *time.Time `gorm:"index" json:"auto_complete_at,omitempty"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook untuk generate transaction code
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
    if t.TransactionCode == "" {
        t.TransactionCode = GenerateTransactionCode()
    }
    t.TotalAmount = t.Amount + t.Fee
    return nil
}

func GenerateTransactionCode() string {
    // Implementasi generate code, misal: REKBER-<timestamp>-<random>
    return "REKBER-" + time.Now().Format("20060102150405") // contoh sederhana
}