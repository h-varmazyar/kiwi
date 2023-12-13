package repositories

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type State string

const (
	StateAddSeries        = "add_series"
	StateAddMovie         = "add_movie"
	StateAddMovieBanner   = "add_movie_banner"
	StateAddEpisode       = "add_episode"
	StateAddEpisodeBanner = "add_episode_banner"
	StateAddMedia         = "add_media"
)

type UserStateRepository struct {
	redisClient *redis.Client
}

func NewUserState(redisClient *redis.Client) *UserStateRepository {
	return &UserStateRepository{redisClient: redisClient}
}

func (r *UserStateRepository) SetState(ctx context.Context, userId int64, state State) error {
	_, err := r.redisClient.Set(ctx, generateKey("state", userId), []byte(state), 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserStateRepository) GetState(ctx context.Context, userId int64) (State, error) {
	res, err := r.redisClient.Get(ctx, generateKey("state", userId)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return State(res), nil
}

func (r *UserStateRepository) DeleteState(ctx context.Context, userId int64) error {
	_, err := r.redisClient.Del(ctx, generateKey("state", userId)).Result()
	if err != nil {
		return err
	}
	return nil
}
