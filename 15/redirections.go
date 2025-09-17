package main

import (
	"os"
	"os/exec"
)

type redirection struct {
	Stdout       string
	Stderr       string
	Stdin        string
	StdoutAppend bool
	StderrAppend bool
}

func parseRedirection(input []string) ([]string, *redirection) {
	redirect := &redirection{}
	cleanArgs := make([]string, 0)

	for i := 0; i < len(input); i++ {
		arg := input[i]

		switch arg {
		case ">":
			if i < len(input)-1 {
				redirect.Stdout = input[i+1]
				redirect.StdoutAppend = false
				i++
			}
		case ">>":
			if i < len(input)-1 {
				redirect.Stdout = input[i+1]
				redirect.StdoutAppend = true
				i++
			}
		case "2>":
			if i < len(input)-1 {
				redirect.Stderr = input[i+1]
				redirect.StderrAppend = false
				i++
			}
		case "2>>":
			if i < len(input)-1 {
				redirect.Stderr = input[i+1]
				redirect.StderrAppend = true
				i++
			}
		case "&>":
			if i < len(input)-1 {
				redirect.Stdout = input[i+1]
				redirect.Stderr = input[i+1]
				redirect.StdoutAppend = false
				redirect.StderrAppend = false
				i++
			}
		case "<":
			if i < len(input)-1 {
				redirect.Stdin = input[i+1]
				i++
			}
		default:
			if !isRedirectOperator(arg) {
				cleanArgs = append(cleanArgs, input[i])
			}
		}
	}

	return cleanArgs, redirect
}

func applyRedirection(cmd *exec.Cmd, redirect *redirection) error {
	if redirect.Stdout != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if redirect.StdoutAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err := os.OpenFile(redirect.Stdout, flags, 0644)
		if err != nil {
			return err
		}
		cmd.Stdout = file
	}

	if redirect.Stderr != "" && redirect.Stderr != redirect.Stdout {
		flags := os.O_CREATE | os.O_WRONLY
		if redirect.StderrAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err := os.OpenFile(redirect.Stderr, flags, 0644)
		if err != nil {
			return err
		}
		cmd.Stderr = file
	}

	if redirect.Stdin != "" {
		file, err := os.Open(redirect.Stdin)
		if err != nil {
			return err
		}
		cmd.Stdin = file
	}

	return nil
}

func isRedirectOperator(arg string) bool {
	operators := []string{">", ">>", "2>", "2>>", "&>", "<"}
	for _, operator := range operators {
		if arg == operator {
			return true
		}
	}
	return false
}
