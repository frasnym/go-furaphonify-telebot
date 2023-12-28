package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/frasnym/go-furaphonify-telebot/common"
	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/pkg/session"
	"github.com/frasnym/go-furaphonify-telebot/pkg/truecaller"
	"github.com/frasnym/go-furaphonify-telebot/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// WhatsappService is an interface for managing WhatsApp-related actions.
type WhatsappService interface {
	Request(ctx context.Context, userID int, chatID int64) error
	Processor(ctx context.Context, userID int, input string) error
}

type whatsappSvc struct {
	botRepo    repository.BotRepository
	gsheetRepo repository.GSheetRepository
}

// Request initiates a request for a user to enter their phone number for WhatsAppification.
func (s *whatsappSvc) Request(ctx context.Context, userID int, chatID int64) error {
	var err error
	now := time.Now()
	defer func() {
		logger.LogService(ctx, "WhatsappRequest", err, &now)
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
	now := time.Now()
	defer func() {
		logger.LogService(ctx, "WhatsappProcessor", err, &now)
	}()

	if session.IsInteractionTimedOut(userID) {
		err = s.notifyError(ctx, userID, "Request Timeout")
		if err != nil {
			err = fmt.Errorf("err notifyError: %w", err)
		}

		session.DeleteUserSession(userID)
		return err
	}

	chatID, err := session.GetChatID(userID)
	if err != nil {
		err = fmt.Errorf("err session.GetChatID: %w", err)
		return err
	}

	messageID, err := session.GetMessageID(userID)
	if err != nil {
		err = fmt.Errorf("err session.GetMessageID: %w", err)
		return err
	}

	// Process phone number
	replyText := fmt.Sprintf("Processing: %s\n\n", input)
	phoneNumber := common.RemoveNonNumeric(input)
	if len(phoneNumber) <= 0 {
		if err = session.ResetTimer(userID); err != nil {
			err = s.notifyError(ctx, userID, "No Session")
			if err != nil {
				err = fmt.Errorf("err notifyError: %w", err)
			}
			return nil
		}
		replyText = fmt.Sprint(replyText, "Unable to parse your number, please try again.")
		err = s.notifyError(ctx, userID, replyText)
		if err != nil {
			err = fmt.Errorf("err notifyError: %w", err)
		}

		return err
	}

	var replyMarkup *tgbotapi.InlineKeyboardMarkup
	defer func() {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, replyText)
		if replyMarkup != nil {
			editMsg.ReplyMarkup = replyMarkup
		}

		_, err = s.botRepo.SendMessage(ctx, editMsg)
		if err != nil {
			err = fmt.Errorf("err botRepo.SendMessage: %w", err)
		}

		session.DeleteUserSession(userID)
	}()

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

	// Phone Information
	// TODO: Get from sheet db first
	tcResp, err := s.getPhoneInformation(ctx, phoneNumber)
	if err != nil {
		err = fmt.Errorf("err notifyError: %w", err)
		return err
	}

	// Update message
	replyText = fmt.Sprint(replyText, "\n", tcResp.ParseInformationMessage())

	// Save to db
	if err = s.gsheetRepo.AppendRow(ctx, phoneNumber, tcResp.Name, tcResp.Raw); err != nil {
		err = fmt.Errorf("err gsheetRepo.AppendRow: %w", err)
		return err
	}

	err = nil
	return nil
}

// NewWhatsappService creates a new WhatsappService using the provided bot repository.
func NewWhatsappService(botRepo *repository.BotRepository, gsheetRepo *repository.GSheetRepository) WhatsappService {
	return &whatsappSvc{botRepo: *botRepo, gsheetRepo: *gsheetRepo}
}

func (s *whatsappSvc) notifyError(ctx context.Context, userID int, msg string) error {
	var err error
	now := time.Now()
	defer func() {
		logger.LogService(ctx, "WhatsappNotifyError", err, &now)
	}()

	chatID, err := session.GetChatID(userID)
	if err != nil {
		err = fmt.Errorf("err session.GetChatID: %w", err)
		return err
	}

	messageID, err := session.GetMessageID(userID)
	if err != nil {
		err = fmt.Errorf("err session.GetMessageID: %w", err)
		return err
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, msg)
	_, err = s.botRepo.SendMessage(ctx, editMsg)
	if err != nil {
		err = fmt.Errorf("err botRepo.SendMessage: %w", err)
		return err
	}

	return nil
}

func (s *whatsappSvc) getPhoneInformation(ctx context.Context, phoneNumber string) (*truecaller.SearchResponse, error) {
	var err error
	now := time.Now()
	defer func() {
		logger.LogService(ctx, "WhatsappGetPhoneInformation", err, &now)
	}()

	tcResp, err := truecaller.GetPhoneNumberInformation(ctx, phoneNumber)
	if err != nil {
		err = fmt.Errorf("err truecaller.GetPhoneNumberInformation: %w", err)
		return nil, err
	}

	return tcResp, nil
}
