package main

import (
	"bufio"
	"fmt"
	"os"
)

var supportedCommands = commandsRegistry{}

func init() {
	addHelp(supportedCommands)
	addMap(supportedCommands)
	addExit(supportedCommands)
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
