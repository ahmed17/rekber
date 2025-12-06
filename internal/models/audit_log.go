package models

import (
    "time"
    "gorm.io/gorm"
)

type AuditLogAction string

const (
    ActionCreate    AuditLogAction = "create"
    ActionUpdate    AuditLogAction = "update"
    ActionDelete    AuditLogAction = "delete"
    ActionLogin     AuditLogAction = "login"
    ActionLogout    AuditLogAction = "logout"
    ActionStatusChange AuditLogAction = "status_change"
)

type AuditLog struct {
    ID         uint           `gorm:"primaryKey" json:"id"`
    UserID     *uint          `gorm:"index" json:"user_id,omitempty"`
    User       *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
    
    Action     AuditLogAction `gorm:"size:50;not null" json:"action"`
    EntityType string         `gorm:"size:50;not null" json:"entity_type"`
    EntityID   uint           `json:"entity_id"`
    
    OldData    string         `gorm:"type:text" json:"old_data,omitempty"`
    NewData    string         `gorm:"type:text" json:"new_data,omitempty"`
    Changes    string         `gorm:"type:text" json:"changes,omitempty"`
    
    IPAddress  string         `gorm:"size:45" json:"ip_address,omitempty"`
    UserAgent  string         `json:"user_agent,omitempty"`
    
    CreatedAt  time.Time      `json:"created_at"`
}