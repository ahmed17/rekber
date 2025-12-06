package models

import (
    "time"
)

type Testimony struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    TransactionID uint           `gorm:"uniqueIndex;not null" json:"transaction_id"`
    Transaction   Transaction    `gorm:"foreignKey:TransactionID" json:"transaction"`
    
    UserID        uint           `gorm:"not null;index" json:"user_id"`
    User          User           `gorm:"foreignKey:UserID" json:"user"`
    
    Rating        int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
    Comment       string         `json:"comment"`
    AdminNotes    string         `json:"admin_notes,omitempty"`
    
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
}