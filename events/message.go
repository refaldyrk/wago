package events

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"wago/command"
	"wago/log"

	"github.com/mdp/qrterminal"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(v *events.Message, client *whatsmeow.Client) {
	if v.Info.Sender.String() == "6288809462517@s.whatsapp.net" || v.Info.Sender.String() == "6283899673331@s.whatsapp.net" {
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
			fmt.Println("Login Disini")
			err = os.Remove(qr)
			if err != nil {
				fmt.Println(err)
			}
			return
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
			fmt.Println("Login Disana")
			err = os.Remove(qr)
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
		case "reconnect":
			client.Disconnect()
			err := client.Connect()
			if err != nil {
				fmt.Errorf("Failed to connect: %v", err)
			}
			return
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
			os.Remove(fmt.Sprintf("./whatsapp%s.db", v.Info.Chat.String()))
			return
		case "getlink":
			if !strings.Contains(v.Info.Chat.String(), "@g.us") {
				message := "Only In Group"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				return
			}
			if len(args) > 1 {
				go func() {}()
				groups, _ := client.GetJoinedGroups()
				remotes, _ := strconv.Atoi(args[1])
				group, _ := client.GetGroupInviteLink(groups[remotes-1].JID, false)
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &group})
				return
			}
			link, _ := client.GetGroupInviteLink(v.Info.Chat, false)
			client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &link})
			return
			break
		case "mygroup":
			go func() {
				var nameGroup []string
				groups, _ := client.GetJoinedGroups()
				for _, v := range groups {
					nameGroup = append(nameGroup, v.Name)
				}
				var builder strings.Builder
				builder.WriteString("")
				for i, elem := range nameGroup {
					builder.WriteString(fmt.Sprintf("%d.%s\n", i+1, elem))
				}

				joinedString := builder.String()
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &joinedString})
				return
			}()
			return
		case "leaveremote":
			if len(args) < 2 {
				joinedString := "Please Give Me Argument"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &joinedString})
				return
			}
			go func() {
				groups, _ := client.GetJoinedGroups()
				remotes, _ := strconv.Atoi(args[1])
				message := fmt.Sprintf("Success Leave Group %s", groups[remotes-1].Name)
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				client.LeaveGroup(groups[remotes-1].JID)
			}()
			return
		case "inforemote":
			if len(args) < 2 {
				group, _ := client.GetGroupInfo(v.Info.Chat)
				message := fmt.Sprintf("INFO GROUP\n\nName: %s\nName Set By: %s\nCreated At: %s\nParticipant: %d\nOwner: %s", group.Name, group.NameSetBy.User, group.GroupCreated, len(group.Participants), group.OwnerJID.User)
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				return
			} else {
				go func() {
					groups, _ := client.GetJoinedGroups()
					remotes, _ := strconv.Atoi(args[1])
					group, _ := client.GetGroupInfo(groups[remotes-1].JID)
					message := fmt.Sprintf("INFO GROUP\n\nName: %s\n\nName Set By: wa.me/%s\n\nCreated At: %s\n\nParticipant: %d\n\nOwner: wa.me/@%s", group.Name, group.NameSetBy.User, group.GroupCreated, len(group.Participants), group.OwnerJID.User)
					client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				}()
				return
			}
			return
		case "fakeremote":
			if len(args) < 2 {
				message := "no argument"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				return
			}
			go func() {
				groups, _ := client.GetJoinedGroups()
				remotes, _ := strconv.Atoi(args[1])
				group, _ := client.GetGroupInfo(groups[remotes-1].JID)
				var participants []string
				for _, v := range group.Participants {
					participants = append(participants, v.JID.String())
				}
				message := "🙂"
				client.SendMessage(context.Background(), group.JID, &proto.Message{ExtendedTextMessage: &proto.ExtendedTextMessage{Text: &message, ContextInfo: &proto.ContextInfo{MentionedJid: participants}}})
			}()
			return
		case "creategc":
			if len(args) < 2 {
				message := "no argument"
				client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
				return
			}
			newGroup, _ := client.CreateGroup(whatsmeow.ReqCreateGroup{Name: args[1]})
			link, _ := client.GetGroupInviteLink(newGroup.JID, false)
			client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &link})
			return
		case "spamtext":
			go func() {
				if len(args) < 2 {
					message := "no argument"
					client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
					return
				}
				if args[1] == "" && args[2] == "" {
					message := "no argument"
					client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &message})
					return
				}
				//apmas := strings.Split(v.Message.GetConversation(), "|")
				// Mencocokkan teks setelah "/spamtext" dan sebelum "|"
				re := regexp.MustCompile(`/spamtext\s*(.*?)\s*\|`)
				match := re.FindStringSubmatch(v.Message.GetConversation())
				var result string
				if len(match) > 1 {
					result = match[1]
					fmt.Println(result) // Output: halo semuanya
				}

				// Mencocokkan teks setelah "|"
				re = regexp.MustCompile(`\|\s*(.*?)\s*$`)
				match = re.FindStringSubmatch(v.Message.GetConversation())

				afterPipe := ""
				if len(match) > 1 {
					afterPipe = match[1]
					fmt.Println(afterPipe) // Output: 10
				}
				jumlah, _ := strconv.Atoi(afterPipe)
				for i := 1; i < int(jumlah); i++ {
					go func() {
						client.SendMessage(context.Background(), v.Info.Chat, &proto.Message{Conversation: &result})
					}()
				}
			}()

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
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
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

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			fmt.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}
