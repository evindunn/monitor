package main

import "regexp"

var reKeyValDelim = regexp.MustCompile(`\s*:\s*`)
var reMultiSpace = regexp.MustCompile(`\s+`)

func strSliceContains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

func sum(ints ...int) int {
	sumVal := 0
	for _, val := range ints {
		sumVal += val
	}
	return sumVal
}
