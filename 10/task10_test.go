package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHumanReadable(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{"100", 100, false},
		{"1K", 1024, false},
		{"1M", 1048576, false},
		{"1G", 1073741824, false},
		{"a", 0, true},
		{"1k", 1024, false},
	}

	for _, test := range tests {
		result, err := parseHumanReadable(test.input)
		assert.Equal(t, result, test.expected)
		assert.Equal(t, test.wantErr, err != nil)
	}
}

func TestParseMonth(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"jan", 1, false},
		{"Feb", 2, false},
		{"mAR", 3, false},
		{"january", 0, true},
		{"12", 0, true},
	}

	for _, test := range tests {
		result, err := parseMonth(test.input)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, test.wantErr, err != nil)
	}
}

func TestGetColumnValue(t *testing.T) {
	tests := []struct {
		input    string
		column   int
		expected string
	}{
		{"ab c d", 1, "ab"},
		{"ab c d", 2, "c"},
		{"ab c d", 3, "d"},
		{"ab c d", 4, ""},
	}

	for _, test := range tests {
		result := getColumnValue(test.input, test.column, false)
		assert.Equal(t, test.expected, result)
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		input1   string
		input2   string
		opts     sortOptions
		expected bool
		wantErr  bool
	}{
		{"a", "b", sortOptions{}, true, false},
		{"b", "a", sortOptions{}, false, false},
		{"1", "2", sortOptions{numeric: true}, true, false},
		{"1K", "1G", sortOptions{humanReadable: true}, true, false},
		{"feb", "JaN", sortOptions{month: true}, false, false},
		{"1", "jan", sortOptions{month: true}, false, false},
		{"jan", "1", sortOptions{month: true}, true, false},
		{"a", "a ", sortOptions{ignoreTrailing: true}, false, false},
		{"a", "a ", sortOptions{}, true, false},
	}

	for _, test := range tests {
		result, err := compareValues(test.input1, test.input2, test.opts)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, test.wantErr, err != nil)
	}
}

func TestIsSorted(t *testing.T) {
	tests := []struct {
		input    []string
		opts     sortOptions
		expected bool
	}{
		{[]string{"a", "b"}, sortOptions{}, true},
		{[]string{"2", "1"}, sortOptions{numeric: true}, false},
		{[]string{"jan", "feb"}, sortOptions{month: true}, true},
		{[]string{"1G", "1M"}, sortOptions{humanReadable: true}, false},
	}

	for _, test := range tests {
		result := isSorted(test.input, test.opts)
		assert.Equal(t, test.expected, result)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"a", "b", "b"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"1", "2", "2"}, []string{"1", "2"}},
	}

	for _, test := range tests {
		result := removeDuplicates(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestSortLines(t *testing.T) {
	tests := []struct {
		input    []string
		opts     sortOptions
		expected []string
	}{
		{[]string{"a", "c", "b"}, sortOptions{}, []string{"a", "b", "c"}},
		{[]string{"a", "a", "a"}, sortOptions{}, []string{"a", "a", "a"}},
		{[]string{"a", "c", "b"}, sortOptions{reverse: true}, []string{"c", "b", "a"}},
		{[]string{"2", "1", "3"}, sortOptions{numeric: true}, []string{"1", "2", "3"}},
		{[]string{"2", "1", "3"}, sortOptions{numeric: true, reverse: true}, []string{"3", "2", "1"}},
		{[]string{"mar", "jan", "dec"}, sortOptions{month: true}, []string{"jan", "mar", "dec"}},
		{[]string{"mar", "jan", "dec"}, sortOptions{month: true, reverse: true}, []string{"dec", "mar", "jan"}},
		{[]string{"1M", "1G", "1K"}, sortOptions{humanReadable: true}, []string{"1K", "1M", "1G"}},
		{[]string{"1M", "1G", "1K"}, sortOptions{humanReadable: true, reverse: true}, []string{"1G", "1M", "1K"}},
		{[]string{"a b", "b c", "c a"}, sortOptions{column: 2}, []string{"c a", "a b", "b c"}},
		{[]string{"a b", "b c", "c a"}, sortOptions{column: 2, reverse: true}, []string{"b c", "a b", "c a"}},
		{[]string{"a mar", "b dec", "c jan"}, sortOptions{column: 2, month: true, reverse: true}, []string{"b dec", "a mar", "c jan"}},
	}

	for _, test := range tests {
		sortLines(test.input, test.opts)
		assert.Equal(t, test.expected, test.input)
	}
}
