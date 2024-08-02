package main

import (
	"fmt"

	"github.com/venturarome/DaftWatch/internal/bot"
	"github.com/venturarome/DaftWatch/internal/server"
)

func main() {
	bot := bot.InstanceTelegramBot()
	go bot.StartLongPolling()

	server := server.InstanceServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("Server error: %s", err))
	}
}
