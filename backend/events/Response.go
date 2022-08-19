package events

type ResponseEvent struct {
	Type string    `json:"type" validate:"required,min=1"`
	To   string    `json:"to" validate:"required,min=1"`
	Data EventData `json:"data"`
}
