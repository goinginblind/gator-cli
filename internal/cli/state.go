package cli

import (
	"github.com/goinginblind/gator-cli/internal/config"
	"github.com/goinginblind/gator-cli/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func NewState(db *database.Queries, cfg *config.Config) *state {
	return &state{db: db, cfg: cfg}
}
