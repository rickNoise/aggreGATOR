package app

import (
	"context"
	"fmt"

	"github.com/rickNoise/aggreGATOR/rss"
)

// Get the next feed to fetch from the DB.
// Mark it as fetched.
// Fetch the feed using the URL.
// Iterate over the items in the feed and print their titles to the console.
func scrapeFeeds(s *State) error {
	feedToFetch, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not GetNextFeedToFetch: %w", err)
	}

	err = s.Db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return fmt.Errorf("problem marking feed as fetched: %w", err)
	}

	feedRSSData, err := rss.FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("problem fetching RSS feed: %w", err)
	}

	for _, RSSItem := range feedRSSData.Channel.Item {
		fmt.Println("-", RSSItem.Title)
	}
	return nil
}
