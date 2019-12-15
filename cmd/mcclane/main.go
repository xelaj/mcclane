package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"github.com/xelaj/mcclane/internal/app"
	"github.com/xelaj/mcclane/internal/db/pg"
	"github.com/xelaj/mcclane/internal/telegram"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	tgBot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("TG_TOKEN"), http.DefaultClient)
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DSN")
	time.Sleep(time.Second * 3) // waiting for db up
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	b := telegram.New(tgBot)
	db := pg.NewDB(conn)

	a := app.NewApp(b, db)
	log.Fatal(a.Work(context.Background()))
}
