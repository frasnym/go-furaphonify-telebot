package tgbotapi

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/frasnym/go-furaphonify-telebot/common"
	"github.com/frasnym/go-furaphonify-telebot/config"
	"github.com/frasnym/go-furaphonify-telebot/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewTgBotApi() model.TeleBot {
	cfg := config.GetConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		panic(fmt.Errorf("unable to init tgbotapi: %w", err))
	}

	return &tgBotApi{
		cfg: cfg,
		bot: bot,
	}
}

type tgBotApi struct {
	cfg *config.Config
	bot *tgbotapi.BotAPI
}

// GetUpdate implements model.TeleBot.
func (*tgBotApi) GetUpdate(r io.Reader) (*model.TeleBotUpdate, error) {
	update := tgbotapi.Update{}
	if err := json.NewDecoder(r).Decode(&update); err != nil {
		return nil, fmt.Errorf("error decoding update: %w", err)
	}

	return &model.TeleBotUpdate{
		UpdateID: update.UpdateID,
		Message: &model.TeleBotMessage{
			MessageID: update.Message.MessageID,
			Chat: &model.TeleBotChat{
				ID:   update.Message.Chat.ID,
				Type: update.Message.Chat.Type,
			},
			Text: update.Message.Text,
		},
	}, nil
}

// SetWebhook Set up the webhook for receiving updates from Telegram.
func (b *tgBotApi) SetWebhook() error {
	webhookURL := fmt.Sprintf("https://%s/webhook", b.cfg.VercelUrl)

	info, err := b.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("unable to GetWebhookInfo: %w", err)
	}
	if info.URL == webhookURL {
		return common.ErrNoChanges
	}

	_, err = b.bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		return fmt.Errorf("unable to SetWebhook: %w", err)
	}

	return nil
}

// SendTextMessage implements telebot.Bot.
func (b *tgBotApi) SendTextMessage(chatID int64, text string) (*model.TeleBotMessage, error) {
	stringMsg := tgbotapi.NewMessage(chatID, text)
	msg, err := b.bot.Send(stringMsg)
	if err != nil {
		return nil, fmt.Errorf("unable to SendTextMessage: %w", err)
	}

	return &model.TeleBotMessage{
		MessageID: msg.MessageID,
	}, nil
}
