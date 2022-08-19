package events

type EventData map[string]string

type GenericEvent struct {
	Type string    `json:"type" validate:"required"`
	Data EventData `json:"data"`
}
