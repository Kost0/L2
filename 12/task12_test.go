package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountLines(t *testing.T) {
	tests := []struct {
		lines   []string
		substr  string
		ignore  bool
		inverse bool
		fix     bool
		want    int
	}{
		{[]string{"ab", "cab", "abc", "bb"}, "ab", false, false, false, 3},
		{[]string{"ad", "Ab", "bc", "bb"}, "ab", false, false, false, 0},
		{[]string{"ab", "Ab", "aBc", "bb"}, "ab", true, false, false, 3},
		{[]string{"ab", "Ab", "aBc", "bb"}, "ab", false, true, false, 3},
		{[]string{"ab", "Ab", "aBc", "bb"}, "ab", true, true, false, 1},
		{[]string{"ab", "Ab", "aBc", "abc"}, "ab", false, false, true, 1},
		{[]string{"ab", "Ab", "aBc", "abc"}, "ab", false, true, true, 3},
		{[]string{"ab", "Ab", "aBc", "abc"}, "ab", true, false, true, 2},
		{[]string{"ab", "Ab", "aBc", "abc"}, "ab", true, true, true, 2},
	}

	for _, test := range tests {
		res := countLines(test.lines, test.substr, test.ignore, test.inverse, test.fix)
		assert.Equal(t, test.want, res)
	}
}

func TestSearchLines(t *testing.T) {
	tests := []struct {
		lines  []string
		substr string
		opts   searchOptions
		want   []string
	}{
		{[]string{"ab", "cab", "abc", "Ab"}, "ab", searchOptions{}, []string{"ab", "cab", "abc"}},
		{[]string{"ab", "cab", "abc", "Ab"}, "ab", searchOptions{ignore: true}, []string{"ab", "cab", "abc", "Ab"}},
		{[]string{"ab", "cab", "abc", "Ab"}, "ab", searchOptions{fix: true}, []string{"ab"}},
		{[]string{"ab", "cab", "abc", "bg"}, "ab", searchOptions{inverse: true}, []string{"bg"}},
		{[]string{"ab", "cab", "abc", "bg"}, "ab", searchOptions{number: true}, []string{"0 ab", "1 cab", "2 abc"}},
		{[]string{"ad", "cab", "adc", "bg"}, "ab", searchOptions{around: 1}, []string{"ad", "cab", "adc"}},
		{[]string{"ad", "cab", "adc", "bg"}, "ab", searchOptions{before: 1}, []string{"ad", "cab"}},
		{[]string{"ad", "cab", "adc", "bg"}, "ab", searchOptions{after: 1}, []string{"cab", "adc"}},
		{[]string{"ad", "cab", "adc", "bg"}, "ab", searchOptions{around: 1, after: 2}, []string{"ad", "cab", "adc", "bg"}},
		{[]string{"ad", "cab", "adc", "bg"}, "ab", searchOptions{around: 1, after: 2, number: true}, []string{"0 ad", "1 cab", "2 adc", "3 bg"}},
		{[]string{"ab", "cab", "adc", "bg"}, "ab", searchOptions{around: 1, after: 2, number: true, fix: true}, []string{"0 ab", "1 cab", "2 adc"}},
		{[]string{"ab", "cab", "adc", "bg"}, "ab", searchOptions{around: 1, after: 2, number: true, fix: true, inverse: true}, []string{"0 ab", "1 cab", "2 adc", "3 bg"}},
		{[]string{"ab", "cab", "adc", "bg"}, "ab", searchOptions{before: 0, around: 1, after: 2, number: true, fix: true, inverse: true}, []string{"0 ab", "1 cab", "2 adc", "3 bg"}},
	}

	for _, test := range tests {
		res := searchLines(test.lines, test.substr, test.opts)
		assert.Equal(t, test.want, res)
	}
}
