package main

import (
	"fmt"
	"os"

	"github.com/goinginblind/gator-cli/cli"
	"github.com/goinginblind/gator-cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	state := cli.NewState(&cfg)
	commands := cli.NewCommands()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no command given")
		os.Exit(1)
	}
	command := &cli.Command{Name: args[1], Args: args[2:]}
	if err = commands.Run(state, *command); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
