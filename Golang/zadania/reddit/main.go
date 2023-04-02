package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reddit/fetcher"
	"time"
)

const (
	host = "https://www.reddit.com/r/golang.json"
)

func main() {
	var f fetcher.RedditFetcher // do not change
	var w io.Writer             // do not change

	f = &RedditClient{}

	if err := f.Fetch(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := f.Save(w); err != nil {
		log.Fatal(err)
	}

}

type RedditClient struct {
	data fetcher.Response
}

func (r *RedditClient) Fetch(ctx context.Context) error {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host, http.NoBody)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", nil)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&r.data)
	if err != nil {
		return err
	}
	return err
}
func (r RedditClient) Save(w io.Writer) error {
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, child := range r.data.Data.Children {
		if _, err := file.WriteString(fmt.Sprintf("%s\n%s\n", child.Data.Title, child.Data.URL)); err != nil {
			return err
		}
	}
	return nil
}
