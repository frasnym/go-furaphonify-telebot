package service

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/frasnym/go-furaphonify-telebot/common"
	"github.com/frasnym/go-furaphonify-telebot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotService interface {
	SetWebhook() error
	GetUpdate(r io.Reader) (*tgbotapi.Update, error)
	SendTextMessage(chatID int64, text string) (*tgbotapi.Message, error)
}

type botSvc struct {
	bot *tgbotapi.BotAPI
	cfg *config.Config
}

// GetUpdate implements BotService.
func (*botSvc) GetUpdate(r io.Reader) (*tgbotapi.Update, error) {
	update := tgbotapi.Update{}
	if err := json.NewDecoder(r).Decode(&update); err != nil {
		return nil, fmt.Errorf("error decoding update: %w", err)
	}

	return &update, nil
}

// SendTextMessage implements BotService.
func (s *botSvc) SendTextMessage(chatID int64, text string) (*tgbotapi.Message, error) {
	stringMsg := tgbotapi.NewMessage(chatID, text)
	msg, err := s.bot.Send(stringMsg)
	if err != nil {
		return nil, fmt.Errorf("unable to SendTextMessage: %w", err)
	}

	return &msg, nil
}

// SetWebhook implements BotService.
func (s *botSvc) SetWebhook() error {
	webhookURL := fmt.Sprintf("https://%s/webhook", s.cfg.VercelUrl)

	info, err := s.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("unable to GetWebhookInfo: %w", err)
	}
	if info.URL == webhookURL {
		return common.ErrNoChanges
	}

	_, err = s.bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		return fmt.Errorf("unable to SetWebhook: %w", err)
	}

	return nil
}

func NewBotService(cfg *config.Config, bot *tgbotapi.BotAPI) BotService {
	return &botSvc{cfg: cfg, bot: bot}
}
