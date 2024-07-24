package scraper

import (
	"math/rand"
)

func randomUserAgentGenerator() func() string {
	var userAgents = []string{
		"UserAgent1",
		"UserAgent2",
		"UserAgent3",
		"UserAgent4",
	}
	return func() string {
		return userAgents[rand.Int()%len(userAgents)]
	}
}
