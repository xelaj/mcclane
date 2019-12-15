package app

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/xelaj/mcclane/internal/api"
	"github.com/xelaj/mcclane/internal/db/pg"
	"github.com/xelaj/mcclane/internal/model"
	"github.com/xelaj/mcclane/internal/telegram"
	"log"
	"net/http"
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

				_, err = app.db.AddCoordinates(ctx, model.Coordinates{
					CreatedAt:     time.Now(),
					UserChatID:    uc.ID,
					Latitude:      loc.Location.Latitude,
					Longitude:     loc.Location.Longitude,
					HotLocationID: hl.ID,
				})
				if err != nil {
					log.Printf("")
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

	r := mux.NewRouter()
	r.HandleFunc("/api/hot_location/create", api.AddHotLocation(app.db)).Methods(http.MethodPost)
	r.HandleFunc("/api/news/create", api.AddNews(app.db)).Methods(http.MethodPost)
	return http.ListenAndServe(":80", r)
}
