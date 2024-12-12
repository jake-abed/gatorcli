package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jake-abed/gatorcli/internal/database"
	"os"
	"time"
	"strconv"
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

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		fmt.Println("Agg command expects one argument!")
		os.Exit(1)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.Arguments[0])

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <- ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	postQty := 2
	if len(cmd.Arguments) >= 1 {
		postsQty, err := strconv.ParseInt(cmd.Arguments[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		postQty = int(postsQty)
	}

	postParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(postQty),
	}
	posts, err := s.Db.GetPostsForUser(context.Background(), postParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(posts)

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title.String)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("Published: %s\n", post.PublishedAt.String)
	}
	return nil
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 2 {
		fmt.Println("addFeed expects two arguments!")
		os.Exit(1)
	}

	now := time.Now()

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID : uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(
		context.Background(), feedFollowParams,
	)
	if err != nil {
		fmt.Println("Something went wrong while trying follow the feed:")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Println("feeds expects no arguments!")
		os.Exit(1)
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed: %s | URL: %s | User: %s\n",
			feed.Name, feed.Url, feed.Name_2,
		)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Println("Follow expects exactly one argument!")
		os.Exit(1)
	}

	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(
		context.Background(), feedFollowParams,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("User <%s> now following feed <%s>", user.Name, feed.Name)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Println("Unfollow expects exactly one argument!")
		os.Exit(1)
	}

	unfollowParams := database.DeleteFeedFollowParams{
		Name: user.Name,
		Url: cmd.Arguments[0],
	}

	err := s.Db.DeleteFeedFollow(context.Background(), unfollowParams)
	if err != nil {
		fmt.Println("Could not unfollow feed:")
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		fmt.Println("Following requires no argument!")
		os.Exit(1)
	}

	followingFeeds, err := s.Db.GetFeedFollowsForUser(
		context.Background(), s.Config.CurrentUserName)
	if err != nil {
		fmt.Println("Something went wrong grabbing Follow Feeds!")
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ff := range followingFeeds {
		fmt.Printf("* %s\n", ff.FeedName)
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
