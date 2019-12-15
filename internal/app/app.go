package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xelaj/mcclane/internal/api"
	"github.com/xelaj/mcclane/internal/db/pg"
	"github.com/xelaj/mcclane/internal/model"
	"github.com/xelaj/mcclane/internal/telegram"
	"log"
	"net/http"
	"strings"
	"time"
)

type App struct {
	bot *telegram.Bot
	db  *pg.DB
}

func NewApp(bot *telegram.Bot, db *pg.DB) *App {
	return &App{
		db:  db,
		bot: bot,
	}
}

func (app *App) Work(ctx context.Context) error {

	waitForReg := make(map[int64]struct{})
	waitForContacts := make(map[int64]struct{})
	waitForAnswer := make(map[int64]struct{})

	go app.bot.Listen()

	go func() {
		for {
			select {
			case loc := <-app.bot.LocC:
				uc, err := app.db.GetUserChatByChatID(ctx, loc.ChatID)
				if err != nil {
					log.Println(err)
					continue
				}
				if uc.ID == 0 {
					_, err = app.db.AddUserChat(ctx, model.UserChat{
						ChatID:   loc.ChatID,
						UserName: loc.UserName,
					})
				}
				hl, err := app.db.GetHotLocationByPoint(ctx, loc.Location.Latitude, loc.Location.Longitude)
				if err != nil {
					hl = &model.HotLocation{}
					err := app.db.UpdateUserChatWarning(ctx, true, uc.ID)
					if err != nil {
						log.Println(err)
					}
				} else {
					err := app.db.UpdateUserChatWarning(ctx, false, uc.ID)
					if err != nil {
						log.Println(err)
					}
				}
				_, err = app.db.AddCoordinates(ctx, model.Coordinates{
					CreatedAt:     time.Now(),
					UserChatID:    uc.ID,
					Latitude:      loc.Location.Latitude,
					Longitude:     loc.Location.Longitude,
					HotLocationID: hl.ID,
				})
				if err != nil {
					log.Println(err)
					continue
				}
			case c := <-app.bot.ComC:
				switch c.Text {
				case "/start":
					_, err := app.db.AddUserChat(ctx, model.UserChat{
						UserName: c.UserName,
						ChatID:   c.ChatID,
					})
					if err != nil {
						log.Println(err)
					}
				case "/reg":
					err := app.bot.Send(c.ChatID, "Отправьте свои ФИО")
					if err != nil {
						log.Println(err)
						continue
					}
					waitForReg[c.ChatID] = struct{}{}
				default:
					if _, ok := waitForReg[c.ChatID]; ok {
						err := app.db.UpdateUserChatName(ctx, c.ChatID, c.Text)
						if err != nil {
							log.Println(err)
						}
						delete(waitForReg, c.ChatID)
						_ = app.bot.Send(c.ChatID, "Оставьте телеграм-аккаунты, кому сообщить в случае, если связь с Вами пропадет")
						waitForContacts[c.ChatID] = struct{}{}
					}

					if _, ok := waitForContacts[c.ChatID]; ok {
						contacts := strings.Split(strings.TrimPrefix(strings.Trim(c.Text, " "), "@"), " ")
						uc, err := app.db.GetUserChatByChatID(ctx, c.ChatID)
						if err != nil {
							log.Println(err)
							continue
						}
						for _, contact := range contacts {
							_, _ = app.db.AddContact(ctx, model.Contacts{
								UserChatID: uc.ID,
								Contact:    strings.Trim(contact, " "),
							})
						}

						delete(waitForContacts, c.ChatID)
					}

					if _, ok := waitForAnswer[c.ChatID]; ok {
						uc, err := app.db.GetUserChatByChatID(ctx, c.ChatID)
						if err != nil {
							log.Println(err)
							continue
						}
						switch c.Text {
						case "нет", "Нет":
							cc, err := app.db.GetContactsByUserChatID(ctx, uc.ID)
							if err != nil {
								log.Println(err)
								continue
							}
							for _, contact := range cc {
								cnt, err := app.db.GetUserChatByUserName(ctx, contact.Contact)
								if err != nil {
									log.Println(err)
									continue
								}
								app.bot.Send(cnt.ChatID, fmt.Sprintf("@%s в опасности!", uc.UserName))
							}
						case "да", "Да":
							err := app.db.UpdateUserChatWarning(ctx, false, uc.ID)
							if err != nil {
								log.Println(err)
							}
						}
						delete(waitForAnswer, c.ChatID)
					}
				}
			}
		}
	}()

	go func() {
		lastNewsSent := 0
		for {
			news, err := app.db.GetLastNews(ctx, lastNewsSent)
			if err != nil {
				log.Printf(err.Error())
				continue
			}
			for _, n := range news {
				if n.ID > lastNewsSent {
					lastNewsSent = n.ID
				}

				cc, err := app.db.GetLastCoordinatesByHotLocationID(ctx, n.HotLocationID)
				if err != nil {

				}
				chats := make(map[int]struct{})
				for _, c := range cc {
					chats[c.UserChatID] = struct{}{}
				}

				for chat := range chats {
					uc, err := app.db.GetUserChat(ctx, chat)
					if err != nil {

					}
					err = app.bot.Send(uc.ChatID, n.Text)
					if err != nil {

					}
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Minute * 2)
			cc, err := app.db.GetUserChatsWarning(ctx)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, c := range cc {
				err := app.bot.Send(c.ChatID, "Все хорошо?")
				if err != nil {
					log.Println(err)
					continue
				}
				waitForAnswer[c.ChatID] = struct{}{}
			}
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/hot_location/create", api.AddHotLocation(app.db)).Methods(http.MethodPost)
	r.HandleFunc("/api/news/create", api.AddNews(app.db)).Methods(http.MethodPost)
	return http.ListenAndServe(":80", r)
}
