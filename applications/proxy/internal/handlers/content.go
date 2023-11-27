package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/proxy/pkg/entities"
	"github.com/h-varmazyar/kiwi/pkg/tgBotHelpers"
)

func (h *Handler) handleMedia(ctx context.Context, b *bot.Bot, update *models.Update) {
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
		tgBotHelpers.SendError(ctx, b, params)
		return
	}

	photo := update.Message.Photo[0]
	for _, p := range update.Message.Photo[1:] {
		if photo.FileSize < p.FileSize {
			photo = p
		}
	}

	post := &entities.Post{
		FileId: photo.FileID,
	}
	if err := h.postRepo.Create(ctx, post); err != nil {
		params := &tgBotHelpers.ErrParams{
			ChatId: chatId,
			Err:    err,
		}
		tgBotHelpers.SendError(ctx, b, params)
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   responseContentSaved,
	}); err != nil {
		params := &tgBotHelpers.ErrParams{
			ChatId:   chatId,
			IsSilent: true,
			Err:      err,
		}
		tgBotHelpers.SendError(ctx, b, params)
	}
}
