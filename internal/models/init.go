package models

import (
    "log"
    "gorm.io/gorm"
    "rekber/internal/config"
)

func AutoMigrate() error {
    db := config.DB
    
    err := db.AutoMigrate(
        &User{},
        &Transaction{},
        &Testimony{},
        &AuditLog{},
    )
    
    if err != nil {
        log.Printf("Auto migration failed: %v", err)
        return err
    }
    
    log.Println("Auto migration completed successfully")
    return nil
}