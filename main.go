package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/lulzshadowwalker/pokedex/cache"
	"github.com/lulzshadowwalker/pokedex/repl"
)

func main() {
	registry := registry{
		config: &config{
			next:  new("https://pokeapi.co/api/v2/location-area"),
			cache: cache.New(5 * time.Second),
		},
		commands: make(map[string]command),
	}
	commands := []command{
		{
			name:        "exit",
			description: "Exit the Pokdex",
			callback: func(config *config, args []string) error {
				fmt.Printf("Closing the Pokedex... Goodbye!")
				os.Exit(0)
				return nil
			},
		},

		{
			name:        "help",
			description: "Displays a help message",
			callback: func(config *config, args []string) error {
				fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
				for _, command := range registry.commands {
					fmt.Printf("%s: %s\n", command.name, command.description)
				}

				return nil
			},
		},

		{
			name:        "map",
			description: "Displays a paginated list of areas in the Pokemon world",
			callback: func(config *config, args []string) error {
				if config.next == nil {
					return nil
				}

				url := *config.next

				resolver := func() []byte {
					response, err := http.Get(url)
					if err != nil {
						return nil
					}
					defer response.Body.Close()

					if response.StatusCode != http.StatusOK {
						return nil
					}

					bytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil
					}

					return bytes
				}

				bytes := config.cache.Remember(url, resolver, 5*time.Minute)
				if bytes == nil {
					return ErrUnknown
				}

				data := struct {
					Previous *string
					Next     *string
					Results  []struct {
						Name string
					}
				}{}
				if err := json.Unmarshal(bytes, &data); err != nil {
					return err
				}

				config.previous = data.Previous
				config.next = data.Next

				for _, result := range data.Results {
					fmt.Println(result.Name)
				}

				return nil
			},
		},

		{
			name:        "mapb",
			description: "Displays a paginated list of areas in the Pokemon world (map backward)",
			callback: func(config *config, args []string) error {
				if config.previous == nil {
					return nil
				}

				url := *config.previous

				resolver := func() []byte {
					response, err := http.Get(url)
					if err != nil {
						return nil
					}
					defer response.Body.Close()

					if response.StatusCode != http.StatusOK {
						return nil
					}

					bytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil
					}

					return bytes
				}

				bytes := config.cache.Remember(url, resolver, 5*time.Minute)
				if bytes == nil {
					return ErrUnknown
				}

				data := struct {
					Previous *string
					Next     *string
					Results  []struct {
						Name string
					}
				}{}
				if err := json.Unmarshal(bytes, &data); err != nil {
					return err
				}

				config.previous = data.Previous
				config.next = data.Next

				for _, result := range data.Results {
					fmt.Println(result.Name)
				}

				return nil
			},
		},

		{
			name:        "explore",
			description: "Explor an area to list potential Pokemons you may encounter",
			callback: func(config *config, args []string) error {
				if len(args) == 0 {
					return ErrInvalidArgument
				}

				fmt.Printf("Exploring %s...\n", args[0])

				url := "https://pokeapi.co/api/v2/location-area/" + args[0]

				resolver := func() []byte {
					response, err := http.Get(url)
					if err != nil {
						return nil
					}
					defer response.Body.Close()

					if response.StatusCode != http.StatusOK {
						return nil
					}

					bytes, err := io.ReadAll(response.Body)
					if err != nil {
						return nil
					}

					return bytes
				}

				bytes := config.cache.Remember(url, resolver, 5*time.Minute)
				if bytes == nil {
					return ErrUnknown
				}

				data := struct {
					PokemonEncounters []struct {
						Pokemon struct {
							Name string `json:"name"`
						} `json:"pokemon"`
					} `json:"pokemon_encounters"`
				}{}
				if err := json.Unmarshal(bytes, &data); err != nil {
					return err
				}

				fmt.Println("Found Pokemon:")
				for _, result := range data.PokemonEncounters {
					fmt.Println("-", result.Pokemon.Name)
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
		fmt.Print("Pokedex > ")
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

		if err := registry.run(splits[0], splits[1:]); err != nil {
			fmt.Println(err.Error())
		}
	}
}

type config struct {
	next     *string
	previous *string
	cache    *cache.Cache
}

type command struct {
	name        string
	description string
	callback    func(config *config, args []string) error
}

type registry struct {
	config   *config
	commands map[string]command
}

var (
	ErrCommandUnknown = errors.New("Unknown command")
	ErrUnknown        = errors.New("Unknown error has occurred")
	ErrInvalidArgument = errors.New("Invalid argument")
)

func (r *registry) run(name string, args []string) error {
	command, ok := r.commands[name]
	if !ok {
		return ErrCommandUnknown
	}

	return command.callback(r.config, args)
}
