package utils

import "strings"

func Empty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

func NotEmpty(s *string) bool {
	return !Empty(s)
}
