package imdb

import (
	"context"
	"testing"
)

func TestIMDB_GetRating(t *testing.T) {
	rating := func(t testing.TB, client *IMDB, movieId string, want float32) {
		rate, err := client.GetRating(context.Background(), movieId)
		if err != nil {
			t.Errorf("failed to get rating: %v", err)
			return
		}

		if rate.Data.Title.RatingsSummery.AggregateRating != want {
			t.Errorf("rating response mismatch: fetch(%v) != want(%v)", rate.Data.Title.RatingsSummery.AggregateRating, want)
		}
	}

	t.Run("calculating", func(t *testing.T) {
		client := NewIMDB(context.Background(), Configs{
			BaseUrl: "https://api.graphql.imdb.com/",
		})
		rating(t, client, "tt14998742", 5.6)
	})
}
