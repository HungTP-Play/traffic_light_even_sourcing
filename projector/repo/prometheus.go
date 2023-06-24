package repo

import (
	"fmt"
	"os"
	"projector/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func init() {
	db = Connect()
	// Enable PostGIS extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error; err != nil {
		panic("Failed to enable PostGIS extension")
	}

	AutoMigrate(db)
}

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
	db.AutoMigrate(&model.LightState{})

	fmt.Println("[[Projector]] Auto migration completed!! 🎉🎉🎉")
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

func UpsertTrafficLight(lightID string, lat, lon float64, color int) {
	db := GetDB()
	trafficLight := model.LightState{
		LightID:   lightID,
		Location:  gorm.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", lon, lat),
		Color:     color,
		ChangedAt: time.Now().Unix(),
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "light_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"location", "color", "changed_at"}),
	}).Create(&trafficLight)

	if result.Error != nil {
		fmt.Printf("Error upserting traffic light: %v\n", result.Error)
	}
}
