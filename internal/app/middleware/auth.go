package middleware

import (
	"context"
	"fmt"

	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/database"
)

type LoggedInHandler func(
	s *common.State,
	cmd common.Command,
	user database.User,
) error

type CommandHandler func(
	s *common.State,
	cmd common.Command,
) error

func LoggedIn(handler LoggedInHandler) CommandHandler {
	return func(s *common.State, cmd common.Command) error {
		ctx := context.Background()
		user, err := s.DB.GetUserByName(ctx, s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("fail to get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
