package tgBotHelpers

import (
	"context"
	"github.com/go-telegram/bot"
)

const defaultSuccessMsg = `عملیات با موفقیت انجام شد`

type SuccessParams struct {
	ChatId           int64
	Msg              string
	IsSilent         bool
	ReplyToMessageId int
}

func SendSuccess(ctx context.Context, b *bot.Bot, params *SuccessParams) {
	if params.Msg == "" {
		params.Msg = defaultSuccessMsg
	}

	if !params.IsSilent {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           params.ChatId,
			Text:             params.Msg,
			ReplyToMessageID: params.ReplyToMessageId,
		})
	}
}
