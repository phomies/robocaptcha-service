package main

import (
	"math/rand"
	"time"
)

func generateRandomWord() string {
	wordList := []string{
		"singular",
		"designer",
		"degenerate",
		"commitment",
		"tradition",
		"conference",
		"potential",
		"producer",
	}
	rand.Seed(time.Now().UTC().UnixNano())
	randomIdx := rand.Intn(len(wordList))
	return wordList[randomIdx]
}

func generateRandomNumber() int {
	rand.Seed(time.Now().UTC().UnixNano())
	randomNumber := 1000 + rand.Intn(9000)
	return randomNumber
}

func addSpaces(s string) string {
	spacedString := string(s[0])
	for idx, c := range s {
		if idx == 0 {
			continue
		}
		spacedString += " "
		spacedString += string(c)
	}
	return spacedString
}
