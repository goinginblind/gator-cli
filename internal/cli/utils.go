package cli

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goinginblind/gator-cli/internal/database"
)

func followFeed(s *State, ctx context.Context, feedUrl string) error {
	feed, err := s.db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("fail to get feed with url: '%v'; error: %w", feedUrl, err)
	}
	user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("fail to get user: %w", err)
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err == sql.ErrNoRows {
		fmt.Printf("you are already following this feed\n")
		return nil
	} else if err != nil {
		return fmt.Errorf("fail to create a follow: %w", err)
	}

	fmt.Printf("user '%v' is now following '%v'\n", user.Name, feed.Name)
	return nil
}
