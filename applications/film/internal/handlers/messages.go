package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) stateCheck(ctx context.Context, b *bot.Bot, update *models.Update) {
	//userId := update.Message.Chat.ID
	//state, err := h.userStateRepository.GetState(ctx, userId)
	//if err != nil {
	//	params := &helpers.ErrParams{
	//		ChatId:           update.Message.Chat.ID,
	//		Msg:              "",
	//		IsSilent:         false,
	//		Err:              err,
	//		ReplyToMessageId: update.Message.ID,
	//		Metadata:         update,
	//	}
	//	helpers.SendError(ctx, b, params)
	//}
}
