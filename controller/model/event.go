package model

import "time"

type TrafficLight struct {
	ID           uint   `gorm:"autoIncrement"`
	Location     string `gorm:"not null;size:255"`
	CurrentState string `gorm:"not null;size:255"`
}

type EventCore struct {
	EventName string    `gorm:"not null;size:255"`
	EmittedAt time.Time `gorm:"not null"`
}

type RegistrationEvent struct {
	ID       uint   `gorm:"autoIncrement"`
	LightID  uint   `gorm:"not null"`
	Location string `gorm:"not null;size:255"`
	EventCore
}

// After receive this event, the traffic light state is already changed
type StateChangeEvent struct {
	ID        uint   `gorm:"autoIncrement"`
	LightID   uint   `gorm:"not null"`
	FromState string `gorm:"not null;size:255"`
	ToState   string `gorm:"not null;size:255"`
	EventCore
}

// After receive this event, the traffic light is command to change but not sure if it's changed
type LightStateOverride struct {
	ID      uint   `gorm:"autoIncrement"`
	LightID uint   `gorm:"not null"`
	State   string `gorm:"not null;size:255"`
	EventCore
}

// Response to LightStateOverride event that the traffic light is changed
type LightStateOverrideResponse struct {
	ID      uint   `gorm:"autoIncrement"`
	LightID uint   `gorm:"not null"`
	State   string `gorm:"not null;size:255"`
	EventCore
}
