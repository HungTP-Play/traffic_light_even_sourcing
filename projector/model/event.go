package model

type EventEmitDto struct {
	EventID   string      `json:"event_id"`
	EventName string      `json:"event_name"` // Could be registration_event, state_change_event, light_state_override, light_state_override_response
	EventData interface{} `json:"event_data"`
	EmittedAt int64       `json:"emitted_at"`
}

type LightState struct {
	LightID   string      `gorm:"primaryKey" json:"light_id"`
	Color     int         `gorm:"index" json:"color"` // RED=1,GREEN=2,YELLOW=3
	Location  interface{} `gorm:"type:geometry(Point,4326)"`
	ChangedAt int64       `gorm:"index,autoUpdateTime" json:"changed_at"`
}
