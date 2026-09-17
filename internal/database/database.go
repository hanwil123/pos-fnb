package database

import (
	"fmt"
	"log"

	"github.com/yourorg/pos-fnb-backend/internal/config"
	"github.com/yourorg/pos-fnb-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect membuka koneksi ke PostgreSQL menggunakan GORM.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal konek ke database: %w", err)
	}

	log.Println("database terkoneksi")
	return db, nil
}

// AutoMigrate menjalankan migrasi otomatis untuk semua model.
// Untuk production, disarankan pindah ke tool migration eksplisit (golang-migrate)
// supaya perubahan schema bisa di-review & di-rollback dengan aman.
func AutoMigrate(db *gorm.DB) error {
	log.Println("menjalankan auto-migration...")
	return db.AutoMigrate(
		&models.Restaurant{},
		&models.Table{},
		&models.TableSession{},
		&models.MenuCategory{},
		&models.MenuItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.StaffUser{},
		&models.Review{},
	)
}
