package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/MikkelvtK/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	uniqueConstraintError = "23505"
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

func fetchFeed(ctx context.Context, url string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return new(RSSFeed), fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return new(RSSFeed), fmt.Errorf("error fetching response: %v", err)
	}
	defer res.Body.Close()

	var feed RSSFeed
	dec := xml.NewDecoder(res.Body)

	if err = dec.Decode(&feed); err != nil {
		return new(RSSFeed), fmt.Errorf("error decoding response: %v", err)
	}

	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)

	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	fmt.Printf("fetched %d items from %s\n", len(feed.Channel.Item), feedToFetch.Url)

	for _, item := range feed.Channel.Item {
		err = storePost(s, item, feedToFetch.ID)

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueConstraintError {
			continue
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func storePost(s *state, i RSSItem, feedId uuid.UUID) error {
	publishedAt, err := time.Parse(time.RFC1123, i.PubDate)
	if err != nil {
		return err
	}

	params := database.CreatePostParams{
		ID:          uuid.New(),
		Title:       i.Title,
		Url:         i.Link,
		Description: i.Description,
		PublishedAt: publishedAt,
		FeedID:      feedId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.db.CreatePost(context.Background(), params)
}
