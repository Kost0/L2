package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFields(t *testing.T) {
	tests := []struct {
		input   string
		output  []int
		wantErr bool
	}{
		{"1", []int{0}, false},
		{"1,3,5", []int{0, 2, 4}, false},
		{"1-5", []int{0, 1, 2, 3, 4}, false},
		{"1-3,5", []int{0, 1, 2, 4}, false},
		{"1,2-3,5", []int{0, 1, 2, 4}, false},
		{"", []int(nil), true},
		{"a", []int(nil), true},
	}

	for _, test := range tests {
		res, err := parseFields(test.input)
		assert.Equal(t, test.output, res)
		assert.Equal(t, test.wantErr, err != nil)
	}
}

func TestCutLines(t *testing.T) {
	tests := []struct {
		input  []string
		opts   *cutOptions
		output [][]string
	}{
		{[]string{"ab\tcd\tef\tgh"}, &cutOptions{fields: []int{0, 1}, delimiter: "\t"}, [][]string{{"ab", "cd"}}},
		{[]string{"ab,cd,ef,gh"}, &cutOptions{fields: []int{0, 1, 2}, delimiter: ","}, [][]string{{"ab", "cd", "ef"}}},
		{[]string{"ab,cd,ef,gh", "ab,cd"}, &cutOptions{fields: []int{0, 1, 2}, delimiter: ","}, [][]string{{"ab", "cd", "ef"}, {"ab", "cd"}}},
		{[]string{"ab,cd,ef,gh", "ab"}, &cutOptions{fields: []int{0, 1, 2}, delimiter: ","}, [][]string{{"ab", "cd", "ef"}, {"ab"}}},
		{[]string{"ab,cd,ef,gh", "ab"}, &cutOptions{fields: []int{0, 1, 2}, delimiter: ",", separated: true}, [][]string{{"ab", "cd", "ef"}}},
	}

	for _, test := range tests {
		res := cutLines(test.input, test.opts)
		assert.Equal(t, test.output, res)
	}
}
