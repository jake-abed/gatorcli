package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jake-abed/gatorcli/internal/database"
	"os"
	"time"
)

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	AllCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.AllCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("Login handler expects a user argument!")
	}
	user, err := s.Db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %s \n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("Register handler expects a user argument!")
	}

	if len(cmd.Arguments) != 1 {
		return errors.New("Login handler expects exactly one name arg!")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Arguments[0],
	}

	user, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Println("Hmm! It looks like that user may already exist?!")
		os.Exit(1)
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("User has been created")
	fmt.Printf("User has been set to: %s \n", cmd.Arguments[0])
	fmt.Println(user)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Println("Users command does not accept arguments!")
		os.Exit(1)
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	username := s.Config.CurrentUserName

	for _, user := range users {
		if username == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Println("Reset does not accept arguments!")
		os.Exit(1)
	}

	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Users successfully reset in Gator DB!")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	cmdFn, ok := c.AllCommands[cmd.Name]
	if !ok {
		msg := fmt.Sprintf("Command <%s> not found!", cmd.Name)
		return errors.New(msg)
	}
	err := cmdFn(s, cmd)
	return err
}
