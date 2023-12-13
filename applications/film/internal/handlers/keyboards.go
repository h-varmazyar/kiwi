package handlers

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/ui/keyboard/inline"
	entities2 "github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
)

func (h *Handler) keyboardAdd(b *bot.Bot) *inline.Keyboard {
	return inline.New(b).
		Row().
		Button("سریال", []byte("سریال"), h.addEpisode).
		//Button("سریال", []byte("سریال"), h.addSeries).
		Button("قسمت سریال", []byte("قسمت سریال"), h.addEpisode).
		Row().
		Button("سینمایی", []byte("سینمایی"), h.addMovie).
		Button("افزودن کیفیت", []byte("افزودن کیفیت"), h.addEpisode).
		//Button("افزودن کیفیت", []byte("افزودن کیفیت"), h.addQuality).
		Row().
		Button("لغو", []byte{}, h.cancelAddition)
}

func (h *Handler) keyboardCancel(b *bot.Bot) *inline.Keyboard {
	return inline.New(b).
		Row().
		Button("لغو", []byte{}, h.cancelAddition)
}

func (h *Handler) keyboardSearchedSeries(_ context.Context, b *bot.Bot, searchResults []*SearchResult) *inline.Keyboard {
	keyboard := inline.New(b)

	for _, result := range searchResults {
		row := keyboard.Row()
		if result.Type == SearchResultMovie {
			row.Button(fmt.Sprintf("%v | %v", result.FaName, result.EnName), []byte(fmt.Sprintf("%v", result.Id)), h.selectMovie)
		} else {
			row.Button(fmt.Sprintf("%v | %v", result.FaName, result.EnName), []byte(fmt.Sprintf("%v", result.Id)), h.selectSeries)
		}
	}

	keyboard.Row().Button("لغو", []byte(""), h.cancelAddition)

	return keyboard
}

func (h *Handler) keyboardSeriesInfo(_ context.Context, b *bot.Bot, seriesId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("فصل‌ها", []byte(fmt.Sprintf("%v", seriesId)), h.showSeasons).
		Button("افزودن فصل جدید", []byte(fmt.Sprintf("%v", seriesId)), h.addSeason)

	return keyboard
}

func (h *Handler) keyboardMovieInfo(_ context.Context, b *bot.Bot, isAdmin bool, movieId uint) *inline.Keyboard {
	keyboard := inline.New(b)
	if isAdmin {
		keyboard.Row().
			Button("افزودن کیفیت جدید", []byte(fmt.Sprintf("%v", movieId)), h.addMovieVideo).
			Button("ارسال در کانال اصلی", []byte(fmt.Sprintf("%v", movieId)), h.sendMovieToPublicChannel)
	}
	keyboard.Row().Button("مشاهده", []byte(fmt.Sprintf("%v", movieId)), h.showMovie)

	return keyboard
}

func (h *Handler) keyboardEpisodeAction(_ context.Context, b *bot.Bot, episodeId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("افزودن کیفیت جدید", []byte(fmt.Sprintf("%v", episodeId)), h.addEpisodeVideo).
		Row().
		Button("ارسال در کانال خصوصی", []byte(fmt.Sprintf("%v", episodeId)), h.sendEpisodeToPrivateChannel).
		Button("ارسال در کانال اصلی", []byte(fmt.Sprintf("%v", episodeId)), h.sendEpisodeToPublicChannel)

	return keyboard
}

func (h *Handler) keyboardEpisodeSubmit(_ context.Context, b *bot.Bot, episodeId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("تایید", []byte(fmt.Sprintf("%v", episodeId)), h.submitEpisodeAddition).
		Button("لغو", []byte(fmt.Sprintf("%v", episodeId)), h.cancelEpisodeAddition)

	return keyboard
}

func (h *Handler) keyboardMovieSubmit(_ context.Context, b *bot.Bot, movieId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("تایید", []byte(fmt.Sprintf("%v", movieId)), h.submitMovieAddition).
		Button("لغو", []byte(fmt.Sprintf("%v", movieId)), h.cancelMovieAddition)

	return keyboard
}

func (h *Handler) keyboardAddEpisodeVideo(_ context.Context, b *bot.Bot, episodeId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("افزودن کیفیت جدید", []byte(fmt.Sprintf("%v", episodeId)), h.addEpisodeVideo)

	return keyboard
}

func (h *Handler) keyboardAddMovieVideo(_ context.Context, b *bot.Bot, movieId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("افزودن کیفیت جدید", []byte(fmt.Sprintf("%v", movieId)), h.addMovieVideo)

	return keyboard
}

func (h *Handler) keyboardSeasonList(_ context.Context, b *bot.Bot, seasons []*entities2.Season) *inline.Keyboard {
	keyboard := inline.New(b)

	for _, season := range seasons {
		keyboard.Row().Button(fmt.Sprintf("%v", season.Title), []byte(fmt.Sprintf("%v", season.ID)), h.showSeason)
	}

	return keyboard
}

func (h *Handler) keyboardSeasonInfo(_ context.Context, b *bot.Bot, seasonId uint) *inline.Keyboard {
	keyboard := inline.New(b).
		Row().
		Button("قسمت‌ها", []byte(fmt.Sprintf("%v", seasonId)), h.showEpisodes).
		Button("افزودن قسمت جدید", []byte(fmt.Sprintf("%v", seasonId)), h.addEpisode)

	return keyboard
}

func (h *Handler) keyboardEpisodeList(_ context.Context, b *bot.Bot, episodes []*entities2.Episode) *inline.Keyboard {
	keyboard := inline.New(b)

	for _, episode := range episodes {
		keyboard.Row().Button(fmt.Sprintf("%v", episode.Title), []byte(fmt.Sprintf("%v", episode.ID)), h.showEpisode)
	}

	return keyboard
}

func (h *Handler) keyboardMediaQualities(_ context.Context, b *bot.Bot) *inline.Keyboard {
	keyboard := inline.New(b).Row().
		Button(fmt.Sprintf("%v", entities2.Quality480HQ), []byte(fmt.Sprintf("%v", entities2.Quality480HQ)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality720HQ), []byte(fmt.Sprintf("%v", entities2.Quality720HQ)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality1080HQ), []byte(fmt.Sprintf("%v", entities2.Quality1080HQ)), h.setQuality).
		Row().
		Button(fmt.Sprintf("%v", entities2.Quality480BluRay), []byte(fmt.Sprintf("%v", entities2.Quality480BluRay)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality720BluRay), []byte(fmt.Sprintf("%v", entities2.Quality720BluRay)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality1080BluRay), []byte(fmt.Sprintf("%v", entities2.Quality1080BluRay)), h.setQuality).
		Row().
		Button(fmt.Sprintf("%v", entities2.Quality480WebDL), []byte(fmt.Sprintf("%v", entities2.Quality480WebDL)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality720WebDL), []byte(fmt.Sprintf("%v", entities2.Quality720WebDL)), h.setQuality).
		Button(fmt.Sprintf("%v", entities2.Quality1080WebDL), []byte(fmt.Sprintf("%v", entities2.Quality1080WebDL)), h.setQuality)

	return keyboard
}

func (h *Handler) keyboardContinueQuality(_ context.Context, b *bot.Bot) *inline.Keyboard {
	keyboard := inline.New(b).Row().
		Button("اتمام", []byte(fmt.Sprintf("%v", entities2.Quality480HQ)), h.completeAddQualities)

	return keyboard
}
