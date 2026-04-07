package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/plugin/opentelemetry/tracing"
)

var DB *gorm.DB

func ConnectDatabase() {

	db, err := gorm.Open(sqlite.Open("hotel.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Enable OpenTelemetry tracing
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatal(err)
	}

	DB = db
}