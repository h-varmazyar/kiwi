package handlers

type SearchResultType string

const (
	SearchResultSeries SearchResultType = "series"
	SearchResultMovie  SearchResultType = "movie"
)

type SearchResult struct {
	FaName string
	EnName string
	Id     uint
	Type   SearchResultType
}
