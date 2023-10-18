package handler

import (
	"fmt"
	"net/http"

	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprintf(w, "<b>Webhook OK</b>")

	telebot.InitTgBotApi()
	bot := telebot.GetBot()

	update, err := bot.GetUpdate(r.Body)
	if err != nil {
		logger.Error(r.Context(), "unable to SetWebhook", err)
		return
	}

	if update.Message != nil {
		_, err = bot.SendTextMessage(update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			logger.Error(r.Context(), "unable to SendTextMessage", err)
			return
		}
	}

}
