package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/goinginblind/gator-cli/internal/database"
	"github.com/goinginblind/gator-cli/internal/rss"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("'login' expects a single argument")
	}

	ctx := context.Background()
	// check if the name is already taken
	if _, err := s.db.GetUserByName(ctx, cmd.Args[0]); err == sql.ErrNoRows {
		return fmt.Errorf("username '%v' does not exist", cmd.Args[0])
	} else if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("the user has been set as '%v'\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	// check if there are arguments at all
	if len(cmd.Args) == 0 {
		return fmt.Errorf("'register' expects a single argument")
	}
	ctx := context.Background()
	username := strings.TrimSpace(cmd.Args[0])

	// check if the name is already taken
	if _, err := s.db.GetUserByName(ctx, username); err == nil {
		return fmt.Errorf("username '%v' already taken", username)
	}

	// create new user in the db
	_, err := s.db.CreateUser(ctx, username)
	if err != nil {
		return fmt.Errorf("fail to create a user in database: %w", err)
	}
	// set new current user
	s.cfg.SetUser(username)
	fmt.Printf("user '%v' has been created and set as current\n", username)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	if err := s.db.ResetRows(ctx); err != nil {
		return fmt.Errorf("failt to reset rows: %w", err)
	}
	fmt.Printf("table rows have been reset\n")
	return nil
}

func handlerGetUsers(s *State, cmd Command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("fail to get users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAggregator(s *State, cmd Command) error {
	ctx := context.Background()
	res, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	for i := range res.Channel.Item {
		fmt.Print(res.Channel.Item[i].Title)
		fmt.Println()
	}
	return nil
}

func handlerCreateFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("'addfeed' expects two arguments")
	}
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("fail to get users id: %w", err)
	}
	feedParams := database.CreateFeedParams{
		Name:   strings.TrimSpace(cmd.Args[0]),
		Url:    strings.TrimSpace(cmd.Args[1]),
		UserID: user.ID,
	}

	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("fail to create feed: %w", err)
	}

	fmt.Printf("ID:		%s \n", feed.ID)
	fmt.Printf("Name:		%s \n", feed.Name)
	fmt.Printf("UserID:		%s \n", feed.UserID)
	fmt.Printf("URL:		%s \n", feed.Url)
	fmt.Printf("Created at: 	%s \n", feed.CreatedAt)

	if err = followFeed(s, ctx, feed.Url); err != nil {
		return fmt.Errorf("feed created, but fail to follow it: %w", err)
	}

	return nil
}

func handlerGetFeedsWithUNames(s *State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedsWithUNames(ctx)
	if err != nil {
		return fmt.Errorf("fail to get feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("feed name: %v\nurl: %v\nadded by: %v\n\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerCreateFollow(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("'follow' expects one argument")
	}
	feedUrl := strings.TrimSpace(cmd.Args[0])
	ctx := context.Background()
	return followFeed(s, ctx, feedUrl)
}

func handlerGetFeedFollowsForUser(s *State, cmd Command) error {
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("fail to get user: %w", err)
	}
	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("fail to get feed follows: %w", err)
	}
	for i, feed := range feeds {
		fmt.Printf("%v. %v\n", i+1, feed)
	}
	return nil
}
