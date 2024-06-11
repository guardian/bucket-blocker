package main

import "strings"

func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	split := strings.Split(s, ",")
	for i, v := range split {
		split[i] = strings.TrimSpace(v)
	}
	return split
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
