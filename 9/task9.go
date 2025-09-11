package main

import (
	"errors"
	"unicode"
)

func unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	res := ""
	cur := ""
	num := 0
	isStr := false
	for i := 0; i < len(str)-1; i++ {
		if unicode.IsLetter(rune(str[i])) {
			if isStr {
				cur += "/"
			}
			cur += string(str[i])

		} else if isStr {
			res += cur
			cur = string(str[i])
			isStr = false
		} else if unicode.IsDigit(rune(str[i])) {
			num *= 10
			num += int(str[i]) - 48
			if !unicode.IsDigit(rune(str[i+1])) {
				for n := 0; n < num; n++ {
					res += cur
				}
				cur = ""
				num = 0
			}
		} else if str[i] == '/' {
			isStr = true
		} else {
			return "", errors.New("wrong char")
		}
	}

	if unicode.IsDigit(rune(str[len(str)-1])) {
		if isStr {
			cur += string(str[len(str)-1])
			res += cur
		} else if cur == "" {
			return "", errors.New("wrong string")
		} else {
			num *= 10
			num += int(str[len(str)-1]) - 48
			for n := 0; n < num; n++ {
				res += cur
			}
		}
	} else {
		cur += string(str[len(str)-1])
		res += cur
	}

	return res, nil
}
