package handler

import (
	"fmt"
	"net/http"

	"github.com/frasnym/go-furaphonify-telebot/common"
	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/config"
	"github.com/frasnym/go-furaphonify-telebot/pkg/session"
	"github.com/frasnym/go-furaphonify-telebot/pkg/telebot"
	"github.com/frasnym/go-furaphonify-telebot/repository"
	"github.com/frasnym/go-furaphonify-telebot/service"
)

// WebhookHandler handles incoming HTTP requests for a Telegram bot's webhook.
// It processes updates, handles commands, and takes appropriate actions based on the received messages.
// If any errors occur during the process, they are logged.
// After processing the request, it writes a "Webhook OK" message to the response writer (w).
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	// Log any errors and write "Webhook OK" as the API response
	defer func() {
		logger.LogService(ctx, "WebhookHandler", err)
		fmt.Fprintf(w, "Webhook OK")
	}()

	// Create a new bot repository with the application's configuration and Telegram bot
	botRepo := repository.NewBotRepository(config.GetConfig(), telebot.GetBot())
	whatsappSvc := service.NewWhatsappService(&botRepo)

	// Get the update from the request body
	update, err := botRepo.GetUpdate(ctx, r.Body)
	if err != nil {
		err = fmt.Errorf("err botRepo.GetUpdate: %w", err)
		return
	}

	// Handle messages and commands
	if update.Message != nil {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		// Handle commands
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case common.CommandWhatsappify:
				if err = whatsappSvc.Request(ctx, userID, chatID); err != nil {
					err = fmt.Errorf("err whatsappSvc.Request: %w", err)
				}
				return
			default:
				err = fmt.Errorf("invalid command: %s", update.Message.Command())
				return
			}
		}

		// Get the user's current action
		action, errSession := session.GetAction(userID)
		if errSession != nil {
			err = fmt.Errorf("err session.GetAction: %w", errSession)
			return
		}

		// Handle requests based on the user's current action
		switch action {
		case common.CommandWhatsappify:
			if err = whatsappSvc.Processor(ctx, userID, update.Message.Text); err != nil {
				err = fmt.Errorf("err whatsappSvc.Processor: %w", err)
			}
			return
		}

		err = fmt.Errorf("unprocessable text: %s", update.Message.Text)
	}
}
