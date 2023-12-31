package tgBotHelpers

import (
	"context"
	"encoding/json"
	e "errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/h-varmazyar/kiwi/pkg/errors"
)

const (
	defaultErrMsg = `خطای پیشبینی نشده ای رخ داده است`
)

type ErrParams struct {
	ChatId           int64
	Msg              string
	IsSilent         bool
	Err              error
	ReplyToMessageId int
	Metadata         interface{}
	ReplyMarkup      *models.ReplyMarkup
}

func SendError(ctx context.Context, b *bot.Bot, params *ErrParams) {
	er := errors.Cast(params.Err)
	if params.Metadata != nil {
		s, _ := json.Marshal(params.Metadata)
		er = er.AddDetail("update", s)
	}

	if params.Msg == "" {
		t := errors.New("")
		if e.As(params.Err, &t) {
			params.Msg = t.AddLang(ctx.Value("lang").(string)).Error()
		} else {
			params.Msg = defaultErrMsg
		}
	}

	fmt.Println("error:", er)

	if !params.IsSilent {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:           params.ChatId,
			Text:             params.Msg,
			ReplyToMessageID: params.ReplyToMessageId,
		})
		if err != nil {
			fmt.Println("failed to send error:", err)
		}
	}
}
