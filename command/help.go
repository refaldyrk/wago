package command

import (
	"context"
	"fmt"
	"time"
	"wago/log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func HelpCommand(client *whatsmeow.Client, ctx context.Context, to types.JID) {
	help := fmt.Sprintf("REFAL SELFBOT\n\nTIME: %s\n\nCOMMAND:\n1. hi -> tagall\n2. wiki [argument]\n3. login\n4. logout\n\nTHANKS", time.Now().Format("2006-01-02 15:04:05"))
	_, _ = client.SendMessage(ctx, to, &proto.Message{Conversation: &help})
	log.LogMe("COMMAND", fmt.Sprintf("Help Command To: %s", to))
}
