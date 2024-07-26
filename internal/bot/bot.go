package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	telegramhandler "github.com/venturarome/DaftWatch/internal/handler/telegram"
)

// type Bot struct {
// 	bot *tgbotapi.BotAPI
// }

func InitTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = os.Getenv("TELEGRAM_MODE") != "release"

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

var userContext map[int64]string = make(map[int64]string)
var userTimestampContext map[int64]int = make(map[int64]int)

const CONTEXT_TTL int = 600 // 10 minutes

func StartLongPolling(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 50 // 50 sec is, apparently, maximum allowed

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		userId := update.Message.From.ID
		currentTimestamp := update.Message.Date
		currentText := strings.Trim(update.Message.Text, " ")

		// DEBUG
		fmt.Println("[DEBUG] Text received: ", currentText)

		if !update.Message.IsCommand() {
			// Probably it comes from a multistep command. Let's check out "in-memory database".
			storedCommand, found := userContext[userId]
			lastInteractionTimestamp := userTimestampContext[userId]
			if !found {
				// No user context.
				update.Message.Text = "/nocommand"
			} else if currentTimestamp-lastInteractionTimestamp > CONTEXT_TTL {
				// Outdated user context.
				delete(userContext, userId)
				delete(userTimestampContext, userId)
				update.Message.Text = "/outdatedcommand"
			} else {
				// There was a valid context. Update and use it.
				updatedCommand := storedCommand + " " + currentText
				userContext[userId] = updatedCommand
				userTimestampContext[userId] = currentTimestamp
				update.Message.Text = updatedCommand
			}
			// HACK: "fake" received update.Message.Entities to seem a command:
			me := tgbotapi.MessageEntity{
				Type:   "bot_command",
				Offset: 0,
				Length: len(strings.Split(update.Message.Text, " ")[0]), // Extract "/command" from "/command param1 param2 ..."
			}
			update.Message.Entities = append(update.Message.Entities, me)
		} else {
			// Message is a command. Update cache.
			userContext[userId] = currentText
			userTimestampContext[userId] = currentTimestamp
		}

		// Create a new MessageConfig. We don't have text yet, so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /createalert, /sayhi and /status."
		case "createalert":
			msg = telegramhandler.HandleCreateAlert(bot, update)
		case "sayhi":
			msg.Text = fmt.Sprintf("Hi %s :)", update.Message.From.FirstName)
		case "status":
			msg.Text = "I'm under construction."
		// TODO: case with donation options.
		// TODO: case to report bugs/suggestions.
		// Meta-commands to handle cornercases.
		case "nocommand":
			msg.Text = "I don't understand you. Please use /help to learn how to interact with me."
		case "outdatedcommand":
			msg.Text = "You took too long to finish your previous command. Please, start again."
		default:
			msg.Text = "I don't know that command, sorry!"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
