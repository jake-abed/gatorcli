package main

import (
	"context"
	"fmt"
	"github.com/jake-abed/gatorcli/internal/database"
	"os"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return handler(s, cmd, currentUser)
	}
}
