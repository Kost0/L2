package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func identCommand(input string) error {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	command := args[0]

	var err error

	switch command {
	case "cd":

		err = cd(args[1:])
	case "pwd":
		err = pwd(args[1:])
	case "echo":
		err = echo(args[1:])
	case "kill":
		err = kill(args[1:])
	case "ps":
		err = ps(args[1:])
	default:
		err = externalCommand(args)
	}

	return err
}

func cd(args []string) error {
	var path string

	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = home
	} else {
		path = args[0]
	}
	if path == "~" {
		path = os.Getenv("HOME")
	}
	err := os.Chdir(path)
	if err != nil {
		return err
	}

	return nil
}

func pwd(args []string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	_, redirect := parseRedirection(args)
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
		defer func() {
			err = file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		_, err = file.WriteString(path)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(path)
	}

	return nil
}

func echo(args []string) error {
	newArgs, redirect := parseRedirection(args)
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

		defer func() {
			err = file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		for _, arg := range newArgs {
			_, err = file.WriteString(arg + " ")
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	} else {
		for _, arg := range newArgs {
			fmt.Print(arg + " ")
		}
		fmt.Print("\n")
	}

	return nil
}

func kill(args []string) error {
	if len(args) == 0 {
		return errors.New("the argument is missing")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}

func ps(args []string) error {
	args, redirect := parseRedirection(args)
	cmd := exec.Command("ps")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := applyRedirection(cmd, redirect)
	if err != nil {
		return err
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func externalCommand(args []string) error {
	newArgs, redirect := parseRedirection(args)
	var cmd *exec.Cmd

	if len(newArgs) > 1 {
		cmd = exec.Command(newArgs[0], newArgs[1:]...)
	} else {
		cmd = exec.Command(newArgs[0])
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := applyRedirection(cmd, redirect)
	if err != nil {
		return err
	}

	err = cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			fmt.Println(string(exitError.Stderr))
		} else {
			fmt.Printf("%s: command not found...\n", args[0])
		}
		return err
	}

	return nil
}
