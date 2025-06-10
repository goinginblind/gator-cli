package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/goinginblind/gator-cli/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd Command) error {
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

func handlerRegister(s *state, cmd Command) error {
	// check if there are arguments at all
	if len(cmd.Args) == 0 {
		return fmt.Errorf("'register' expects a single argument")
	}
	ctx := context.Background()

	// check if the name is already taken
	if _, err := s.db.GetUserByName(ctx, cmd.Args[0]); err == nil {
		return fmt.Errorf("username '%v' already taken", cmd.Args[0])
	}

	// create a new user struct for db
	newUser := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]}
	// create new user in the db
	_, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("fail to create a user in database: %w", err)
	}
	// set new current user
	s.cfg.SetUser(cmd.Args[0])
	fmt.Printf("user '%v' has been created and set as current\n", cmd.Args[0])
	return nil
}

func handlerReset(s *state, cmd Command) error {
	ctx := context.Background()
	if err := s.db.ResetRows(ctx); err != nil {
		return fmt.Errorf("failt to reset rows: %w", err)
	}
	fmt.Printf("table rows have been reset\n")
	return nil
}

func handlerGetUsers(s *state, cmd Command) error {
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
