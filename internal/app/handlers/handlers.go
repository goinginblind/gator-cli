package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/database"
)

func Login(s *common.State, cmd common.Command) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}

	ctx := context.Background()
	// check if the name is already taken
	if _, err := s.DB.GetUserByName(ctx, cmd.Args[0]); err == sql.ErrNoRows {
		return fmt.Errorf("username '%v' does not exist", cmd.Args[0])
	} else if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if err := s.Config.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Printf("the user has been set as '%v'\n", cmd.Args[0])
	return nil
}

func Register(s *common.State, cmd common.Command) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}

	ctx := context.Background()
	username := strings.TrimSpace(cmd.Args[0])

	// check if the name is already taken
	if _, err := s.DB.GetUserByName(ctx, username); err == nil {
		return fmt.Errorf("username '%v' already taken", username)
	}

	// create new user in the db
	_, err := s.DB.CreateUser(ctx, username)
	if err != nil {
		return fmt.Errorf("fail to create a user in database: %w", err)
	}
	// set new current user
	s.Config.SetUser(username)
	fmt.Printf("user '%v' has been created and set as current\n", username)
	return nil
}

func Reset(s *common.State, cmd common.Command) error {
	ctx := context.Background()
	if err := s.DB.ResetRows(ctx); err != nil {
		return fmt.Errorf("fail to reset rows: %w", err)
	}
	fmt.Printf("table rows have been reset\n")
	return nil
}

func GetUsers(s *common.State, cmd common.Command) error {
	ctx := context.Background()
	users, err := s.DB.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("fail to get users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func Aggregator(s *common.State, cmd common.Command) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}
	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("fail to parse time between reqs: %w", err)
	}
	fmt.Printf("collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func AddFeed(s *common.State, cmd common.Command, user database.User) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}
	ctx := context.Background()
	feedParams := database.CreateFeedParams{
		Name:   strings.TrimSpace(cmd.Args[0]),
		Url:    strings.TrimSpace(cmd.Args[1]),
		UserID: user.ID,
	}

	feed, err := s.DB.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("fail to create feed: %w", err)
	}

	fmt.Printf("ID:		%s \n", feed.ID)
	fmt.Printf("Name:		%s \n", feed.Name)
	fmt.Printf("UserID:		%s \n", feed.UserID)
	fmt.Printf("URL:		%s \n", feed.Url)
	fmt.Printf("Created at: 	%s \n", feed.CreatedAt)

	if err = followFeed(s, user, ctx, feed.Url); err != nil {
		return fmt.Errorf("feed created, but fail to follow it: %w", err)
	}

	return nil
}

func GetFeedsWithUNames(s *common.State, cmd common.Command) error {
	ctx := context.Background()
	feeds, err := s.DB.GetFeedsWithUNames(ctx)
	if err != nil {
		return fmt.Errorf("fail to get feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("feed name: %v\nurl: %v\nadded by: %v\n\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func CreateFollow(s *common.State, cmd common.Command, user database.User) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}
	feedUrl := strings.TrimSpace(cmd.Args[0])
	ctx := context.Background()
	return followFeed(s, user, ctx, feedUrl)
}

func UsersFollows(s *common.State, cmd common.Command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.DB.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("fail to get feed follows: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Printf("no feeds followed")
		return nil
	}
	for i, feed := range feeds {
		fmt.Printf("%v. %v\n", i+1, feed)
	}
	return nil
}

func DeleteFollow(s *common.State, cmd common.Command, user database.User) error {
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		return err
	}
	ctx := context.Background()
	feed, err := s.DB.GetFeedByUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("fail to get feed by url: %w", err)
	}
	err = s.DB.DeleteFollow(ctx, database.DeleteFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("fail to delete to unfollow: %w", err)
	}
	fmt.Printf("unfollowed feed %s\n", feed.Name)
	return nil
}

func Browse(s *common.State, cmd common.Command, user database.User) error {
	var lim int
	if err := requireArgs(cmd, 1, cmd.Name); err != nil {
		lim = 2
	} else {
		lim, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("%v is not a valid limit value", cmd.Args[0])
		}
	}

	ctx := context.Background()
	posts, err := s.DB.GetPostsUser(ctx, database.GetPostsUserParams{ID: user.ID, Limit: int32(lim)})
	if err != nil {
		return fmt.Errorf("fail to fetch posts for user: %w", err)
	}
	for _, post := range posts {
		fmt.Printf("\033[31m------- %v ------\033[0m\n", post.Title)
		fmt.Printf("\033[33m* URL: %v\033[0m\n", post.Url)
		fmt.Printf("* Description: %v\n", post.Description.String)
		fmt.Printf("* Published on: %v\n", post.PublishedAt.Time)
		fmt.Println("\033[31m------------------------------------------\033[0m")
		fmt.Println()
		fmt.Println()
	}
	return nil
}
