package model

type EventEmitDto struct {
	EventID   string      `json:"event_id"`
	EventName string      `json:"event_name"` // Could be registration_event, state_change_event, light_state_override, light_state_override_response
	EventData interface{} `json:"event_data"`
	EmittedAt int64       `json:"emitted_at"`
}

type EventReceiveDto struct {
	Status string `json:"status"` // Could be success, failed
	EventEmitDto
}

type EventCore struct {
	EventName string `gorm:"index" json:"event_name"`
	EmittedAt int64  `gorm:"index" json:"emitted_at"`
}

// A Data is also a DTO

// Use as EventData of RegistrationEvent
type RegistrationEventData struct {
	LightID  string `json:"light_id"`
	Location string `json:"location"`
}

type RegistrationEvent struct {
	ID        string `gorm:"primaryKey" json:"id"`
	LightID   string `gorm:"index" json:"light_id"`
	Location  string `gorm:"index" json:"location"`
	EventCore `gorm:"embedded"`
}

// Use as EventData of StateChangeEvent
type StateChangeData struct {
	LightID   string `json:"light_id"`
	FromState string `json:"from_state"`
	ToState   string `json:"to_state"`
}

type StateChangeEvent struct {
	ID        string `gorm:"primaryKey" json:"id"`
	LightID   string `gorm:"index" json:"light_id"`
	FromState string `gorm:"index" json:"from_state"`
	ToState   string `gorm:"index" json:"to_state"`
	EventCore `gorm:"embedded"`
}

// Use as EventData of LightStateOverrideEvent
type LightStateOverrideData struct {
	LightID string `json:"light_id"`
	ToState string `json:"to_state"`
}

type LightStateOverrideEvent struct {
	ID        string `gorm:"primaryKey" json:"id"`
	LightID   string `gorm:"index" json:"light_id"`
	ToState   string `gorm:"index" json:"to_state"`
	EventCore `gorm:"embedded"`
}

// Use as EventData of LightStateOverrideDoneEvent
type LightStateOverrideDoneData struct {
	LightID string `json:"light_id"`
}

type LightStateOverrideDoneEvent struct {
	ID      string `gorm:"primaryKey" json:"id"`
	LightID string `gorm:"index" json:"light_id"`
	EventCore
}

type LightState struct {
	LightID   string      `gorm:"primaryKey" json:"light_id"`
	Color     int         `gorm:"index" json:"color"` // RED=1,GREEN=2,YELLOW=3
	Location  interface{} `gorm:"type:geometry(Point,4326)"`
	ChangedAt int64       `gorm:"index,autoUpdateTime" json:"changed_at"`
}
