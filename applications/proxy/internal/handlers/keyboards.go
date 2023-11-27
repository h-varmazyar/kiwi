package handlers

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) proxyKeyboard(_ *bot.Bot, proxyLinks []string) models.InlineKeyboardMarkup {
	kb := models.InlineKeyboardMarkup{
		InlineKeyboard: make([][]models.InlineKeyboardButton, 0),
	}
	row := make([]models.InlineKeyboardButton, 0)
	for i, link := range proxyLinks {
		btn := models.InlineKeyboardButton{
			Text: "اتصال ✅",
			URL:  link,
		}
		row = append(row, btn)
		if i%2 == 1 || i == len(proxyLinks)-1 {
			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
			row = make([]models.InlineKeyboardButton, 0)
		}
	}

	return kb
}
