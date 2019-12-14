package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	EmptyLocationErr = fmt.Errorf("message without location")
	EmptyFromErr     = fmt.Errorf("message without user information")
	EmptyChatErr     = fmt.Errorf("message without chat information")
	EmptyMessageErr  = fmt.Errorf("update without message & edited message")
)

type Bot struct {
	cl *tgbotapi.BotAPI
}

func New(cl *tgbotapi.BotAPI) *Bot {
	return &Bot{
		cl: cl,
	}
}

func (b *Bot) Listen(exit chan UserLocation) error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // TODO: extract
	updates, err := b.cl.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		loc, err := b.GetLocation(update)
		if err != nil {
			continue
		}
		exit <- *loc
	}
	return nil
}

type UserLocation struct {
	UserName string
	ChatID   int64
	Location *tgbotapi.Location
}

func (b *Bot) GetLocation(u tgbotapi.Update) (*UserLocation, error) {
	if u.Message == nil {
		if u.EditedMessage == nil {
			return nil, EmptyMessageErr
		}
		return b.getLocation(u.EditedMessage)
	}
	return b.getLocation(u.Message)
}

func (b *Bot) getLocation(m *tgbotapi.Message) (*UserLocation, error) {
	if m.Location == nil {
		return nil, EmptyLocationErr
	}

	if m.From == nil {
		return nil, EmptyFromErr
	}

	if m.Chat == nil {
		return nil, EmptyChatErr
	}

	return &UserLocation{
		UserName: m.From.UserName,
		ChatID:   m.Chat.ID,
		Location: m.Location,
	}, nil
}
