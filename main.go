package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", hello)
	router.GET("/list_rand_uas", listRandUas)

	router.Run("localhost:8080")
}

func hello(c *gin.Context) {
	var salutation = []string{"hello", "world", "!"}
	c.IndentedJSON(http.StatusOK, salutation)
}

func listRandUas(c *gin.Context) {
	uas := make([]string, 5)
	gen := randomUserAgentGenerator()

	for i := 0; i < 5; i++ {
		uas[i] = gen()
	}
	c.IndentedJSON(http.StatusOK, uas)
}

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
