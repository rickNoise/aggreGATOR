package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	fmt.Printf("attempting to FetchFeed from %s...\n", feedURL)

	req, err := http.NewRequestWithContext(
		context.Background(),
		"GET",
		feedURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	fmt.Println("Using http.DefaultClient...")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	fmt.Printf("DefaultClient Response Status: %s\n", res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var resJSON RSSFeed
	err = xml.Unmarshal(resBody, &resJSON)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling xml: %w", err)
	}

	// decode escaped HTML entities (like &ldquo;)
	// need to run the Title and Description fields (of both the entire channel as well as the items) through this function.
	// Channel fields
	resJSON.Channel.Title = html.UnescapeString(resJSON.Channel.Title)
	resJSON.Channel.Description = html.UnescapeString(resJSON.Channel.Description)
	// Item fields
	for i := range resJSON.Channel.Item {
		resJSON.Channel.Item[i].Title = html.UnescapeString(resJSON.Channel.Item[i].Title)
		resJSON.Channel.Item[i].Description = html.UnescapeString(resJSON.Channel.Item[i].Description)
	}

	return &resJSON, nil
}
