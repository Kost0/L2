package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

type sortOptions struct {
	column         int
	numeric        bool
	reverse        bool
	unique         bool
	month          bool
	ignoreTrailing bool
	checkSorted    bool
	humanReadable  bool
	inputFile      string
}

var humanReadableSuffix = map[string]int64{
	"K": 1 << 10,
	"M": 1 << 20,
	"G": 1 << 30,
	"T": 1 << 40,
	"P": 1 << 50,
	"E": 1 << 60,
}

var monthNames = map[string]time.Month{
	"jan": time.January,
	"feb": time.February,
	"mar": time.March,
	"apr": time.April,
	"may": time.May,
	"jun": time.June,
	"jul": time.July,
	"aug": time.August,
	"sep": time.September,
	"oct": time.October,
	"nov": time.November,
	"dec": time.December,
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

func getColumnValue(line string, column int, ignoreTrailing bool) string {
	if column == 0 {
		if ignoreTrailing {
			return strings.TrimRight(line, " \t")
		}
		return line
	}

	columns := strings.Split(line, " ")
	if column > 0 && column <= len(columns) {
		value := columns[column-1]
		if ignoreTrailing {
			return strings.TrimRight(value, " \t")
		}
		return value
	}

	return ""
}

func parseMonth(s string) (int, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if month, ok := monthNames[s]; ok {
		return int(month), nil
	}
	return 0, fmt.Errorf("invalid month: %s", s)
}

func parseHumanReadable(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))

	if num, err := strconv.ParseInt(s, 10, 64); err == nil {
		return num, nil
	}

	re := regexp.MustCompile(`^(\d+)([KMGTPE]?)B?$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid human readable format: %s", s)
	}

	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	if suffix := matches[2]; suffix != "" {
		multiplier, ok := humanReadableSuffix[suffix]
		if !ok {
			return 0, fmt.Errorf("unknow suffix: %s", suffix)
		}
		num *= multiplier
	}

	return num, nil
}

func compareValues(first, second string, opts sortOptions) (bool, error) {
	if opts.ignoreTrailing {
		first = strings.TrimRight(first, " \t")
		second = strings.TrimRight(second, " \t")
	}

	if opts.numeric {
		numA, errA := strconv.ParseFloat(first, 64)
		numB, errB := strconv.ParseFloat(second, 64)

		if errA == nil && errB == nil {
			return numA < numB, nil
		} else if errA == nil {
			return true, errA
		} else if errB == nil {
			return false, errB
		}
	}

	if opts.month {
		monthA, errA := parseMonth(first)
		monthB, errB := parseMonth(second)

		if errA == nil && errB == nil {
			return monthA < monthB, nil
		} else if errA == nil {
			return true, errA
		} else if errB == nil {
			return false, errB
		}
	}

	if opts.humanReadable {
		sizeA, errA := parseHumanReadable(first)
		sizeB, errB := parseHumanReadable(second)

		if errA == nil && errB == nil {
			return sizeA < sizeB, nil
		} else if errA == nil {
			return true, errA
		} else if errB == nil {
			return false, errB
		}
	}

	return first < second, nil
}

func isSorted(lines []string, opts sortOptions) bool {
	for i := 1; i < len(lines); i++ {
		prevVal := getColumnValue(lines[i-1], opts.column, opts.ignoreTrailing)
		curVal := getColumnValue(lines[i], opts.column, opts.ignoreTrailing)

		less, err := compareValues(prevVal, curVal, opts)
		if err != nil {
			continue
		}

		if opts.reverse {
			if less {
				return false
			}
		} else {
			if !less && prevVal != curVal {
				return false
			}
		}
	}

	return true
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return result
}

func sortLines(lines []string, opts sortOptions) {
	sort.SliceStable(lines, func(i, j int) bool {
		valA := getColumnValue(lines[i], opts.column, opts.ignoreTrailing)
		valB := getColumnValue(lines[j], opts.column, opts.ignoreTrailing)

		less, err := compareValues(valA, valB, opts)
		if err != nil {
			return lines[i] < lines[j]
		}

		if opts.reverse {
			return !less
		}

		return less
	})
}

func main() {
	column := pflag.IntP("key", "k", 0, "sort by column")
	numeric := pflag.BoolP("numbers", "n", false, "sort by numeric")
	reverse := pflag.BoolP("reverse", "r", false, "reverse sorting")
	unique := pflag.BoolP("unique", "u", false, "only unique")
	month := pflag.BoolP("month", "M", false, "sort by month")
	ignoreTrailing := pflag.BoolP("ignore blanks", "b", false, "ignore trailing blanks")
	checkSorted := pflag.BoolP("check", "c", false, "check is it sorted")
	humanReadable := pflag.BoolP("human readable", "h", false, "sort by human-readable sizes")

	pflag.Parse()

	var inputFile string

	if pflag.NArg() > 0 {
		inputFile = pflag.Arg(0)
	}

	lines, err := readLines(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	opts := sortOptions{
		column:         *column,
		numeric:        *numeric,
		reverse:        *reverse,
		unique:         *unique,
		month:          *month,
		ignoreTrailing: *ignoreTrailing,
		checkSorted:    *checkSorted,
		humanReadable:  *humanReadable,
		inputFile:      inputFile,
	}

	if opts.checkSorted {
		if isSorted(lines, opts) {
			os.Exit(0)
		} else {
			_, err = fmt.Fprintf(os.Stderr, "Input is not sorted\n")
			os.Exit(1)
		}
	}

	if opts.unique {
		lines = removeDuplicates(lines)
	}

	sortLines(lines, opts)
	for _, line := range lines {
		fmt.Println(line)
	}
}
