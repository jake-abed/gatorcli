package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"github.com/jake-abed/gatorcli/internal/config"
	"github.com/jake-abed/gatorcli/internal/database"
	"os"
	"database/sql"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)
	appState := &state{Config: &cfg, Db: dbQueries}
	cmds := &commands{AllCommands: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", handlerFollowing)

	args := os.Args

	if len(args) < 2 {
		fmt.Println(fmt.Errorf("Argument required!"))
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	currentCmd := command{Name: cmdName, Arguments: cmdArgs}

	err = cmds.run(appState, currentCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
