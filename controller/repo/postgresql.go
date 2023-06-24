package repo

import (
	"controller/model"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetPostgresDSN() string {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")
	return "host=" + postgresHost + " port=" + postgresPort + " user=" + postgresUser + " password=" + postgresPassword + " dbname=" + postgresDatabase + " sslmode=disable"
}

func Connect() *gorm.DB {
	dsn := GetPostgresDSN()
	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("[[Event Store]] failed to connect database")
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&model.TrafficLight{})
	fmt.Println("[[Controller]] Auto migration completed!! 🎉🎉🎉")
}

func GetDB() *gorm.DB {
	if db == nil {
		db = Connect()
		AutoMigrate(db)
	}

	return db
}

func CloseDB() {
	if db != nil {
		dbb, _ := db.DB()
		dbb.Close()
	}
}

func StoreTrafficLight(light model.TrafficLight) error {
	db := GetDB()
	db.Create(&light)
	fmt.Println("[[Controller]] Stored traffic light: ", light)
	return nil
}
