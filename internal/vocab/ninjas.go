package vocab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	definition_endpoint  = "http://api.api-ninjas.com/v1/dictionary"
	random_word_endpoint = "https://api.api-ninjas.com/v1/randomword"
)

type ApiNinjas struct {
	apikey  string
	timeout time.Duration
}

func NewApiNinjas(key string, connTimeout time.Duration) *ApiNinjas {
	return &ApiNinjas{
		apikey:  key,
		timeout: connTimeout,
	}
}

func (an ApiNinjas) GetDefinition(word string) (string, error) {
	ctx, done := context.WithTimeout(context.Background(), an.timeout)
	defer done()

	req, err := http.NewRequestWithContext(ctx, "GET", definition_endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Api-Key", an.apikey)
	q := req.URL.Query()
	q.Add("word", word)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading http body: %w", err)
	}

	return parseDefResponse(body)
}

func (an ApiNinjas) GetRandom() (string, error) {
	ctx, done := context.WithTimeout(context.Background(), an.timeout)
	defer done()

	req, err := http.NewRequestWithContext(ctx, "GET", random_word_endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error at creating http request: %w", err)
	}
	req.Header.Set("X-Api-Key", an.apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error at sending http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error in reading http body: %w", err)
	}

	return parseRWResponse(body)
}

func parseRWResponse(resp []byte) (string, error) {
	type result struct {
		Text []string `json:"word"`
	}

	var results result
	if err := json.Unmarshal(resp, &results); err != nil {
		log.Println("ninja returned: ", string(resp))
		return "", fmt.Errorf("unmarshaling error: %w", err)
	}

	return results.Text[0], nil
}

func parseDefResponse(resp []byte) (string, error) {
	type result struct {
		Definition string `json:"definition"`
		Word       string `json:"word"`
		Valid      bool   `json:"valid"`
	}

	var results result
	if err := json.Unmarshal(resp, &results); err != nil {
		log.Println("ninja returned: ", string(resp))
		return "", fmt.Errorf("unmarshaling error: %w", err)
	}

	return results.Definition, nil
}
