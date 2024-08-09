package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (th *TelegramHandler) HandleDonate(update tgbotapi.Update) (msg tgbotapi.MessageConfig, clearContext bool) {
	if update.Message == nil { // panic on non-Message updates
		panic("Received non-Message Update")
	}

	messageText := strings.Trim(update.Message.Text, " ")
	if strings.Index(messageText, "/donate") != 0 {
		panic("Badly routed Message Update")
	}

	fmt.Println("[DEBUG] donate > received: ", messageText)

	// TODO research in the future: https://core.telegram.org/bots/payments
	// TODO consider adding Bitcoin/crypto, too.

	chatId := update.Message.Chat.ID
	msgText := "<b>You can donate via Revolut:</b> revolut.me/venturamendo\n\n" +
		"This bot was initially developed for personal use. " +
		"Source code can be found on GitHub: https://github.com/venturarome/DaftWatch.\n\n" +
		"Once working, it was decided to make it accessible for everyone, at no charge. " +
		"However, if you feel like making a donation to the bot developer, it will be greatly appreciated. " +
		"Your support helps keep this bot running smoothly and allows for continuous improvements and new features.\n\n" +
		"Plus, you can reach out the bot developer with the /feedback or /contact commands. " +
		"As a sign of gratitude for the donation, he will remove some bot limitations (ie. nomber of simultaneous alerts)," +
		"initially set to ensure a fair use of the tool.\n\n" +
		"Thank you for your support! 🙏"

	msg = tgbotapi.NewMessage(chatId, msgText)
	msg.ParseMode = "html"
	msg.DisableWebPagePreview = true

	return msg, true
}
