package imdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Configs struct {
	BaseUrl string `yaml:"baseUrl"`
}

type IMDB struct {
	httpClient *http.Client
	configs    Configs
}

func NewIMDB(_ context.Context, conf Configs) *IMDB {
	return &IMDB{
		httpClient: &http.Client{},
		configs:    conf,
	}
}

func (c *IMDB) GetRating(_ context.Context, movieId string) (*RateResponse, error) {
	query := fmt.Sprintf(`
query {
  title(id: "%v") {
    ratingsSummary {
      aggregateRating
    }
  }
}
`, movieId)

	req, err := http.NewRequest(http.MethodPost, c.configs.BaseUrl, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	rateResponse := new(RateResponse)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, rateResponse)
	if err != nil {
		return nil, err
	}

	return rateResponse, nil
}
