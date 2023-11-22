package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
)

func (h *Handler) selectSeries(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	seriesId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	series, err := h.seriesRepo.Return(ctx, uint(seriesId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	genres := ""
	for _, genre := range series.Genres {
		genres = fmt.Sprintf("%v | %v", genres, genre.FaName)
	}

	genres = strings.Trim(strings.TrimSpace(genres), "|")

	text := fmt.Sprintf(MsgSeriesInfo, series.Title, series.FaName, series.EnName, series.Presentation, series.Year, series.IMDB, genres)

	photo := &models.InputFileString{
		Data: series.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Photo:       photo,
		Caption:     text,
		ReplyMarkup: h.keyboardSeriesInfo(ctx, b, series.ID),
	})
}
