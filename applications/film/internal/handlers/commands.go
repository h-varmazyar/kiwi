package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
)

func (h *Handler) startCmd(ctx context.Context, b *bot.Bot, update *models.Update) {
	cmdValue := strings.Trim(update.Message.Text, "/start ")

	id, err := strconv.Atoi(cmdValue)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Message.Chat.ID,
			Err:      err,
			Metadata: update.Message.Text,
		})
		return
	}
	media, err := h.mediaRepo.Return(ctx, uint(id))
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Message.Chat.ID,
			Err:      err,
			Metadata: media,
		})
		return
	}

	var video *models.InputFileString
	caption := ""
	if media.OwnerType == "series" {
		episode, err := h.episodeRepo.Return(ctx, media.OwnerID)
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
			})
			return
		}

		caption = prepareEpisodeCaptionForWatch(episode, media.Quality)
		video = &models.InputFileString{
			Data: media.TelegramFileId,
		}
		if err = h.episodeRepo.Visit(ctx, episode.SeasonId); err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
				IsSilent: true,
			})
		}

	} else if media.OwnerType == "movie" {
		movie, err := h.moviesRepo.Return(ctx, media.OwnerID)
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
			})
			return
		}

		caption = prepareMovieQualityCaption(movie, media)
		video = &models.InputFileString{
			Data: media.TelegramFileId,
		}
		if err = h.moviesRepo.Visit(ctx, movie.ID); err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
				IsSilent: true,
			})
		}
	}

	if video != nil {
		_, _ = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:         update.Message.Chat.ID,
			Video:          video,
			Caption:        caption,
			ProtectContent: true,
		})
	}
}

func (h *Handler) addCmd(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userId := update.Message.Chat.ID
	if !h.isAdmin(userId) {
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        MsgAddContent,
		ReplyMarkup: h.keyboardAdd(b),
	})
}
