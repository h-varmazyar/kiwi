package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func (h *Handler) showSeasons(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
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

	text := fmt.Sprintf(MsgSeriesInfo, series.Title, series.FaName, series.EnName, series.Presentation, series.Year, series.IMDB, genres)

	seasons, err := h.seasonRepo.SeriesSeasons(ctx, uint(seriesId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	photo := &models.InputFileString{
		Data: series.Banner.TelegramFileId,
	}

	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Caption:     text,
		Photo:       photo,
		ReplyMarkup: h.keyboardSeasonList(ctx, b, seasons),
	})
}

func (h *Handler) showSeason(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	chatId := update.Chat.ID
	seasonId, err := strconv.Atoi(string(data))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	season, err := h.seasonRepo.Return(ctx, uint(seasonId))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   chatId,
			Err:      err,
			Metadata: update,
		})
		return
	}

	text := fmt.Sprintf(MsgSeasonInfo, season.Title, season.Presentation, season.Year, season.IMDB)

	photo := &models.InputFileString{
		Data: season.Banner.TelegramFileId,
	}
	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:      update.Chat.ID,
		Caption:     text,
		Photo:       photo,
		ReplyMarkup: h.keyboardSeasonInfo(ctx, b, season.ID),
	})
}

func (h *Handler) addSeason(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	//chatId := update.Chat.ID
	//episodeId, err := strconv.Atoi(string(data))
	//if err != nil {
	//	helpers.SendError(ctx, b, &helpers.ErrParams{
	//		ChatId:   chatId,
	//		Err:      err,
	//		Metadata: update,
	//	})
	//	return
	//}
	//
	//episode, err := h.episodesRepo.Return(ctx, uint(episodeId))
	//if err != nil {
	//	helpers.SendError(ctx, b, &helpers.ErrParams{
	//		ChatId:   chatId,
	//		Err:      err,
	//		Metadata: update,
	//	})
	//	return
	//}
	//
	//text := fmt.Sprintf(MsgSeriesInfo, series.Title, series.FaName, series.EnName, series.Presentation, series.Year, series.IMDB, genres)
	//
	//_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
	//	ChatID:      update.Chat.ID,
	//	Text:        text,
	//	ReplyMarkup: keyboardSeriesInfo,
	//})
}
