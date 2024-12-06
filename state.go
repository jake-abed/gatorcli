package main

import (
	"github.com/jake-abed/gatorcli/internal/config"
	"github.com/jake-abed/gatorcli/internal/database"
)

type state struct {
	Db *database.Queries
	Config *config.Config
}
