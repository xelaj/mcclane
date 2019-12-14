package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"mcclane/internal/telegram"
)

const botTocken = 

type Dialog struct {
	Choice     string
	Answer     string
	NextAnswer []Dialog
	RegUser    func(int, string) string
	DelUser    func(int)
}

var DLG []Dialog
var NextDLG = make(map[int]*[]Dialog)

func main() {
	DLG = loadContent()

	tgBot, err := tgbotapi.NewBotAPI(botTocken)
	if err != nil {
		fmt.Println("error while connecting to telegram", err)

	}
	x := telegram.New(tgBot)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println("Error GetUpdatesChan")
	}
	exit := make(chan telegram.UserLocation)
	for update := range updates {
		loc, err := x.GetLocation(update)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, Answer(update.Message.From.ID, update.Message.Text))
			msg.ParseMode = tgbotapi.ModeMarkdown
			_, _ = tgBot.Send(msg)

			continue
		}
		exit <- *loc
		x.Listen(exit)
	}
}
func Answer(userID int, incomeMsg string) string {
	switch incomeMsg {
	case "/start":
		NextDLG[userID] = &DLG
		return `Добро пожаловать в протестный бот.
	Отправьте свою локацию (трансляцию на 8-12 часов) в месте скопления граждан и переодически проверяйте поступающую от нас информацию.
	Вы также можете "зарегистрироваться" (команда /reg) и указать несколько доверенных контактов, которым будет сообщено в случае обнаружения вас вне "горячей" зоны, что подразумевает задержание и транспортировку в ОВД`
	default:

		AnswerThree, ok := NextDLG[userID]
		if !ok {
			AnswerThree = &DLG
		}

		for _, currentAnswer := range *AnswerThree {
			if currentAnswer.RegUser != nil {
				fmt.Println(RegUser(userID, incomeMsg))
				if len(currentAnswer.NextAnswer) > 0 {
					NextDLG[userID] = &currentAnswer.NextAnswer
					return currentAnswer.Answer
				}
				NextDLG[userID] = &DLG
				return currentAnswer.Answer
			}

			if currentAnswer.Choice == incomeMsg {
				if len(currentAnswer.NextAnswer) > 0 {
					NextDLG[userID] = &currentAnswer.NextAnswer
				} else {
					NextDLG[userID] = &DLG
				}
				if currentAnswer.DelUser != nil {
					DelUser(userID)
				}

				return currentAnswer.Answer
			}
		}

	}
	return "incorrect comand"
}

func loadContent() []Dialog {

	DLG = []Dialog{
		{
			Choice: "/reg",
			Answer: "Отправьте свое ФИО и дд.мм.гггг рождения, это информация, которую требуют в ОВД",
			NextAnswer: []Dialog{{
				RegUser: RegUser,
				Answer:  "Теперь отправьте телеграм логины лиц, которых следует оповестить при вашем задержании. Пример @durov, @navalny, @msvetov",
				NextAnswer: []Dialog{
					{
						RegUser: RegUser,
						Answer:  "Принято. Учтите, довереные лица должны хоть раз написать боту любое сообщение и не удлять диалог с ним",
					}}}},
		},
		{
			Choice: "/del",
			Answer: "Подтвердите желание удалить информацию о себе. Отправьте сумму чисел 6+12 ",
			NextAnswer: []Dialog{
				{
					DelUser: DelUser,
					Choice:  "18",
				}},
		},
	}
	return DLG
}

func RegUser(id int, data string) string {

	return "Reg call"
}
func DelUser(id int) {
	fmt.Println("DelUser called by", id)
}
