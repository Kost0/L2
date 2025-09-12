package main

import (
	"fmt"
	"sort"
	"strings"
)

func findAnagrams(strs []string) map[string][]string {
	m1 := map[string][]string{}
	m2 := map[string][]string{}
	for _, str := range strs {
		runes := []rune(strings.ToLower(str))
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		strNew := string(runes)

		if _, ok := m1[strNew]; ok {
			m1[strNew] = append(m1[strNew], str)
		} else {
			m1[strNew] = []string{str}
		}
	}

	for _, v := range m1 {
		if len(v) > 1 {
			m2[v[0]] = append(m2[v[0]], v...)
		}
	}

	return m2
}

func main() {
	strs := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	m := findAnagrams(strs)
	fmt.Println(m)
}
