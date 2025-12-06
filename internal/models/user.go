package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
    Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
    PasswordHash string         `gorm:"size:255;not null" json:"-"`
    FullName     string         `gorm:"size:100;not null" json:"full_name"`
    Phone        string         `gorm:"size:20" json:"phone"`
    AvatarURL    string         `gorm:"size:255" json:"avatar_url"`
    
    Role       string `gorm:"size:20;default:'buyer';check:role IN ('buyer','seller','admin')" json:"role"`
    IsVerified bool   `gorm:"default:false" json:"is_verified"`
    IsActive   bool   `gorm:"default:true" json:"is_active"`
    
    BankName       string `gorm:"size:50" json:"bank_name,omitempty"`
    AccountNumber  string `gorm:"size:50" json:"account_number,omitempty"`
    AccountHolder  string `gorm:"size:100" json:"account_holder,omitempty"`
    
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserResponse untuk menghindari menampilkan sensitive data
type UserResponse struct {
    ID            uint      `json:"id"`
    Username      string    `json:"username"`
    Email         string    `json:"email"`
    FullName      string    `json:"full_name"`
    Phone         string    `json:"phone"`
    AvatarURL     string    `json:"avatar_url"`
    Role          string    `json:"role"`
    IsVerified    bool      `json:"is_verified"`
    IsActive      bool      `json:"is_active"`
    BankName      string    `json:"bank_name,omitempty"`
    AccountNumber string    `json:"account_number,omitempty"`
    AccountHolder string    `json:"account_holder,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:            u.ID,
        Username:      u.Username,
        Email:         u.Email,
        FullName:      u.FullName,
        Phone:         u.Phone,
        AvatarURL:     u.AvatarURL,
        Role:          u.Role,
        IsVerified:    u.IsVerified,
        IsActive:      u.IsActive,
        BankName:      u.BankName,
        AccountNumber: u.AccountNumber,
        AccountHolder: u.AccountHolder,
        CreatedAt:     u.CreatedAt,
        UpdatedAt:     u.UpdatedAt,
    }
}