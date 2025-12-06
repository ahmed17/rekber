package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *sql.DB, dbName string) error {
    // Create driver instance
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create migration driver: %w", err)
    }

    // Get absolute path to migrations folder
    workDir, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("could not get working directory: %w", err)
    }
    
    migrationsPath := filepath.Join(workDir, "migrations")
    
    // Check if migrations folder exists
    if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
        return fmt.Errorf("migrations folder does not exist: %s", migrationsPath)
    }

    // Create migrate instance with sslmode=disable
    m, err := migrate.NewWithDatabaseInstance(
        "file://"+migrationsPath,
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migration instance: %w", err)
    }

    // Run migrations
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("could not run migrations: %w", err)
    }

    log.Println("Migrations completed successfully")
    return nil
}