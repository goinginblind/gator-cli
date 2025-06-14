package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/goinginblind/gator-cli/internal/app"
	"github.com/goinginblind/gator-cli/internal/app/common"
	"github.com/goinginblind/gator-cli/internal/config"
	"github.com/goinginblind/gator-cli/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbURL := cfg.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state := common.NewState(dbQueries, cfg)
	cmds := common.NewRoutes()
	app.RegisterCommands(cmds)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no command provided")
		os.Exit(1)
	}
	cmd := common.Command{Name: args[1], Args: args[2:]}
	if err := cmds.Run(state, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
