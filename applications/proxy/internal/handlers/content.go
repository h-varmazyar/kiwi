package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/proxy/pkg/entities"
	"github.com/h-varmazyar/kiwi/pkg/tgBotHelpers"
)

func (h *Handler) handleMedia(ctx context.Context, bot *bot.Bot, update *models.Update) {
	if update.Message == nil {
		h.log.WithError(errUnsupportedMessage.AddDetail("msg", update))
		return
	}
	chatId := update.Message.Chat.ID
	if update.Message.Photo == nil || len(update.Message.Photo) == 0 {
		params := &tgBotHelpers.ErrParams{
			ChatId: chatId,
			Err:    errInvalidContent,
		}
		tgBotHelpers.SendError(ctx, bot, params)
		return
	}
	for _, photo := range update.Message.Photo {
		post := &entities.Post{
			FileId: photo.FileID,
		}
		if err := h.postRepo.Create(ctx, post); err != nil {
			params := &tgBotHelpers.ErrParams{
				ChatId: chatId,
				Err:    err,
			}
			tgBotHelpers.SendError(ctx, bot, params)
			return
		}
	}
}
