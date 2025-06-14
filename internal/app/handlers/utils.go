package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/database"
	"github.com/goinginblind/gator-cli/internal/rss"
)

func requireArgs(cmd common.Command, n int, usage string) error {
	if len(cmd.Args) < n {
		return fmt.Errorf("'%s' expects %d arguments", usage, n)
	}
	return nil
}

func followFeed(s *common.State, user database.User, ctx context.Context, feedUrl string) error {
	feed, err := s.DB.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("fail to get feed with url: '%v'; error: %w", feedUrl, err)
	}

	_, err = s.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err == sql.ErrNoRows {
		fmt.Printf("you are already following this feed\n")
		return nil
	} else if err != nil {
		return fmt.Errorf("fail to create a follow: %w", err)
	}

	fmt.Printf("user '%v' is now following '%v'\n", user.Name, feed.Name)
	return nil
}

func scrapeFeeds(s *common.State) error {
	ctx := context.Background()
	feed, err := s.DB.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("fail to get next feed to fetch: %w", err)
	}
	if err = s.DB.MarkFeedFetched(ctx, feed.ID); err != nil {
		return fmt.Errorf("fail to mark feed as fetched: %w", err)
	}
	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		nullTime := toNullTime(item.PubDate)
		_, err = s.DB.CreatePost(ctx, database.CreatePostParams{
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: nullTime,
			FeedID:      feed.ID,
		})
		if err == nil {
			fmt.Printf("post '%v' saved\n", item.Title)
		} else if err != sql.ErrNoRows {
			fmt.Printf("fail to create post: %v\n", err)
		}
	}
	return nil
}

func toNullTime(s string) sql.NullTime {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822,
		time.RFC3339,
	}
	for _, f := range formats {
		if parsed, err := time.Parse(f, s); err == nil {
			return sql.NullTime{Time: parsed, Valid: true}
		}
	}
	return sql.NullTime{}
}
