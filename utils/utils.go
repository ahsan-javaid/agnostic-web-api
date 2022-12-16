package utils

import "strings"

func Check(e error) {
	if e != nil {
			panic(e)
	}
}

func GetCollectionName(url string) string {
	return strings.Split(url, "/")[1]
}

func GetParams(url string) []string {
	return strings.Split(url, "/")
}