package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var (
	EmptyLocationErr = fmt.Errorf("message without location")
	EmptyFromErr     = fmt.Errorf("message without user information")
	EmptyChatErr     = fmt.Errorf("message without chat information")
	EmptyMessageErr  = fmt.Errorf("update without message & edited message")
)

type Bot struct {
	cl   *tgbotapi.BotAPI
	LocC chan UserLocation
	ComC chan Command
	RegC chan RegInfo
}

func New(cl *tgbotapi.BotAPI) *Bot {
	return &Bot{
		cl:   cl,
		ComC: make(chan Command),
		LocC: make(chan UserLocation),
		RegC: make(chan RegInfo),
	}
}

func (b *Bot) Listen() error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // TODO: extract
	updates, err := b.cl.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		loc, err := b.GetLocation(update)
		if err != nil {
			if update.Message != nil {
				switch update.Message.Text {
				case "/start":
					err := b.Send(update.Message.Chat.ID, startMsg)
					if err != nil {
						log.Println(err)
						continue
					}
					b.ComC <- Command{
						ChatID:   update.Message.Chat.ID,
						UserName: update.Message.From.UserName,
						Text:     "/start",
					}
				case "/reg":
					b.ComC <- Command{
						ChatID: update.Message.Chat.ID,
						Text:   "/reg",
					}
				case "/del":
				default:
					b.ComC <- Command{
						ChatID: update.Message.Chat.ID,
						Text:   update.Message.Text,
					}
				}
			}
			continue
		}
		b.LocC <- *loc
	}
	return nil
}

type UserLocation struct {
	UserName string
	ChatID   int64
	Location *tgbotapi.Location
}

type RegInfo struct {
	Name   string
	ChatID int64
}

type Command struct {
	ChatID   int64
	UserName string
	Text     string
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

func (b *Bot) Send(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.cl.Send(msg)
	return err
}

const startMsg = `Добро пожаловать в протестный бот.
Отправьте трансляцию своей локации в месте скопления граждан и переодически проверяйте поступающую от нас информацию.
Вы также можете "зарегистрироваться" (команда /reg) и указать несколько доверенных контактов, которым будет сообщено в случае обнаружения вас вне "горячей" зоны, что подразумевает задержание и транспортировку в ОВД.
Команда /del удалит имеющуюся информацию`
