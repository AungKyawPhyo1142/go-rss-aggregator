package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AungKyawPhyo1142/rss-aggregator/internal"
	"github.com/AungKyawPhyo1142/rss-aggregator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type State struct {
	config *internal.Config
	db     *database.Queries
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	if err := c.handlers[cmd.name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) {

	c.handlers[name] = f

}

func handlerLogin(s *State, cmd Command) error {

	if len(cmd.args) <= 0 {
		return errors.New("login handler expects <username> as argument")
	}

	username := cmd.args[0]

	_, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {

		if err == sql.ErrNoRows {
			fmt.Println("user does not exist")
			os.Exit(1)
		}
		return err
	}

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("user has been set & logged in")

	return nil

}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) <= 0 {
		return errors.New("login handler expects <username> as argument")
	}

	name := cmd.args[0]

	dbUser, err := s.db.GetUserByName(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			user := database.CreateUserParams{
				Name:      name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ID:        uuid.New(),
			}
			u, err := s.db.CreateUser(context.Background(), user)
			if err != nil {
				fmt.Printf("error creating user: %v\n", err)
				return err
			}
			if err := s.config.SetUser(u.Name); err != nil {
				fmt.Printf("error setting user: %v\n", err)
				return err
			}

			fmt.Printf("user created: %v\n", u)
		} else {
			fmt.Printf("error getting user: %v\n", err)
			return err
		}
	} else {
		// found a user, so reject
		fmt.Printf("user already exists: %v\n", dbUser)
		return errors.New("user already exists")
	}

	return nil
}

func handlerReset(s *State, cmd Command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

func main() {

	var state State

	cfg, err := internal.Read()
	if err != nil {
		fmt.Println(err)
	}

	state.config = &cfg

	// open db connection
	db, err := sql.Open("postgres", state.config.Db_Url)
	if err != nil {
		fmt.Printf("error opening db connection: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	state.db = dbQueries

	// register cli commands
	commands := Commands{
		handlers: make(map[string]func(*State, Command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)

	// with program name
	raw_args := os.Args

	if len(raw_args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmd_name := raw_args[1]
	cmd_args := raw_args[2:]
	cmd := Command{
		name: cmd_name,
		args: cmd_args,
	}

	if err := commands.run(&state, cmd); err != nil {
		fmt.Printf("error running command: %v", err)
		os.Exit(1)
	}

}
