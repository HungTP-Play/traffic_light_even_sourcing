# Overview Architecture

![Architecture](https://i.ibb.co/sykzVs0/2023-06-23-01-35.png)

## Models

```go
type TrafficLight struct {
  gorm:"primary_key:id"
  ID        uint      `gorm:"autoIncrement"`
  Location  string    `gorm:"not null;size:255"`
  CurrentState  string    `gorm:"not null;size:255"`
}

type RegistrationEvent struct {
  gorm:"primary_key:id"
  ID      uint      `gorm:"autoIncrement"`
  LightID   uint      `gorm:"not null"` 
  Location  string    `gorm:"not null;size:255"`
  EmittedAt time.Time `gorm:"not null"`
}

type StateChangeEvent struct {
  gorm:"primary_key:id"
  ID      uint      `gorm:"autoIncrement"`
  LightID   uint      `gorm:"not null"` 
  FromState   string    `gorm:"not null;size:255"`
  ToState   string    `gorm:"not null;size:255"`
  EmittedAt time.Time `gorm:"not null"`
}

type LightStateOverride struct {
  gorm:"primary_key:id"
  ID      uint      `gorm:"autoIncrement"`   
  LightID   uint      `gorm:"not null"`
  State     string    `gorm:"not null;size:255"` 
  CommandedAt time.Time `gorm:"not null"` 
}
```

- `TrafficLight` represents a single Traffic Light, with its ID, location and current state.

- `StateChangeEvent` represents an event emitted when a Light changes state, with the `FromState` and `ToState`, and timestamp.

- `LightStateOverride` represents a command from the Controller to a Light to change state, with the commanded state and timestamp.

Then in your code, you could do:

- Have a TrafficLight GORM model control its own state, updating `CurrentState` and emitting `StateChangeEvent`s as it transitions.

- The Controller could query `TrafficLight` to see the current state of all Lights.

- The Controller could insert `LightStateOverride` records to command a Light to change state. The Light would check for any outstanding commands before changing state on its own.

- The Event Store and Projector listen for `StateChangeEvent`s and update their stores appropriately.

## Queues

- `Metadata` Queue for Light registration
- `StateChange` Queue for Light state changes

## Projections

Projections for the control team would be:

1. Current state of all lights:

```sql
SELECT * FROM TrafficLights;
```

This would show the control team the current state (red, green, yellow) of every light for monitoring.

2. State change history for a given light:

```sql
SELECT * FROM StateChangeEvents 
WHERE LightID = {light_id}
ORDER BY EmittedAt DESC;
```

This would show the timeline of state changes for a particular light, useful for auditing or troubleshooting.

3. List of pending light override commands:

```sql
SELECT * FROM LightStateOverrides
WHERE CommandedAt > (NOW() - INTERVAL '5' MINUTE) 
AND (SELECT CurrentState FROM TrafficLights WHERE ID = LightID) != State; 
```

This would show any light override commands from the past 5 minutes that have not been implemented yet by the lights. Helpful to ensure all commands are processed properly.

4. Override command history for a light:

```sql
SELECT * FROM LightStateOverrides
WHERE LightID = {light_id}
ORDER BY CommandedAt DESC;
```

This would show the history of all light override commands for a particular light, useful for auditing purposes.

5. List of lights currently in a given state (e.g. red lights):

```sql
SELECT * FROM TrafficLights 
WHERE CurrentState = 'red';
```

This would allow the control team to quickly see all lights in a particular state for monitoring.

6. Count of lights in each state:

```sql
SELECT CurrentState, COUNT(*) as Count 
FROM TrafficLights 
GROUP BY CurrentState;
```
