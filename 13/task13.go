package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/pflag"
)

type cutOptions struct {
	fields    []int
	delimiter string
	separated bool
}

func readLines() ([]string, error) {
	reader := os.Stdin

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseFields(fields string) ([]int, error) {
	res := make([]int, 0)
	prevNum := ""
	curNum := ""
	for i := range fields {
		if unicode.IsDigit(rune(fields[i])) {
			curNum += string(fields[i])
		} else if fields[i] == ',' {
			if prevNum == "" {
				digit, err := strconv.Atoi(curNum)
				if err != nil {
					return nil, err
				}
				res = append(res, digit-1)
			} else {
				prevDigit, err := strconv.Atoi(prevNum)
				if err != nil {
					return nil, err
				}
				lastDigit, err := strconv.Atoi(curNum)
				if err != nil {
					return nil, err
				}

				for j := prevDigit; j <= lastDigit; j++ {
					res = append(res, j-1)
				}
			}
			prevNum = ""
			curNum = ""
		} else if fields[i] == '-' {
			prevNum = curNum
			curNum = ""
		} else {
			return nil, errors.New("invalid field")
		}
	}

	if prevNum == "" {
		digit, err := strconv.Atoi(curNum)
		if err != nil {
			return nil, err
		}
		res = append(res, digit-1)
	} else {
		prevDigit, err := strconv.Atoi(prevNum)
		if err != nil {
			return nil, err
		}
		lastDigit, err := strconv.Atoi(curNum)
		if err != nil {
			return nil, err
		}

		for j := prevDigit; j <= lastDigit; j++ {
			res = append(res, j-1)
		}
	}

	if len(res) == 0 {
		return nil, errors.New("need fields")
	}

	return res, nil
}

func cutLines(lines []string, opts *cutOptions) [][]string {
	res := make([][]string, 0)

	for _, line := range lines {
		if strings.Contains(line, opts.delimiter) {
			parts := strings.Split(line, opts.delimiter)
			resParts := make([]string, 0)

			for _, field := range opts.fields {
				if field >= len(parts) {
					break
				}
				resParts = append(resParts, parts[field])
			}

			res = append(res, resParts)
		} else {
			if !opts.separated {
				res = append(res, []string{line})
			}
		}
	}

	return res
}

func main() {
	fields := pflag.StringP("fields", "f", "", "specify fields to output")
	delimiter := pflag.StringP("delimiter", "d", "\t", "symbol for separation")
	separated := pflag.BoolP("separated", "s", false, "only separated")

	pflag.Parse()

	numFields, err := parseFields(*fields)
	if err != nil {
		log.Fatal(err)
	}

	lines, err := readLines()
	if err != nil {
		log.Fatal(err)
	}

	opts := cutOptions{
		fields:    numFields,
		delimiter: *delimiter,
		separated: *separated,
	}

	res := cutLines(lines, &opts)

	for _, line := range res {
		for i := range len(line) - 1 {
			fmt.Print(line[i] + opts.delimiter)
		}
		fmt.Print(line[len(line)-1] + "\n")
	}
}
