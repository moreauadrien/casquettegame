package payloads

import (
	"encoding/json"
	"timesup/events"

	"github.com/go-playground/validator/v10"
)

var responseValidator = validator.New()

type ResponsePayload struct {
	Event events.ResponseEvent `json:"event"`
}

func (p ResponsePayload) Validate() error {
	return responseValidator.Struct(p)
}

func (r ResponsePayload) Marshal() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return b
}
