package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	apihandler "github.com/venturarome/DaftWatch/internal/handler/api"
)

func InstanceServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      apihandler.InstanceHandler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
