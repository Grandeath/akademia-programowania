package main

import (
	"io"
	"reddit/fetcher"
)

func main() {
	var f fetcher.RedditFetcher
	var w io.Writer

	f.Fetch()
	f.Save(w)
}
