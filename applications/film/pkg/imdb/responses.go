package imdb

type RateResponse struct {
	Data struct {
		Title struct {
			RatingsSummery struct {
				AggregateRating float32
			}
		}
	}
}
