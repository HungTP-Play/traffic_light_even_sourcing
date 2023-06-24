package model

type EventEmitDto struct {
	EventID   string      `json:"event_id"`
	EventName string      `json:"event_name"`
	EventData interface{} `json:"event_data"`
	EmittedAt int64       `json:"emitted_at"`
}

// This is model also an DTO
type TrafficLight struct {
	LightID      string `gorm:"primaryKey" json:"light_id"`
	Location     string `gorm:"index" json:"location"`
	RegisteredAt int64  `gorm:"index,autoCreateTime" json:"registered_at"`
}
