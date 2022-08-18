package game

import "encoding/json"

type EventData interface {
	Wrap() Event
}

type Event struct {
	Type string    `json:"type"`
	Data EventData `json:"data"`
}

func (e Event) Marshal() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return b
}

type JoinEvent struct {
	Players []PlayerInfos `json:"players"`
	Host    string        `json:"host"`
}

func (e JoinEvent) Wrap() Event {
	return Event{
		Type: "playerJoin",
		Data: e,
	}
}

type InfosEvent struct {
	Players []PlayerInfos `json:"players"`
	Host    string        `json:"host"`
}

func (e InfosEvent) Wrap() Event {
	return Event{
		Type: "infos",
		Data: e,
	}
}
