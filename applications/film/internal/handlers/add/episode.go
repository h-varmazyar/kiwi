package add

//
//import (
//	"context"
//	"github.com/go-telegram/bot"
//	"github.com/go-telegram/bot/models"
//	"github.com/h-varmazyar/kiwi/applications/film/internal/handlers"
//	"github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
//)
//
//func (h *Handler) addEpisode(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	if err := h.userStateRepo.SetState(ctx, update.Chat.ID, repositories.StateAddEpisode); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   update.Chat.ID,
//			Err:      err,
//			Metadata: update,
//		})
//		return
//	}
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID: update.Chat.ID,
//		Text:   MsgSelectSeriesForEpisode,
//	})
//}
//
//func (h *Handler) selectSeries(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	//chatId := update.Chat.ID
//	//seriesId, err := strconv.Atoi(string(data))
//	//if err != nil {
//	//	helpers.SendError(ctx, b, &helpers.ErrParams{
//	//		ChatId:   chatId,
//	//		Err:      err,
//	//		Metadata: update,
//	//	})
//	//	return
//	//}
//	//
//	//series, err := h.seriesRepo.Return(ctx, uint(seriesId))
//	//if err != nil {
//	//	helpers.SendError(ctx, b, &helpers.ErrParams{
//	//		ChatId:   chatId,
//	//		Err:      err,
//	//		Metadata: update,
//	//	})
//	//	return
//	//}
//	//
//	//h.addContentRepo.GetSeries()
//	//
//	//_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//	//	ChatID: update.Chat.ID,
//	//	Text:   MsgSelectSeriesForEpisode,
//	//})
//}
//
//func (h *Handler) AddEpisodeState(ctx context.Context, b *bot.Bot, update *models.Update) {
//	//series, err := h.seriesRepo.Search(ctx, update.Message.Text)
//	//if err != nil {
//	//	helpers.SendError(ctx, b, &helpers.ErrParams{
//	//		ChatId:   update.Message.Chat.ID,
//	//		Err:      err,
//	//		Metadata: update,
//	//	})
//	//	return
//	//}
//	//
//	//if len(series) == 0 {
//	//	helpers.SendError(ctx, b, &helpers.ErrParams{
//	//		ChatId:   update.Message.Chat.ID,
//	//		Err:      ErrNoSeriesFound,
//	//		Metadata: update,
//	//	})
//	//	return
//	//}
//	//
//	//_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//	//	ChatID:      update.Message.Chat.ID,
//	//	Text:        MsgSelectSeriesBetweenSearch,
//	//	ReplyMarkup: h.keyboardSearchedSeries(ctx, b, series),
//	//})
//}
