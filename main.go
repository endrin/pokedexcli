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
			keyword := commandWithArguments[0]
			// args := commandWithArguments[1:]
			if err := executeCommand(keyword); err != nil {
				fmt.Printf("Error executing `%s`: %v\n", keyword, err)
			}
		}
	}
}
