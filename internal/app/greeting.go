package app

import "strings"

func Greeting(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		trimmed = "world"
	}

	return "hello, " + trimmed
}
