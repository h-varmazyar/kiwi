package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
)

func (h *Handler) stateAddMedia(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {

		return
	}

	media, err := h.addContentRepo.GetMedia(ctx, update.Message.Chat.ID)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Message.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	if update.Message.Video != nil {
		media.TelegramFileId = update.Message.Video.FileID
	} else {
		//if url, err:=url.Parse(update.Message.Text);err!=nil{
		//	helpers.SendError(ctx, b, &helpers.ErrParams{
		//		ChatId:   update.Chat.ID,
		//		Err:      err,
		//		Metadata: string(data),
		//	})
		//	return
		//}
		//}else{
		//	send invalid err
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        " حال حاضر این امکان پشتیبانی نمیشود",
			ReplyMarkup: h.keyboardCancel(b),
		})
	}

	if err = h.addContentRepo.SetMedia(ctx, update.Message.Chat.ID, media); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Message.Chat.ID,
			Err:      err,
			Metadata: media,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        MsgSetVideoQuality,
		ReplyMarkup: h.keyboardMediaQualities(ctx, b),
	})
}

func (h *Handler) setQuality(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	media, err := h.addContentRepo.GetMedia(ctx, update.Chat.ID)
	if err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	media.Quality = entities.MediaQuality(data)

	if err = h.mediaRepo.Create(ctx, media); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Chat.ID,
		Text:        MsgSetOtherVideoQuality,
		ReplyMarkup: h.keyboardContinueQuality(ctx, b),
	})
}

func (h *Handler) completeAddQualities(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
	if err := h.userStateRepo.DeleteState(ctx, update.Chat.ID); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}
	if err := h.addContentRepo.DeleteMedia(ctx, update.Chat.ID); err != nil {
		SendError(ctx, b, &ErrParams{
			ChatId:   update.Chat.ID,
			Err:      err,
			Metadata: update,
		})
		return
	}

	SendSuccess(ctx, b, &SuccessParams{
		ChatId:   update.Chat.ID,
		IsSilent: false,
	})
}
