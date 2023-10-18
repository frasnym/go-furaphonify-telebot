package handler

import (
	"fmt"
	"net/http"

	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	telebot.InitTgBotApi()
	bot := telebot.GetBot()

	err := bot.SetWebhook()
	if err != nil {
		logger.Error(r.Context(), "unable to SetWebhook", err)
	}

	fmt.Fprintf(w, "Index OK")
}
