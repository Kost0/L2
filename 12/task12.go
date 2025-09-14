package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type searchOptions struct {
	after   int
	before  int
	around  int
	count   bool
	ignore  bool
	inverse bool
	fix     bool
	number  bool
}

func searchLines(lines []string, substr string, opts searchOptions) []string {
	res := make([]string, 0)

	left := max(opts.before, opts.around)
	right := max(opts.after, opts.around)
	last := -1

	for i, line := range lines {
		if opts.ignore {
			line = strings.ToLower(line)
			substr = strings.ToLower(substr)
		}

		var founded bool

		if opts.inverse {
			if opts.fix {
				founded = substr != line
			} else {
				founded = !strings.Contains(line, substr)
			}
		} else {
			if opts.fix {
				founded = substr == line
			} else {
				founded = strings.Contains(line, substr)
			}
		}

		if founded {
			if left > 0 {
				for j := 0; j < left; j++ {
					if i-left+j <= last {
						continue
					}
					if opts.number {
						res = append(res, strconv.Itoa(i-left+j)+" "+lines[i-left+j])
					} else {
						res = append(res, lines[i-left+j])
					}
				}
			}
			if last < i {
				if opts.number {
					res = append(res, strconv.Itoa(i)+" "+lines[i])
				} else {
					res = append(res, lines[i])
				}
			}
			if right > 0 {
				for j := 1; j <= right && i+j < len(lines); j++ {
					if last >= i+j {
						continue
					}
					if opts.number {
						res = append(res, strconv.Itoa(i+j)+" "+lines[i+j])
					} else {
						res = append(res, lines[i+j])
					}
				}
			}

			last = i + right
		}
	}

	return res
}

func countLines(lines []string, substr string, ignore, inverse, fix bool) int {
	res := 0

	for _, line := range lines {
		if ignore {
			line = strings.ToLower(line)
			substr = strings.ToLower(substr)
		}

		var founded bool

		if inverse {
			if fix {
				founded = substr != line
			} else {
				founded = !strings.Contains(line, substr)
			}
		} else {
			if fix {
				founded = substr == line
			} else {
				founded = strings.Contains(line, substr)
			}
		}

		if founded {
			res++
		}
	}

	return res
}

func readLines(inputFile string) ([]string, error) {
	var reader io.Reader

	if inputFile == "-" || inputFile == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		reader = file
	}

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

func main() {
	after := pflag.IntP("lines after", "A", 0, "get lines after")
	before := pflag.IntP("lines before", "B", 0, "get lines before")
	around := pflag.IntP("lines around", "C", 0, "get lines around")
	count := pflag.BoolP("count", "c", false, "count lines")
	ignore := pflag.BoolP("ignore", "i", false, "ignore register")
	inverse := pflag.BoolP("inversion", "v", false, "non-matching strings")
	fix := pflag.BoolP("fix", "F", false, "exact match of the substring")
	number := pflag.BoolP("number", "n", false, "numbers of lines")

	pflag.Parse()

	var inputFile string
	var substr string

	if pflag.NArg() > 0 {
		substr = pflag.Arg(0)
		inputFile = pflag.Arg(1)
	} else {
		_, err := fmt.Fprintln(os.Stderr, "Not enough arguments")
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}

	lines, err := readLines(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	opts := searchOptions{
		after:   *after,
		before:  *before,
		around:  *around,
		count:   *count,
		ignore:  *ignore,
		inverse: *inverse,
		fix:     *fix,
		number:  *number,
	}

	if opts.count {
		fmt.Println(countLines(lines, substr, *ignore, *inverse, *fix))
		os.Exit(0)
	}

	res := searchLines(lines, substr, opts)

	for _, line := range res {
		fmt.Println(line)
	}
}
