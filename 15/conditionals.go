package main

import (
	"strings"
)

type conditionalCommand struct {
	command  []string
	operator string
}

func parseConditionalCommands(input string) []conditionalCommand {
	args := strings.Fields(input)
	var cmds []conditionalCommand

	var curCommand []string

	for _, arg := range args {
		if arg == "&&" || arg == "||" {
			cmds = append(cmds, conditionalCommand{
				command:  curCommand,
				operator: arg,
			})
			curCommand = nil
		} else {
			curCommand = append(curCommand, arg)
		}
	}

	cmds = append(cmds, conditionalCommand{
		command:  curCommand,
		operator: "",
	})

	return cmds
}

func executeConditionalCommands(cmds []conditionalCommand) {
	for _, c := range cmds {
		err := identCommand(strings.Join(c.command, " "))

		if c.operator == "||" {
			if err == nil {
				break
			}
		} else if err != nil {
			break
		}
	}
}
