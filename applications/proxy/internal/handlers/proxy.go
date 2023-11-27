package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/pkg/tgBotHelpers"
)

func (h *Handler) handleProxy(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		h.log.WithError(errUnsupportedMessage.AddDetail("msg", update))
		return
	}
	chatId := update.Message.Chat.ID
	proxyLinks := make([]string, 0)
	if k := update.Message.ReplyMarkup.InlineKeyboard; k != nil {
		for _, row := range k {
			for _, button := range row {
				proxyLinks = append(proxyLinks, button.URL)
			}
		}
	}

	if len(proxyLinks) > 0 {
		photo, err := h.postRepo.NewUnused(ctx)
		if err != nil {
			params := &tgBotHelpers.ErrParams{
				ChatId: chatId,
				Err:    err,
			}
			tgBotHelpers.SendError(ctx, b, params)
			return
		}
		photoFile := &models.InputFileString{
			Data: photo.FileId,
		}

		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      h.configs.PublishChannelId,
			Photo:       photoFile,
			Caption:     respProxyCaption(proxyLinks),
			ParseMode:   models.ParseModeMarkdown,
			ReplyMarkup: h.proxyKeyboard(b, proxyLinks),
		})
		if err != nil {
			params := &tgBotHelpers.ErrParams{
				ChatId: chatId,
				Err:    err,
			}
			tgBotHelpers.SendError(ctx, b, params)
		}
	} else {
		params := &tgBotHelpers.ErrParams{
			ChatId: chatId,
			Err:    errInvalidProxy,
		}
		tgBotHelpers.SendError(ctx, b, params)
	}
}
