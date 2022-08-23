package events

import (
	"encoding/json"
	"reflect"
	"timesup/structures"
)

type EventData interface{}

type TeamColor string

type PlayerInfos struct {
	Username string    `json:"username"`
	Id       string    `json:"id"`
	Team     TeamColor `json:"team"`
}

type ResponseData struct {
	Status  string               `json:"status"`
	Message string               `json:"message,omitempty"`
	RoomId  string               `json:"roomId,omitempty"`
	Players []PlayerInfos        `json:"players,omitempty"`
	Host    string               `json:"host,omitempty"`
	Team    TeamColor            `json:"team,omitempty"`
	Cards   *structures.CardPile `json:"cards,omitempty"`
}

type JoinData struct {
	RoomId string `json:"roomId"`
}

type TurnUpdate struct {
	Cards []string `json:"cards,omitempty"`
}

type StateUpdateData struct {
	State   string        `json:"state"`
	Players []PlayerInfos `json:"players,omitempty"`
	Speaker *PlayerInfos  `json:"speaker,omitempty"`
	Cards   []string      `json:"cards,omitempty"`
}

type GenericEvent struct {
	Type string    `json:"type" validate:"required"`
	Data EventData `json:"data"`
}

type ResponseEvent struct {
	Type string       `json:"type" validate:"required,min=1"`
	To   string       `json:"to" validate:"required,min=1"`
	Data ResponseData `json:"data"`
}

func (e *GenericEvent) UnmarshalJSON(data []byte) error {
	typeName, value, err := UnmarshalCustomValue(data, "type", "data", map[string]reflect.Type{
		"join": reflect.TypeOf(JoinData{}),
	})

	if err != nil {
		return err
	}

	e.Type = typeName
	e.Data = value

	return nil
}

func UnmarshalCustomValue(data []byte, typeJsonField, valueJsonField string, customTypes map[string]reflect.Type) (string, EventData, error) {
	m := map[string]interface{}{}

	if err := json.Unmarshal(data, &m); err != nil {
		return "", nil, err
	}

	typeName := m[typeJsonField].(string)
	var value EventData
	if ty, found := customTypes[typeName]; found {
		value = reflect.New(ty).Interface().(EventData)
	}

	valueBytes, err := json.Marshal(m[valueJsonField])
	if err != nil {
		return "", nil, err
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return "", nil, err
	}

	return typeName, value, nil
}
