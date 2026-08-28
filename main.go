package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/lulzshadowwalker/pokedex/repl"
)

func main() {
	registry := registry{commands: make(map[string]command)}
	commands := []command{
		{
			name:        "exit",
			description: "Exit the Pokdex",
			callback: func() error {
				fmt.Printf("Closing the Pokedex... Goodbye!")
				os.Exit(0)
				return nil
			},
		},

		{
			name:        "help",
			description: "Displays a help message",
			callback: func() error {
				fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
				for _, command := range registry.commands {
					fmt.Printf("%s: %s\n", command.name, command.description)
				}

				return nil
			},
		},
	}

	for _, command := range commands {
		registry.commands[command.name] = command
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokdex > ")
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Printf("failed to read from stdin: %q", err)
				os.Exit(1)
			}

			os.Exit(0)
		}
		splits := repl.Split(scanner.Text())
		if len(splits) == 0 {
			continue
		}

		if err := registry.run(splits[0]); err != nil {
			fmt.Println(err.Error())
		}
	}
}

type command struct {
	name        string
	description string
	callback    func() error
}

type registry struct {
	commands map[string]command
}

var ErrCommandUnknown = errors.New("Unknown command")

func (r *registry) run(name string) error {
	command, ok := r.commands[name]
	if !ok {
		return ErrCommandUnknown
	}

	return command.callback()
}
