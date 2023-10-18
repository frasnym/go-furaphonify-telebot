package telebot

import (
	"errors"

	"github.com/frasnym/go-furaphonify-telebot/model"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot/tgbotapi"
)

var bot *model.TeleBot

func InitTgBotApi() {
	tgBot := tgbotapi.NewTgBotApi()
	bot = &tgBot
}

func GetBot() model.TeleBot {
	if bot == nil {
		panic(errors.New("please init bot first"))
	}

	return *bot
}
