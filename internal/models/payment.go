package models

import (
    "time"

    "gorm.io/gorm"
)

type PaymentMethod string

const (
    MethodBankTransfer PaymentMethod = "bank_transfer"
    MethodEWallet      PaymentMethod = "e_wallet"
)

type PaymentStatus string

const (
    PaymentPending  PaymentStatus = "pending"
    PaymentVerified PaymentStatus = "verified"
    PaymentRejected PaymentStatus = "rejected"
)

type Payment struct {
    ID            uint          `gorm:"primarykey" json:"id"`
    CreatedAt     time.Time     `json:"created_at"`
    UpdatedAt     time.Time     `json:"updated_at"`
    
    TransactionID uint          `gorm:"not null" json:"transaction_id"`
    PaymentMethod PaymentMethod `gorm:"type:varchar(30);default:'bank_transfer'" json:"payment_method"`
    
    // Informasi transfer
    BankName      string        `json:"bank_name"`
    AccountNumber string        `json:"account_number"`
    AccountHolder string        `json:"account_holder"`
    Amount        float64       `gorm:"not null" json:"amount"`
    
    // Bukti transfer
    ProofImage    string        `json:"proof_image"` // Path ke gambar
    Notes         string        `json:"notes"`
    
    Status        PaymentStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
    VerifiedBy    *uint         `json:"verified_by"` // Admin yang memverifikasi
    VerifiedAt    *time.Time    `json:"verified_at"`
    
    // Relations
    Transaction   Transaction   `gorm:"foreignKey:TransactionID" json:"transaction"`
    Verifier      *User         `gorm:"foreignKey:VerifiedBy" json:"verifier"`
}