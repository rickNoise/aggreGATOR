package app

import (
	"github.com/rickNoise/aggreGATOR/internal/config"
	"github.com/rickNoise/aggreGATOR/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}
