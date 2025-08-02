package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rickNoise/aggreGATOR/internal/database"
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

	numItemsToProcess := len(feedRSSData.Channel.Item)
	var numNewPostsCreated, numDuplicatePostsSkipped int

	fmt.Printf("Iterating over items for feed %s...\n", feedToFetch.Name)
	for _, RSSItem := range feedRSSData.Channel.Item {
		// Build database item params
		dbItem := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       parseStringForNullString(RSSItem.Title),
			Url:         RSSItem.Link,
			Description: parseStringForNullString(RSSItem.Description),
			PublishedAt: parseStringForNullTime(RSSItem.PubDate),
			FeedID:      feedToFetch.ID,
		}

		_, err := s.Db.CreatePost(context.Background(), dbItem)
		if err != nil {
			if checkForUniqueConstraintViolation(err) {
				numDuplicatePostsSkipped++
				continue
			} else {
				// If it's any other error, log it and stop the entire scrape.
				// This is important for critical errors like a bad connection or syntax error.
				return fmt.Errorf("could not create post for title %s: %w", RSSItem.Title, err)
			}
		} else {
			fmt.Printf(">> new post created: - %s\n", dbItem.Title.String)
			numNewPostsCreated++
		}
	}
	fmt.Printf("### from total %d feed items, created %d new posts and skipped %d already saved posts\n", numItemsToProcess, numNewPostsCreated, numDuplicatePostsSkipped)
	return nil
}

func checkForUniqueConstraintViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func parseStringForNullString(str string) sql.NullString {
	if str != "" {
		return sql.NullString{String: str, Valid: true}
	} else {
		return sql.NullString{Valid: false}
	}
}

func parseStringForNullTime(dateStr string) sql.NullTime {
	// The layouts to try in order. We'll start with the most common.
	dateLayouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC3339,
	}

	// First, handle the case where the string is empty.
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}

	// Iterate through our list of known formats and try to parse the string.
	for _, layout := range dateLayouts {
		parsedTime, err := time.Parse(layout, dateStr)
		if err == nil {
			// Success! Return a valid sql.NullTime.
			return sql.NullTime{Time: parsedTime, Valid: true}
		}
	}

	// Parsing unsuccessful with all our predefined layouts, return NULL database value.
	return sql.NullTime{Valid: false}
}
