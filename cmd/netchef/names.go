package main

import (
	"fmt"
	"math/rand"
	"time"
)

// List of random adjectives and nouns
var adjectives = []string{
	"ancient", "bold", "brilliant", "clever", "curious", "daring", "eager", "fearless",
	"gentle", "lively", "mystic", "noble", "radiant", "swift", "vivid",
}

var nouns = []string{
	"falcon", "pioneer", "voyager", "explorer", "sentinel", "beacon", "harbinger", "trailblazer",
	"visionary", "guardian", "seeker", "oracle", "ranger", "champion", "wanderer", "sage",
}

// Function to generate a random Heroku-style chain name
func randomChainName() string {
	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s", adj, noun)
}
