package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail creating an http request: %w", err)
	}

	client := http.DefaultClient
	request.Header.Set("User-Agent", "gator")

	response, err := client.Do(request)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail doing an http request: %w", err)
	}
	defer response.Body.Close()

	rawData, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail reading response body: %w", err)
	}

	var feed RSSFeed
	if err = xml.Unmarshal(rawData, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("fail to unmarshal: %w", err)
	}

	// unescaping the whole html thing...
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
