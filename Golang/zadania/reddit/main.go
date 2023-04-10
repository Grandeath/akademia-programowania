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
	"sync"
	"time"
)

var hosts = [][]string{
	{"https://www.reddit.com/r/programming.json", "output1"},
	{"https://www.reddit.com/r/golang.json", "output2"},
	{"https://www.reddit.com/r/java.json", "output3"},
}

func main() {
	var wg sync.WaitGroup
	for _, host := range hosts {
		wg.Add(1)
		go FetchAndSave(host[0], host[1], &wg)
	}
	wg.Wait()
}

type RedditClient struct {
	data fetcher.Response
}

func (r *RedditClient) Fetch(ctx context.Context, host string) error {
	client := &http.Client{}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*10))
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

	for _, child := range r.data.Data.Children {
		if _, err := fmt.Fprintf(w, "%s\n%s\n", child.Data.Title, child.Data.URL); err != nil {
			return err
		}
	}
	return nil
}

func FetchAndSave(host string, fileNmae string, wg *sync.WaitGroup) {
	var f fetcher.RedditFetcher // do not change
	var w io.Writer             // do not change

	f = &RedditClient{}
	file, err := os.Create(fileNmae)
	if err != nil {
		log.Printf("Got error: %s", err)
		return
	}
	defer file.Close()
	w = file

	if err := f.Fetch(context.Background(), host); err != nil {
		log.Printf("Got error: %s", err)
		return
	}

	if err := f.Save(w); err != nil {
		log.Printf("Got error: %s", err)
		return
	}
	log.Printf("Request completed: %s", host)
	wg.Done()
}
