package main

import (
	"bufio"
	"fmt"
	"os"
	gptintegration "pcx-translate/gpt"
)

func main() {

	option := getInput(`
	Portuguese and Spanish - enter 1
	`)

	filePath := getInput("Inform the filepath to the docs you want to check: ")
	err := gptintegration.GptIntegration(option, filePath)
	if err != nil {
		panic(err)
	}

}

// Get the input informed by the user
func getInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
