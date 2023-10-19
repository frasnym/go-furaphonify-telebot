package handler

import (
	"fmt"
	"net/http"

	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/config"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot"
	"github.com/frasnym/go-furaphonify-telebot/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer fmt.Fprintf(w, "Index OK")

	botService := service.NewBotService(config.GetConfig(), telebot.GetBot())
	err := botService.SetWebhook()
	if err != nil {
		logger.Error(r.Context(), "unable to SetWebhook", err)
	}
}
