package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/jake-abed/gatorcli/internal/database"
	"html"
	"io"
	"net/http"
	"time"
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

func scrapeFeeds(s *state) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	markFetchedParams := database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: time.Now(),
	}

	if err = s.Db.MarkFeedFetched(context.Background(), markFetchedParams); err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("* %s\n", item.Title)
	}

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "gator-cli")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := &RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return nil, err
	}

	feedTitle := html.UnescapeString(rssFeed.Channel.Title)
	feedDesc := html.UnescapeString(rssFeed.Channel.Description)

	rssFeed.Channel.Title = feedTitle
	rssFeed.Channel.Description = feedDesc

	for i, item := range rssFeed.Channel.Item {
		itemTitle := html.UnescapeString(item.Title)
		itemDesc := html.UnescapeString(item.Description)
		newItem := RSSItem{
			Title:       itemTitle,
			Description: itemDesc,
			Link:        item.Link,
			PubDate:     item.PubDate,
		}

		rssFeed.Channel.Item[i] = newItem
	}

	return rssFeed, nil
}
