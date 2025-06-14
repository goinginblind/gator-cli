package cli

import (
	"github.com/goinginblind/gator-cli/internal/config"
	"github.com/goinginblind/gator-cli/internal/database"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func NewState(db *database.Queries, cfg *config.Config) *State {
	return &State{db: db, cfg: cfg}
}
