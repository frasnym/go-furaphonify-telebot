package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/frasnym/go-furaphonify-telebot/common/logger"
	"github.com/frasnym/go-furaphonify-telebot/config"
	"github.com/google/uuid"
	"google.golang.org/api/sheets/v4"
)

type GSheetRepository interface {
	AppendRow(ctx context.Context, phone, name, metadata string) error
}

type gsheetRepo struct {
	cfg     *config.Config
	service *sheets.Service
}

// AppendRow implements GSheetRepository.
func (repo *gsheetRepo) AppendRow(ctx context.Context, phone, name, metadata string) error {
	var err error
	now := time.Now()
	defer func() {
		logger.LogRepository(ctx, "GSheetDeleteMessage", err, &now)
	}()

	values := &sheets.ValueRange{
		Values: [][]interface{}{{
			uuid.NewString(),
			phone,
			name,
			metadata,
		}},
	}

	_, err = repo.service.Spreadsheets.Values.
		Append(repo.cfg.GsheetID, "db!A:C", values).ValueInputOption("RAW").Do()

	if err != nil {
		err = fmt.Errorf("err repo.service.Spreadsheets.Values.Append: %w", err)
		return err
	}

	return nil
}

func NewGSheetRepository(cfg *config.Config, service *sheets.Service) GSheetRepository {
	return &gsheetRepo{cfg: cfg, service: service}
}
