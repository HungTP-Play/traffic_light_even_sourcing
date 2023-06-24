package repo

import (
	"fmt"
	"os"

	"event_store/model"

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
	db.AutoMigrate(&model.RegistrationEvent{})
	db.AutoMigrate(&model.StateChangeEvent{})
	db.AutoMigrate(&model.LightStateOverrideEvent{})
	db.AutoMigrate(&model.LightStateOverrideDoneEvent{})

	fmt.Println("[[Event Store]] Auto migration completed!! 🎉🎉🎉")
}

func GetDB() *gorm.DB {
	if db == nil {
		db = Connect()
		AutoMigrate(db)
	}

	return db
}

func CloseDB() {
	db, _ := db.DB()
	db.Close()
}

func StoreEvent(event model.EventEmitDto) error {
	db := GetDB()
	defer CloseDB()

	fmt.Println("[[Event Store]] Storing event: ", event)
	// Map event to model by event_name
	switch event.EventName {
	case "registration_event":
		registrationEvent := model.RegistrationEvent{
			ID:        event.EventID,
			LightID:   event.EventData.(map[string]interface{})["light_id"].(string),
			Location:  event.EventData.(map[string]interface{})["location"].(string),
			EventCore: model.EventCore{EventName: event.EventName, EmittedAt: event.EmittedAt},
		}

		db.Create(&registrationEvent)
	case "state_change_event":
		stateChangeEvent := model.StateChangeEvent{
			ID:        event.EventID,
			LightID:   event.EventData.(map[string]interface{})["light_id"].(string),
			FromState: event.EventData.(map[string]interface{})["from_state"].(string),
			ToState:   event.EventData.(map[string]interface{})["to_state"].(string),
			EventCore: model.EventCore{EventName: event.EventName, EmittedAt: event.EmittedAt},
		}

		db.Create(&stateChangeEvent)
	case "light_state_override":
		lightStateOverrideEvent := model.LightStateOverrideEvent{
			ID:        event.EventID,
			LightID:   event.EventData.(map[string]interface{})["light_id"].(string),
			ToState:   event.EventData.(map[string]interface{})["to_state"].(string),
			EventCore: model.EventCore{EventName: event.EventName, EmittedAt: event.EmittedAt},
		}

		db.Create(&lightStateOverrideEvent)

	case "light_state_override_response":
		lightStateOverrideDoneEvent := model.LightStateOverrideDoneEvent{
			ID:        event.EventID,
			LightID:   event.EventData.(map[string]interface{})["light_id"].(string),
			EventCore: model.EventCore{EventName: event.EventName, EmittedAt: event.EmittedAt},
		}

		db.Create(&lightStateOverrideDoneEvent)
	}

	fmt.Printf("[[Event Store]] Event %s stored successfully!! 🎉🎉🎉\n", event.EventName)
	fmt.Printf("[[Event Store]] Event stored data %v\n", event)
	return nil
}
