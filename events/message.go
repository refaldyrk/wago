package events

import (
	"context"
	"fmt"
	"os"
	"strings"
	"wago/command"
	"wago/log"

	"github.com/mdp/qrterminal"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(v *events.Message, client *whatsmeow.Client) {
	if v.Info.Sender.String() == "YOUR_WHATSAPP" || v.Info.Sender.String() == "YOUR_HELPER" {
		switch v.Message.GetConversation() {
		case "/login":
			fmt.Println("Login")
			qr := Login(v.Info.Chat.String())
			client.SendMessage(context.Background(), v.Info.Sender, &proto.Message{ImageMessage: &proto.ImageMessage{FileSha256: qr}})
			break
		}
	}
	if v.Info.IsFromMe {
		log.LogMe("SEND MESSAGE", fmt.Sprintf("%s: %s -> %s\n", v.Info.Sender, v.Message.GetConversation(), v.Info.Chat))
		commandMessage := command.GetCommand(v.Message.GetConversation())
		args := strings.Split(commandMessage, " ")
		switch args[0] {
		case "help":
			command.HelpCommand(client, context.Background(), v.Info.Chat)
			break
		case "wiki":
			if len(args) == 1 {
				message := "You No Have Arg"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				return
			}
			query := strings.Join(args[1:], " ")
			command.WikipediaCommand(client, context.Background(), v.Info.Chat, query)
			break
		case "logout":
			log.LogMe("LOGOUT", client.Store.PushName)
			message := "Success Logout..."
			client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
			client.Logout()
			break
		}
	} else {
		log.LogMe("RECEIVE MESSAGE", fmt.Sprintf("%s: %s -> %s\n", v.Info.Sender, v.Message.GetConversation(), v.Info.Chat))
	}
}

func Login(id string) []byte {
	//dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite", fmt.Sprintf("file:whatsapp%s.db?_foreign_keys=on&_pragma=busy_timeout=10000", id), nil)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	//clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, nil)

	//Create Handler
	e := NewEventsHandler(MyClient{WAClient: client})
	client.AddEventHandler(e.EventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				var png []byte
				png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
				if err != nil {
					log.LogMe("LOGIN", "Gagal Login "+err.Error())
				}
				return png
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		log.LogMe("LOGIN", fmt.Sprintf("%s", client.Store.PushName))
		if err != nil {
			panic(err)
		}
	}

	return []byte{}
}
