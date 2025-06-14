package common

import (
	"github.com/goinginblind/gator-cli/internal/config"
	"github.com/goinginblind/gator-cli/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

func NewState(db *database.Queries, cfg *config.Config) *State {
	return &State{DB: db, Config: cfg}
}
