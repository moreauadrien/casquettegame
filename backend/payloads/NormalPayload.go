package payloads

import (
	"encoding/json"
	"timesup/events"

	"github.com/go-playground/validator/v10"
)

var normalValidator = validator.New()

type NormalEventPayload struct {
	Credentials Credentials         `json:"credentials" validate:"required"`
	MessageId   string              `json:"messageId" validate:"required,min=1"`
	Event       events.GenericEvent `json:"event" validate:"required"`
}

func (p NormalEventPayload) Validate() error {
	return normalValidator.Struct(p)
}
func (r NormalEventPayload) Marshal() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return b
}
