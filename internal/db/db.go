package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Dhairya3124/coupon-system-go/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() *DB {
	dbURL := os.Getenv("DB_URI")
	if dbURL == "" {
		dbURL = "host=localhost user=admin password=adminpassword dbname=couponsystem port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if db.Migrator().HasTable(&model.Coupon{}) {
		log.Println("Dropping existing coupons table...")
		if err := db.Migrator().DropTable(&model.Coupon{}); err != nil {
			fmt.Printf("failed to drop tables: %v", err)
		}
	}

	log.Println("Creating coupons table...")
	if err := db.AutoMigrate(&model.Coupon{}); err != nil {
		fmt.Printf("failed to create coupons table: %v", err)
	}
	log.Println("Coupons table created successfully")

	return &DB{DB: db}
}
