package boot

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	log.Println("Connecting to database...")
	start := time.Now()

	uri := os.Getenv("SQLITE_DB_FILE")
	if uri == "" {
		log.Fatal("DATABASE_URI environment variable not set")
	}

	DB, err = gorm.Open(sqlite.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	duration := time.Since(start)
	log.Printf("Connected to database in %s\n", duration)
}
