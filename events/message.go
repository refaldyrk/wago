package events

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"wago/command"
	"wago/log"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(v *events.Message, client *whatsmeow.Client) {
	if v.Info.Sender.String() == "6288809462517@s.whatsapp.net" || v.Info.Sender.String() == "6288809462517@s.whatsapp.net" {
		arg := strings.Split(v.Message.ExtendedTextMessage.GetText(), " ")
		switch arg[0] {
		case "/login":
			fmt.Println("Login")
			var jid string
			if v.Message.ExtendedTextMessage.ContextInfo.MentionedJid != nil || len(v.Message.ExtendedTextMessage.ContextInfo.MentionedJid) == 0 {
				jid = v.Message.ExtendedTextMessage.ContextInfo.MentionedJid[0]
			} else {
				jid = v.Info.Chat.String()
			}
			qr := Login(jid)
			jpegImageFile, jpegErr := os.Open(qr)
			if jpegErr != nil {
				fmt.Println(jpegErr)
			}
			defer jpegImageFile.Close()

			jpegFileinfo, _ := jpegImageFile.Stat()
			var jpegSize int64 = jpegFileinfo.Size()
			jpegBytes := make([]byte, jpegSize)

			jpegBuffer := bufio.NewReader(jpegImageFile)
			_, jpegErr = jpegBuffer.Read(jpegBytes)

			resp, err := client.Upload(context.Background(), jpegBytes, whatsmeow.MediaImage)
			if err != nil {
				fmt.Println(err)
			}
			mimetyoe := "image/jpeg"

			imageMsg := &proto.ImageMessage{
				Mimetype: &mimetyoe, // replace this with the actual mime type
				// you can also optionally add other fields like ContextInfo and JpegThumbnail here
				ThumbnailDirectPath: &resp.DirectPath,
				ThumbnailSha256:     resp.FileSHA256,
				ThumbnailEncSha256:  resp.FileEncSHA256,
				JpegThumbnail:       jpegBytes,

				Url:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSha256: resp.FileEncSHA256,
				FileSha256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}

			_, err = client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{ImageMessage: imageMsg})
			if err != nil {
				fmt.Println(err)
			}
			break
		}
		switch v.Message.GetConversation() {
		case "/login":
			fmt.Println("Login")
			qr := Login(v.Info.Chat.String())
			jpegImageFile, jpegErr := os.Open(qr)
			if jpegErr != nil {
				fmt.Println(jpegErr)
			}
			defer jpegImageFile.Close()

			jpegFileinfo, _ := jpegImageFile.Stat()
			var jpegSize int64 = jpegFileinfo.Size()
			jpegBytes := make([]byte, jpegSize)

			jpegBuffer := bufio.NewReader(jpegImageFile)
			_, jpegErr = jpegBuffer.Read(jpegBytes)

			resp, err := client.Upload(context.Background(), jpegBytes, whatsmeow.MediaImage)
			if err != nil {
				fmt.Println(err)
			}
			mimetyoe := "image/jpeg"

			imageMsg := &proto.ImageMessage{
				Mimetype: &mimetyoe, // replace this with the actual mime type
				// you can also optionally add other fields like ContextInfo and JpegThumbnail here
				ThumbnailDirectPath: &resp.DirectPath,
				ThumbnailSha256:     resp.FileSHA256,
				ThumbnailEncSha256:  resp.FileEncSHA256,
				JpegThumbnail:       jpegBytes,

				Url:           &resp.URL,
				DirectPath:    &resp.DirectPath,
				MediaKey:      resp.MediaKey,
				FileEncSha256: resp.FileEncSHA256,
				FileSha256:    resp.FileSHA256,
				FileLength:    &resp.FileLength,
			}

			_, err = client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{ImageMessage: imageMsg})
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}
	if v.Info.IsFromMe {
		log.LogMe("SEND MESSAGE", fmt.Sprintf("%s: %s -> %s\n", v.Info.Sender, v.Message.GetConversation(), v.Info.Chat))
		commandMessage := command.GetCommand(v.Message.GetConversation())
		switch v.Message.GetConversation() {
		case "hi":
			if !strings.Contains(v.Info.Chat.String(), "@g.us") {
				message := "Only In Group"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
			}
			group, _ := client.GetGroupInfo(v.Info.Chat)
			var participants []string
			for _, v := range group.Participants {
				participants = append(participants, v.JID.String())
			}
			message := "🤷‍♂️"
			client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{ExtendedTextMessage: &proto.ExtendedTextMessage{Text: &message, ContextInfo: &proto.ContextInfo{MentionedJid: participants}}})
		}
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

func Login(id string) string {
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
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr"+id+".png")
				if err != nil {
					log.LogMe("LOGIN", "Gagal Login "+err.Error())
				}
				return "./" + "qr" + id + ".png"
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

	return ""
}
