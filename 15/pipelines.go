package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

func parseCommands(input string) [][]string {
	parts := strings.Split(input, "|")
	var cmds [][]string

	for _, part := range parts {
		args := strings.Fields(strings.TrimSpace(part))
		if len(args) == 0 {
			continue
		}
		cmds = append(cmds, args)
	}

	return cmds
}

func executePipeline(commands [][]string) error {
	if len(commands) == 0 {
		return errors.New("no commands in pipeline")
	}

	var pipes []io.ReadCloser
	var cmds []*exec.Cmd

	for i, args := range commands {
		cleanArgs, redirect := parseRedirection(args)

		var cmd *exec.Cmd
		if len(args) > 1 {
			cmd = exec.Command(cleanArgs[0], cleanArgs[1:]...)
		} else if len(args) == 1 {
			cmd = exec.Command(cleanArgs[0])
		} else {
			return errors.New("empty command")
		}

		err := applyRedirection(cmd, redirect)
		if err != nil {
			return err
		}

		if i == 0 {
			if cmd.Stdin == nil {
				cmd.Stdin = os.Stdin

			}
		} else {
			cmd.Stdin = pipes[i-1]
		}

		if i == len(commands)-1 {
			if cmd.Stdout == nil {
				cmd.Stdout = os.Stdout
			}
			if cmd.Stderr == nil {
				cmd.Stderr = os.Stderr
			}
		} else {
			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			pipes = append(pipes, stdoutPipe)
		}

		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				if exitError.ExitCode() == 1 && strings.Contains(cmd.Path, "grep") {
					continue
				}
				return err
			}
			return err
		}
	}

	return nil
}
