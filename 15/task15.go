package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	signal.Ignore(os.Interrupt)

	reader := bufio.NewReader(os.Stdin)

	directory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}

		input = strings.TrimSuffix(input, "\n")
		if input == "" {
			continue
		}

		if strings.Contains(input, "$") {
			input = expandEnv(input)
		}

		if strings.Contains(input, "||") || strings.Contains(input, "&&") {
			cmds := parseConditionalCommands(input)
			executeConditionalCommands(cmds)
		} else if strings.Contains(input, "|") {
			cmds := parseCommands(input)
			err = executePipeline(cmds)
			if err != nil {
				log.Println(err)
			}
		} else {
			err = identCommand(input)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
