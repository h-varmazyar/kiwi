package add

//
//import (
//	"context"
//	"fmt"
//	"github.com/go-telegram/bot"
//	"github.com/go-telegram/bot/models"
//	"github.com/h-varmazyar/kiwi/applications/film/internal/handlers"
//	"github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
//	entities2 "github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
//	"strconv"
//	"strings"
//)
//
//func (h *Handler) addMediaCmd(ctx context.Context, b *bot.Bot, update *models.Update) {
//	userId := update.Message.Chat.ID
//	if userId != h.adminChatId {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   userId,
//			Msg:      MsgUnauthorized,
//			Err:      ErrUnauthorized,
//			Metadata: update,
//		})
//		return
//	}
//
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:           userId,
//		Text:             MsgAddContent,
//		ReplyToMessageID: update.Message.ID,
//		ReplyMarkup:      h.keyboardAdd(b),
//	})
//}
//
//func (h *Handler) addSeries(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddSeries); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   update.Chat.ID,
//			Err:      err,
//			Metadata: update,
//		})
//		return
//	}
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID: update.Chat.ID,
//		Text:   MsgAddSeries,
//	})
//}
//func (h *Handler) addMovie(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//
//}
//func (h *Handler) addQuality(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//
//}
//
//func (h *Handler) selectMediaType(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	////todo: set state
//	//switch string(data) {
//	//case entities.MediaMovie:
//	//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//	//		ChatID:      update.Chat.ID,
//	//		Text:        "کیفیت فیلم را انتخاب کنید",
//	//		ReplyMarkup: h.keyboardQualities(b),
//	//	})
//	//case entities.MediaSeries:
//	//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//	//		ChatID: update.Chat.ID,
//	//		Text:   "نام سریال را وارد کنید",
//	//	})
//	//default:
//	//
//	//}
//}
//
//func (h *Handler) selectQuality(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	//switch string(data) {
//	//case entities.Quality480HQ:
//	//case entities.Quality720HQ:
//	//case entities.Quality1080HQ:
//	//case entities.Quality480BluRay:
//	//case entities.Quality720BluRay:
//	//case entities.Quality1080BluRay:
//	//case entities.Quality480WebDL:
//	//case entities.Quality720WebDL:
//	//case entities.Quality1080WebDL:
//	//default:
//	//
//	//}
//}
//
//func (h *Handler) AddSeriesState(ctx context.Context, b *bot.Bot, update *models.Update) {
//	chatId := update.Message.Chat.ID
//	data := strings.Split(strings.TrimSpace(update.Message.Text), "\n")
//	fmt.Println("data:", len(data), data[0])
//	if len(data) != 6 {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgInvalidSeriesData,
//			Err:      ErrInvalidSeriesData,
//			Metadata: update,
//		})
//		return
//	}
//	series := &entities2.Series{
//		Title:        strings.TrimSpace(data[0]),
//		EnName:       strings.TrimSpace(data[2]),
//		FaName:       strings.TrimSpace(data[1]),
//		Presentation: strings.TrimSpace(data[3]),
//		Year:         strings.TrimSpace(data[4]),
//		IMDBLink:     strings.TrimSpace(data[5]),
//	}
//
//	if err := h.addContentRepo.SetSeries(ctx, chatId, series); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Err:      err,
//			Metadata: series,
//		})
//		return
//	}
//
//	if err := h.userStateRepo.SetState(ctx, chatId, repositories.StateAddBanner); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Err:      err,
//			Metadata: "cannot move to state " + repositories.StateAddBanner,
//		})
//		return
//	}
//
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:      chatId,
//		Text:        MsgAddSeriesBanner,
//		ReplyMarkup: h.keyboardCancel(b),
//	})
//}
//
//func (h *Handler) AddBannerState(ctx context.Context, b *bot.Bot, update *models.Update) {
//	fmt.Println("add banner state")
//	chatId := update.Message.Chat.ID
//
//	if update.Message.Photo == nil || len(update.Message.Photo) == 0 {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgNoBannerAdded,
//			Err:      ErrNoBannerAdded,
//			Metadata: update,
//		})
//		return
//	}
//
//	series, err := h.addContentRepo.GetSeries(ctx, chatId)
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgNoSeriesCached,
//			Err:      ErrNoSeriesCached,
//			Metadata: update,
//		})
//		return
//	}
//
//	series.Banner = &entities2.Media{
//		TelegramFileId: update.Message.Photo[0].FileID,
//		Type:           entities2.MediaTypePicture,
//	}
//
//	if err = h.addContentRepo.SetSeries(ctx, chatId, series); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Err:      err,
//			Metadata: series,
//		})
//		return
//	}
//
//	if err = h.userStateRepo.SetState(ctx, chatId, repositories.StateAddBanner); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Err:      err,
//			Metadata: "cannot move to state " + repositories.StateAddBanner,
//		})
//		return
//	}
//
//	genreKeyboard, err := h.keyboardGenre(ctx, b)
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			IsSilent: true,
//			Err:      err,
//		})
//		h.submitSeries(ctx, b, update.Message, nil)
//		return
//	}
//
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:      chatId,
//		Text:        MsgSelectGenre,
//		ReplyMarkup: genreKeyboard,
//	})
//}
//
//func (h *Handler) selectGenre(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	chatId := update.Chat.ID
//	series, err := h.addContentRepo.GetSeries(ctx, chatId)
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgNoSeriesCached,
//			Err:      ErrNoSeriesCached,
//			Metadata: update,
//		})
//		return
//	}
//
//	text := MsgSelectGenre
//	if series.Genres == nil {
//		series.Genres = make([]*entities2.Genre, 0)
//	}
//
//	id, err := strconv.Atoi(string(data))
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Err:      err,
//			Metadata: update,
//		})
//	} else {
//		genre, err := h.genreRepo.Return(ctx, uint(id))
//		if err != nil {
//			handlers.SendError(ctx, b, &handlers.ErrParams{
//				ChatId:   chatId,
//				Err:      err,
//				Metadata: update,
//			})
//		} else {
//			series.Genres = append(series.Genres, genre)
//			if err = h.userStateRepo.SetState(ctx, chatId, repositories.StateAddBanner); err != nil {
//				handlers.SendError(ctx, b, &handlers.ErrParams{
//					Err:      err,
//					IsSilent: true,
//					Metadata: "cannot move to state " + repositories.StateAddBanner,
//				})
//			}
//			if err = h.addContentRepo.SetSeries(ctx, chatId, series); err != nil {
//				handlers.SendError(ctx, b, &handlers.ErrParams{
//					Err:      err,
//					IsSilent: true,
//					Metadata: "cannot save series",
//				})
//			}
//		}
//	}
//	text = fmt.Sprintf("%v\nژانرهای انتخاب شده:", MsgSelectGenre)
//	for i, genre := range series.Genres {
//		text = fmt.Sprintf("%v\n%v- %v", text, i+1, genre.FaName)
//	}
//
//	genreKeyboard, err := h.keyboardGenre(ctx, b)
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			IsSilent: true,
//			Err:      err,
//		})
//	}
//
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:      chatId,
//		Text:        text,
//		ReplyMarkup: genreKeyboard,
//	})
//}
//
//func (h *Handler) submitSeries(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	chatId := update.Chat.ID
//	series, err := h.addContentRepo.GetSeries(ctx, chatId)
//	if err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgNoSeriesCached,
//			Err:      ErrNoSeriesCached,
//			Metadata: update,
//		})
//		return
//	}
//
//	if err := h.seriesRepo.Create(ctx, series); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   chatId,
//			Msg:      MsgSaveSeriesFailed,
//			Err:      ErrSaveSeriesFailed,
//			Metadata: update,
//		})
//		return
//	}
//
//	if err := h.userStateRepo.DeleteState(ctx, chatId); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			Err:      err,
//			IsSilent: true,
//		})
//	}
//
//	if err := h.addContentRepo.DeleteSeries(ctx, chatId); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			Err:      err,
//			IsSilent: true,
//		})
//	}
//
//	genres := ""
//	for _, genre := range series.Genres {
//		genres = fmt.Sprintf("%v | %v", genres, genre.FaName)
//	}
//	genres = strings.Trim(strings.TrimSpace(genres), "|")
//
//	MsgSeriesSaved = fmt.Sprintf(
//		MsgSeriesSaved,
//		series.Title,
//		series.FaName,
//		series.EnName,
//		series.Presentation,
//		series.Year,
//		series.IMDBLink,
//		genres,
//	)
//
//	photo := &models.InputFileString{
//		Data: series.Banner.TelegramFileId,
//	}
//
//	_, _ = b.SendPhoto(ctx, &bot.SendPhotoParams{
//		ChatID:  chatId,
//		Photo:   photo,
//		Caption: MsgSeriesSaved,
//	})
//}
