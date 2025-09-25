package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/AungKyawPhyo1142/rss-aggregator/internal"
)

type State struct {
	config *internal.Config
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

	if err := s.config.SetUser(username); err != nil {
		return err
	}

	fmt.Println("user has been set")

	return nil

}

func main() {

	var state State

	cfg, err := internal.Read()
	if err != nil {
		fmt.Println(err)
	}

	state.config = &cfg
	/*
		commands := Commands{
			handlers: map[string]func(*State, Command) error{
				"login": handlerLogin,
			},
		}
	*/

	commands := Commands{
		handlers: make(map[string]func(*State, Command) error),
	}
	commands.register("login", handlerLogin)

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
