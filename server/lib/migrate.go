package lib

import (
	"log"

	"github.com/zamachnoi/viewthis/models"
	"gorm.io/gorm"
)

// when making changes to db, automigrate
func AutoMigrate(db *gorm.DB) {
    err := db.AutoMigrate(&models.User{}, &models.Queue{}, &models.Submission{}, &models.Feedback{}, &models.Guild{}, &models.GuildUser{})
    if err != nil {
        log.Printf("Failed to auto-migrate database: %v", err)
    }
}

func MigrateUp() {
    err := InitDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
}
func MigrateDown() {
    log.Default().Println("MigrateDown")
}