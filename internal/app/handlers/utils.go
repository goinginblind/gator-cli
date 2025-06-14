package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/database"
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
