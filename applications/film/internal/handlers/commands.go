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

	if strings.HasPrefix(cmdValue, "ep") {
		videoId := strings.TrimPrefix(cmdValue, "ep")

		id, err := strconv.Atoi(videoId)
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

		episode, err := h.episodeRepo.Return(ctx, media.OwnerID)
		if err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
			})
			return
		}

		if err = h.episodeRepo.Visit(ctx, episode.SeasonId); err != nil {
			SendError(ctx, b, &ErrParams{
				ChatId:   update.Message.Chat.ID,
				Err:      err,
				Metadata: media,
				IsSilent: true,
			})
		}

		video := &models.InputFileString{
			Data: media.TelegramFileId,
		}
		_, _ = b.SendVideo(ctx, &bot.SendVideoParams{
			ChatID:         update.Message.Chat.ID,
			Video:          video,
			Caption:        episode.Title,
			ProtectContent: true,
		})
	}
}
