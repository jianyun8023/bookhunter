package telegram

import (
	"context"
	"time"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/log"
)

type (
	Telegram struct {
		channelID string
		mobile    string
		appID     int64
		appHash   string
		client    *telegram.Client
		ctx       context.Context
	}

	ChannelInfo struct {
		ID         int64
		AccessHash int64
		LastMsgID  int64
	}

	// File is the file info from the telegram channel.
	File struct {
		ID       int64
		Name     string
		Format   file.Format
		Size     int64
		Document *tg.InputDocumentFileLocation
	}
)

// New will create a telegram client.
func New(channelID, mobile string, appID int64, appHash string, sessionPath, proxy string) (*Telegram, error) {
	// Create the http proxy dial.
	dialFunc, err := createProxy(proxy)
	if err != nil {
		return nil, err
	}

	// Create the backend telegram client.
	client := telegram.NewClient(
		int(appID),
		appHash,
		telegram.Options{
			Resolver:       dcs.Plain(dcs.PlainOptions{Dial: dialFunc}),
			SessionStorage: &session.FileStorage{Path: sessionPath},
			Middlewares: []telegram.Middleware{
				floodwait.NewSimpleWaiter().WithMaxRetries(uint(3)),
			},
		},
	)

	ctx := context.Background()
	_, err = bg.Connect(client, bg.WithContext(ctx)) // No need to close this client.
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		ctx:       ctx,
		channelID: channelID,
		mobile:    mobile,
		appID:     appID,
		appHash:   appHash,
		client:    client,
	}

	if err := t.Authentication(); err != nil {
		return nil, err
	}

	t.keepAlive()

	return t, nil
}

// keepAlive will periodically send an API request (account.updateStatus) to the server
// to keep the session alive and update the last active time of the device.
func (t *Telegram) keepAlive() {
	go func() {
		// Immediately update the active time once upon startup.
		if _, err := t.client.API().AccountUpdateStatus(t.ctx, false); err != nil {
			log.Warn("Failed to send initial keep-alive request: ", err)
		}

		// Update every 5 minutes.
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-t.ctx.Done():
				return
			case <-ticker.C:
				if _, err := t.client.API().AccountUpdateStatus(t.ctx, false); err != nil {
					log.Warn("Failed to send keep-alive request: ", err)
				}
			}
		}
	}()
}
