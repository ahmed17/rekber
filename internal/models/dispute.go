package models

import (
    "time"

    "gorm.io/gorm"
)

type DisputeStatus string

const (
    DisputeOpen     DisputeStatus = "open"
    DisputeResolved DisputeStatus = "resolved"
    DisputeClosed   DisputeStatus = "closed"
)

type Dispute struct {
    ID            uint         `gorm:"primarykey" json:"id"`
    CreatedAt     time.Time    `json:"created_at"`
    UpdatedAt     time.Time    `json:"updated_at"`
    
    TransactionID uint         `gorm:"not null" json:"transaction_id"`
    RaisedByID    uint         `gorm:"not null" json:"raised_by_id"` // User yang membuka dispute
    Reason        string       `gorm:"not null" json:"reason"`
    Description   string       `json:"description"`
    Status        DisputeStatus `gorm:"type:varchar(20);default:'open'" json:"status"`
    
    // Relations
    Transaction   Transaction  `gorm:"foreignKey:TransactionID" json:"transaction"`
    RaisedBy      User         `gorm:"foreignKey:RaisedByID" json:"raised_by"`
}