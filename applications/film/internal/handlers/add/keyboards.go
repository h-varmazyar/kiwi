package add

//
//import (
//	"context"
//	"fmt"
//	"github.com/go-telegram/bot"
//	"github.com/go-telegram/ui/keyboard/inline"
//	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
//)
//
//func (h *Handler) keyboardAdd(b *bot.Bot) *inline.Keyboard {
//	return inline.New(b).
//		Row().
//		Button("سریال", []byte("سریال"), h.addSeries).
//		Button("قسمت سریال", []byte("قسمت سریال"), h.addEpisode).
//		Row().
//		Button("سینمایی", []byte("سینمایی"), h.addMovie).
//		Button("افزودن کیفیت", []byte("افزودن کیفیت"), h.addQuality).
//		Row().
//		Button("لغو", []byte{}, h.cancelAddition)
//}
//
//func (h *Handler) keyboardCancel(b *bot.Bot) *inline.Keyboard {
//	return inline.New(b).
//		Row().
//		Button("لغو", []byte{}, h.cancelAddition)
//}
//
//func (h *Handler) keyboardGenre(ctx context.Context, b *bot.Bot) (*inline.Keyboard, error) {
//	genreKeyboard := inline.New(b).Row()
//	genres, err := h.genreRepo.List(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, genre := range genres {
//		genreKeyboard.Button(genre.FaName, []byte(fmt.Sprintf("%v", genre.ID)), h.selectGenre)
//	}
//
//	genreKeyboard.Row().Button("اتمام", []byte(""), h.submitSeries)
//
//	return genreKeyboard, nil
//}
//
//func (h *Handler) keyboardSearchedSeries(_ context.Context, b *bot.Bot, seriesList []*entities.Series) (*inline.Keyboard, error) {
//	seriesKeyboard := inline.New(b)
//
//	for _, series := range seriesList {
//		seriesKeyboard.
//			Row().
//			Button(fmt.Sprintf("%v | %v", series.FaName, series.EnName), []byte(fmt.Sprintf("%v", series.ID)), h.selectSeries)
//	}
//
//	seriesKeyboard.Row().Button("لغو", []byte(""), h.cancelAddition)
//
//	return seriesKeyboard, nil
//}
//
////func (h *Handler) keyboardMediaType(b *bot.Bot) *inline.Keyboard {
////	return inline.New(b).
////		Row().
////		Button(entities.MediaSeries, []byte(entities.MediaSeries), h.selectMediaType).
////		Button(entities.MediaMovie, []byte(entities.MediaMovie), h.selectMediaType).
////		Row().
////		Button("هیچکدام", []byte{}, h.selectMediaType).
////		Row().
////		Button("لغو", []byte{}, h.cancelAddition)
////}
////
////func (h *Handler) keyboardQualities(b *bot.Bot) *inline.Keyboard {
////	return inline.New(b).
////		Row().
////		Button(entities.Quality480HQ, []byte(entities.Quality480HQ), h.selectQuality).
////		Button(entities.Quality720HQ, []byte(entities.Quality720HQ), h.selectQuality).
////		Button(entities.Quality1080HQ, []byte(entities.Quality1080HQ), h.selectQuality).
////		Row().
////		Button(entities.Quality480BluRay, []byte(entities.Quality480BluRay), h.selectQuality).
////		Button(entities.Quality720BluRay, []byte(entities.Quality720BluRay), h.selectQuality).
////		Button(entities.Quality1080BluRay, []byte(entities.Quality1080BluRay), h.selectQuality).
////		Row().
////		Button(entities.Quality480WebDL, []byte(entities.Quality480WebDL), h.selectQuality).
////		Button(entities.Quality720WebDL, []byte(entities.Quality720WebDL), h.selectQuality).
////		Button(entities.Quality1080WebDL, []byte(entities.Quality1080WebDL), h.selectQuality).
////		Row().
////		Button("هیچکدام", []byte{}, h.selectQuality).
////		Row().
////		Button("لغو", []byte{}, h.cancelAddition)
////}
