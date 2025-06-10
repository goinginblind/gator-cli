package main

import (
	"fmt"
	"os"

	"github.com/goinginblind/gator-cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if err = cfg.SetUser("Sam"); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println(cfg.CurrentUserName)
	fmt.Println(cfg.DbURL)
}
