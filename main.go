package main

import (
	"bufio"
	"fmt"
	"os"
)

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
			command := commandWithArguments[0]
			// args := commandWithArguments[1:]
			fmt.Printf("Your command was: %s\n", command)
		}
	}
}
