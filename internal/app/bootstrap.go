package app

import (
	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/app/handlers"
	"github.com/goinginblind/gator-cli/internal/app/middleware"
)

func RegisterCommands(r *common.Routes) {
	r.Register("login", handlers.Login)
	r.Register("register", handlers.Register)
	r.Register("reset", handlers.Reset)
	r.Register("users", handlers.GetUsers)
	r.Register("agg", handlers.Aggregator)
	r.Register("addfeed", middleware.LoggedIn(handlers.AddFeed))
	r.Register("feeds", handlers.GetFeedsWithUNames)
	r.Register("follow", middleware.LoggedIn(handlers.CreateFollow))
	r.Register("following", middleware.LoggedIn(handlers.UsersFollows))
	r.Register("unfollow", middleware.LoggedIn(handlers.DeleteFollow))
}
