package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"mcclane/internal/telegram"
)
const botTocken = "CLASSIFED"

func main() {

	tgBot, err := tgbotapi.NewBotAPI(botTocken)
	if err != nil {
		fmt.Println("error while connecting to telegram", err)

	}
	x := telegram.New(tgBot)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // TODO: extract
	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println("Error GetUpdatesChan")
	}
	exit := make(chan telegram.UserLocation)
	for update := range updates {
		loc, err := x.GetLocation(update)
		if err != nil {
			continue
		}
		exit <- *loc
		x.Listen(exit)
	}

}
