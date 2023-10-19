package handler

import (
	"fmt"
	"net/http"

	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/config"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot"
	"github.com/frasnym/go-furaphonify-telebot/service"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprintf(w, "Webhook OK")

	botService := service.NewBotService(config.GetConfig(), telebot.GetBot())

	// Get update
	update, err := botService.GetUpdate(r.Body)
	if err != nil {
		logger.Error(r.Context(), "unable to SetWebhook", err)
		return
	}

	// Handle message
	if update.Message != nil {
		replyTxt := fmt.Sprintf("reply of: %s", update.Message.Text)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "whatsappify":
				replyTxt = "Enter you phone number"
			default:
				replyTxt = fmt.Sprintf("invalid command: %s", update.Message.Command())
			}
		}

		_, err = botService.SendTextMessage(update.Message.Chat.ID, replyTxt)
		if err != nil {
			logger.Error(r.Context(), "unable to SendTextMessage", err)
			return
		}
	}
}
