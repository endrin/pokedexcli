package main

import (
	"fmt"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type commandsRegistry map[string]cliCommand

func (registry commandsRegistry) register(command cliCommand) {
	registry[command.name] = command
}

func (registry commandsRegistry) executeCommand(keyword string) error {
	if command, ok := registry[keyword]; ok {
		return command.callback()
	}
	return fmt.Errorf("unknown command")
}

func addHelp(registry commandsRegistry) {
	registry.register(cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func() error {
			fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

			for _, command := range registry {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}

			fmt.Println()

			return nil
		},
	})
}
