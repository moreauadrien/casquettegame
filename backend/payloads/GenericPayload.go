package payloads

import (
	"encoding/json"
	"timesup/events"

	"github.com/go-playground/validator/v10"
)

var genericValidator = validator.New()
var credentialsValidator = validator.New()

type Credentials struct {
	Token    string `json:"token" validate:"required,min=3"`
	Id       string `json:"id" validate:"required,min=3"`
	Username string `json:"username" validate:"required,min=3"`
}

type GenericPayload struct {
	MessageId   string               `json:"messageId"`
	Event       *events.GenericEvent `json:"event" validate:"required"`
	Credentials *Credentials         `json:"credentials,omitempty"`
	raw         []byte
}

func (r GenericPayload) Marshal() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return b
}

func (p GenericPayload) Validate() error {
	return genericValidator.Struct(p)
}

func (p *GenericPayload) UnmarshalFrom(msg []byte) error {
	p.raw = msg

	if err := json.Unmarshal(msg, p); err != nil {
		return err
	}

	if err := p.Validate(); err != nil {
		return err
	}

	if p.Event.Type == "response" {
		respPayload := new(ResponsePayload)

		if err := json.Unmarshal(p.raw, respPayload); err != nil {
			return err
		}

		if err := respPayload.Validate(); err != nil {
			return err
		}
	} else {
		normalPayload := new(NormalEventPayload)

		if err := json.Unmarshal(p.raw, normalPayload); err != nil {
			return err
		}

		if err := normalPayload.Validate(); err != nil {
			return err
		}
	}

	return nil
}
