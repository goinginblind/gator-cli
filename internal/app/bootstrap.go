package app

import (
	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/app/handlers"
)

func RegisterCommands(r *common.Routes) {
	r.Register("login", handlers.Login)
	r.Register("register", handlers.Register)
	r.Register("reset", handlers.Reset)
	r.Register("users", handlers.GetUsers)
	r.Register("agg", handlers.Aggregator)
	r.Register("addfeed", handlers.CreateFeed)
	r.Register("feeds", handlers.GetFeedsWithUNames)
	r.Register("follow", handlers.CreateFollow)
	r.Register("following", handlers.GetFeedFollowsForUser)
}
