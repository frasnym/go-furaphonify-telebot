package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/frasnym/go-furaphonify-telebot/common"
	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/pkg/session"
	"github.com/frasnym/go-furaphonify-telebot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// WhatsappService is an interface for managing WhatsApp-related actions.
type WhatsappService interface {
	Request(ctx context.Context, userID int, chatID int64) error
	Processor(ctx context.Context, userID int, input string) error
}

type whatsappSvc struct {
	botRepo repository.BotRepository
}

// NewWhatsappService creates a new WhatsappService using the provided bot repository.
func NewWhatsappService(botRepo *repository.BotRepository) WhatsappService {
	return &whatsappSvc{botRepo: *botRepo}
}

// Request initiates a request for a user to enter their phone number for WhatsAppification.
func (s *whatsappSvc) Request(ctx context.Context, userID int, chatID int64) error {
	var err error
	defer func() {
		logger.LogService(ctx, "WhatsappRequest", err)
	}()

	// Start a new session for the user
	session.NewSession(userID, chatID, common.CommandWhatsappify)

	// Send a request for the phone number
	replyTxt := "Enter your phone number"
	msg, err := s.botRepo.SendTextMessage(ctx, chatID, replyTxt)
	if err != nil {
		err = fmt.Errorf("error sending text message: %w", err)
		return err
	}

	// Set the message ID in the user's session
	if err := session.SetMessageID(userID, msg.MessageID); err != nil {
		err = fmt.Errorf("error setting message ID: %w", err)
		return err
	}
	return nil
}

// Processor processes the user's input (phone number) for WhatsAppification.
func (s *whatsappSvc) Processor(ctx context.Context, userID int, input string) error {
	var err error
	defer func() {
		logger.LogService(ctx, "WhatsappProcessor", err)
	}()

	chatID, err := session.GetChatID(userID)
	if err != nil {
		err = fmt.Errorf("error getting chat ID: %w", err)
		return err
	}

	messageID, err := session.GetMessageID(userID)
	if err != nil {
		err = fmt.Errorf("error getting message ID: %w", err)
		return err
	}

	replyText := "error"
	var replyMarkup *tgbotapi.InlineKeyboardMarkup
	defer func() {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, replyText)
		if replyMarkup != nil {
			editMsg.ReplyMarkup = replyMarkup
		}

		_, err = s.botRepo.SendMessage(ctx, editMsg)
		if err != nil {
			err = fmt.Errorf("error sending message: %w", err)
		}

		session.DeleteUserSession(userID)
	}()

	// Process phone number
	replyText = fmt.Sprintf("Processing: %s\n\n", input)
	phoneNumber := common.RemoveNonAlphanumeric(input)
	if len(phoneNumber) <= 0 {
		errText := "Unable to parse your number, please try again."
		if err = session.ResetTimer(userID); err != nil {
			err = fmt.Errorf("error resetting timer: %w", err)
			errText = "Timeout"
		}
		replyText = fmt.Sprint(replyText, errText)

		return err
	}

	// Normalize the phone number to the WhatsApp format
	pattern := regexp.MustCompile(`^0+`)
	phoneNumber = pattern.ReplaceAllString(phoneNumber, "62")
	whatsappUrl := fmt.Sprintf("wa.me/%s", phoneNumber)
	replyText = fmt.Sprint(replyText, "Success")

	// Create a result with an inline keyboard for the WhatsApp link
	urlButton := tgbotapi.NewInlineKeyboardButtonURL("WhatsApp 🌎", whatsappUrl)
	inlineKeyboard := tgbotapi.NewInlineKeyboardRow(urlButton)
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(inlineKeyboard)
	replyMarkup = &inlineKeyboardMarkup

	err = nil
	return nil
}
