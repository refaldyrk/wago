package command

import (
	"context"
	"fmt"
	"net/url"
	"wago/api"
	"wago/log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func WikipediaCommand(client *whatsmeow.Client, ctx context.Context, to types.JID, query string) {
	help := api.Wikipedia(url.QueryEscape(query))
	_, _ = client.SendMessage(ctx, to, &proto.Message{Conversation: &help})
	log.LogMe("COMMAND", fmt.Sprintf("Wikipedia Command To: %s & Query %s", to, query))
}
