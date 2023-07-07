package events

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}

func (mycli *MyClient) Register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.EventHandler)
}

func (mycli *MyClient) EventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		MessageHandler(v, mycli.WAClient)
	}
}

func NewEventsHandler(ext MyClient) *MyClient {
	return &MyClient{WAClient: ext.WAClient}
}
