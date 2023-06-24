# Overview Architecture

![Architecture](https://i.ibb.co/sykzVs0/2023-06-23-01-35.png)

## Models

### Event Store Models

```go
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
 LightID  string `json:"light_id"`
 Location string `json:"location"`
 ToState  string `json:"to_state"`
}

type StateChangeEvent struct {
 ID        string `gorm:"primaryKey" json:"id"`
 LightID   string `gorm:"index" json:"light_id"`
 Location  string `gorm:"index" json:"location"`
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

```

### Projection Models

```go
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

```

### Controller Models

```go
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

```

## Queues (A.K.A. Channels)

- `Controller` Queue for control
- `Projections` Queue for projection

## Projections

### Current State of Traffic Lights

```sql
SELECT
  light_id as "id",
  ST_Y(location) as "latitude",
  ST_X(location) as "longitude",
  color as "value",
  last_updated as "time"
FROM light_states;
```

### Total number of traffic lights

```sql
SELECT COUNT(*) FROM light_states;
```

## Prometheus

### Scrape Config

Currently in `prometheus.yml` file

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 1m
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'projector'
    metrics_path: /metrics
    static_configs:
      - targets: ['projector:3333']

```

## Run Steps

1. Build all services

```bash
docker-compose build
```

2. Run all services

```bash
docker-compose up
```

3. Open `frontend` in browser

```bash
http://localhost:3000
```

4. Open `projection` (grafana) in browser

```bash
http://localhost:3001
```

User is `admin` and password is `admin`

Go to the `Dashboard` and select `Current Light States` dashboard
