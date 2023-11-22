package add

//
//import (
//	"context"
//	"github.com/go-telegram/bot"
//	"github.com/go-telegram/bot/models"
//	"github.com/h-varmazyar/kiwi/applications/film/internal/handlers"
//	repositories2 "github.com/h-varmazyar/kiwi/applications/film/internal/repositories"
//	log "github.com/sirupsen/logrus"
//)
//
//type Handler struct {
//	log            *log.Logger
//	userStateRepo  *repositories2.UserStateRepository
//	genreRepo      *repositories2.GenreRepository
//	addContentRepo *repositories2.AddContentRepository
//	seriesRepo     *repositories2.SeriesRepository
//	adminChatId    int64
//}
//
//type Dependencies struct {
//	UserStateRepo  *repositories2.UserStateRepository
//	GenreRepo      *repositories2.GenreRepository
//	AddContentRepo *repositories2.AddContentRepository
//	SeriesRepo     *repositories2.SeriesRepository
//	AdminChatId    int64
//}
//
//func NewHandler(log *log.Logger, dependencies *Dependencies) *Handler {
//	return &Handler{
//		log:            log,
//		userStateRepo:  dependencies.UserStateRepo,
//		genreRepo:      dependencies.GenreRepo,
//		addContentRepo: dependencies.AddContentRepo,
//		seriesRepo:     dependencies.SeriesRepo,
//		adminChatId:    dependencies.AdminChatId,
//	}
//}
//
//func (h *Handler) AddCmd(ctx context.Context, b *bot.Bot, update *models.Update) {
//	if update.Message.Chat.ID != h.adminChatId {
//		return
//	}
//
//	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:      update.Message.Chat.ID,
//		Text:        MsgAddContent,
//		ReplyMarkup: h.keyboardAdd(b),
//	})
//}
//
//func (h *Handler) cancelAddition(ctx context.Context, b *bot.Bot, update *models.Message, data []byte) {
//	if err := h.userStateRepo.DeleteState(ctx, update.Chat.ID); err != nil {
//		handlers.SendError(ctx, b, &handlers.ErrParams{
//			ChatId:   update.Chat.ID,
//			Err:      err,
//			Metadata: update,
//		})
//		return
//	}
//
//	handlers.SendSuccess(ctx, b, &handlers.SuccessParams{
//		ChatId:   update.Chat.ID,
//		IsSilent: false,
//	})
//}
