package main

import (
	"fmt"

	"github.com/venturarome/DaftWatch/internal/bot"
	"github.com/venturarome/DaftWatch/internal/server"
)

func main() {

	telegramBot := bot.InitTelegramBot()
	go bot.StartLongPolling(telegramBot)

	server := server.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
