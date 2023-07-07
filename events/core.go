package events

import (
	"fmt"

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
		fmt.Println("\n\nMESSAGE: ", v)
	case *events.AppState:
		fmt.Println("\n\nAPP STATE:", v)
	case *events.Presence:
		fmt.Println("\n\nPRESENCE:", v)
	case *events.ChatPresence:
		fmt.Println("\n\nCHAT PRESENCE:", v)
	case *events.UnknownCallEvent:
		fmt.Println("\n\nUKNOWN:", v)
	case *events.IdentityChange:
		fmt.Println("\n\nIDENTITY CHANGE:", v)
	case *events.GroupInfo:
		fmt.Println("\n\nGROUP INFO:", v)
	case *events.JoinedGroup:
		fmt.Println("\n\nJOIN GROUP:", v)
	case *events.DeleteChat:
		fmt.Println("\n\nDELETE CHAT:", v)
	case *events.DeleteForMe:
		fmt.Println("\n\nDELETE MESSAGE: ", v)
	}
}

func NewEventsHandler(ext MyClient) *MyClient {
	return &MyClient{WAClient: ext.WAClient}
}
