package models

import (
    "log"
    "rekber/internal/config"
)

func AutoMigrate() error {
    db := config.DB
    
    err := db.AutoMigrate(
        &User{},
        &Transaction{},
        &Testimony{},
    )
    
    if err != nil {
        log.Printf("Auto migration failed: %v", err)
        return err
    }
    
    log.Println("Auto migration completed successfully")
    return nil
}