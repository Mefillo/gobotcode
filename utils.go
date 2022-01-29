package main

import (
	"fmt"
	"strconv"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func convertToIndexed(list []string) (result string) {
	var index string
	for i, b := range list {
		index = strconv.Itoa(i)
		if i != 0 {
			result = result + "\n"
		}
		result = fmt.Sprintf("%s[%s] %s", result, index, b)

	}
	return
}
