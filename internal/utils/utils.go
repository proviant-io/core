package utils

import "strings"

func ContainsInt(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ClearString(s string) string{
	return strings.TrimSpace(s)
}