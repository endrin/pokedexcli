package main

import (
	"bufio"
	"fmt"
	"os"
)

var supportedCommands = commandsRegistry{}

func init() {
	addHelp(supportedCommands)
	supportedCommands.register(cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback: func() error {
			fmt.Println("Closing the Pokedex... Goodbye!")
			os.Exit(0)

			return nil
		},
	})
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			commandWithArguments := cleanInput(input)
			if len(commandWithArguments) == 0 {
				continue
			}
			keyword := commandWithArguments[0]
			// args := commandWithArguments[1:]
			if err := supportedCommands.executeCommand(keyword); err != nil {
				fmt.Printf("Error executing `%s`: %v\n", keyword, err)
			}
		}
	}
}
