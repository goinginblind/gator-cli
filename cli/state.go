package cli

import "github.com/goinginblind/gator-cli/internal/config"

type state struct {
	cfg *config.Config
}

func NewState(cfg *config.Config) *state {
	return &state{cfg: cfg}
}
