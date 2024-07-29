package main

import (
	"fmt"

	"github.com/venturarome/DaftWatch/internal/bot"
	"github.com/venturarome/DaftWatch/internal/server"
)

func main() {

	telegramBot := bot.InitTelegramBot()
	go bot.StartLongPolling(telegramBot)

	server := server.InstanceServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("Server error: %s", err))
	}
}
