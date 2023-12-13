package repositories

import (
	"context"
	"encoding/json"
	"github.com/h-varmazyar/kiwi/applications/film/pkg/entities"
	"github.com/redis/go-redis/v9"
)

type AddContentRepository struct {
	redisClient *redis.Client
}

func NewAddContentRepository(redisClient *redis.Client) *AddContentRepository {

	return &AddContentRepository{redisClient: redisClient}
}

func (r *AddContentRepository) SetSeries(ctx context.Context, chatId int64, series *entities.Series) error {
	if _, err := r.redisClient.Set(ctx, generateKey("series", chatId), series.Json(), 0).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) GetSeries(ctx context.Context, chatId int64) (*entities.Series, error) {
	res, err := r.redisClient.Get(ctx, generateKey("series", chatId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	series := new(entities.Series)
	if err = json.Unmarshal([]byte(res), series); err != nil {
		return nil, err
	}

	return series, nil
}

func (r *AddContentRepository) DeleteSeries(ctx context.Context, chatId int64) error {
	if _, err := r.redisClient.Del(ctx, generateKey("series", chatId)).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) SetEpisode(ctx context.Context, chatId int64, episode *entities.Episode) error {
	if _, err := r.redisClient.Set(ctx, generateKey("episode", chatId), episode.Json(), 0).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) GetEpisode(ctx context.Context, chatId int64) (*entities.Episode, error) {
	res, err := r.redisClient.Get(ctx, generateKey("episode", chatId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	episode := new(entities.Episode)
	if err = json.Unmarshal([]byte(res), episode); err != nil {
		return nil, err
	}

	return episode, nil
}

func (r *AddContentRepository) DeleteEpisode(ctx context.Context, chatId int64) error {
	if _, err := r.redisClient.Del(ctx, generateKey("episode", chatId)).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) SetMedia(ctx context.Context, chatId int64, media *entities.Media) error {
	if _, err := r.redisClient.Set(ctx, generateKey("media", chatId), media.Json(), 0).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) GetMedia(ctx context.Context, chatId int64) (*entities.Media, error) {
	res, err := r.redisClient.Get(ctx, generateKey("media", chatId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	media := new(entities.Media)
	if err = json.Unmarshal([]byte(res), media); err != nil {
		return nil, err
	}

	return media, nil
}

func (r *AddContentRepository) DeleteMedia(ctx context.Context, chatId int64) error {
	if _, err := r.redisClient.Del(ctx, generateKey("media", chatId)).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) SetMovie(ctx context.Context, chatId int64, movie *entities.Movie) error {
	if _, err := r.redisClient.Set(ctx, generateKey("movie", chatId), movie.Json(), 0).Result(); err != nil {
		return err
	}
	return nil
}

func (r *AddContentRepository) GetMovie(ctx context.Context, chatId int64) (*entities.Movie, error) {
	res, err := r.redisClient.Get(ctx, generateKey("movie", chatId)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	movie := new(entities.Movie)
	if err = json.Unmarshal([]byte(res), movie); err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *AddContentRepository) DeleteMovie(ctx context.Context, chatId int64) error {
	if _, err := r.redisClient.Del(ctx, generateKey("movie", chatId)).Result(); err != nil {
		return err
	}
	return nil
}
