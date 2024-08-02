package bot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bh "github.com/venturarome/DaftWatch/internal/handler/bot"
)

type Bot struct {
	BotApi               *tgbotapi.BotAPI
	Handler              *bh.TelegramHandler
	userContext          map[int64]string
	userTimestampContext map[int64]int64 // TODO try Redis/Memcached
	contextTtl           int64           // TODO ¿maybe use time.Duration?
}

func InstanceTelegramBot() *Bot {
	tgBotApi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}
	tgBotApi.Debug = os.Getenv("TELEGRAM_MODE") != "release"
	log.Printf("Authorized on account %s", tgBotApi.Self.UserName)

	tgHandler := bh.InstanceHandler()

	bot := Bot{
		BotApi:               tgBotApi,
		Handler:              tgHandler,
		userContext:          make(map[int64]string),
		userTimestampContext: make(map[int64]int64),
		contextTtl:           600, // 10 minutes
	}

	go bot.periodicContextCleanup()

	return &bot
}

func (bot *Bot) StartLongPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 50 // 50 sec is, apparently, maximum allowed

	updates := bot.BotApi.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		userId := update.Message.From.ID
		currentTimestamp := int64(update.Message.Date)
		currentText := strings.Trim(update.Message.Text, " ")

		// DEBUG
		fmt.Println("[DEBUG] Text received: ", currentText)

		if !update.Message.IsCommand() {
			// Probably it comes from a multistep command. Let's check out "in-memory database".
			storedCommand, found := bot.userContext[userId]
			lastInteractionTimestamp := bot.userTimestampContext[userId]
			if !found {
				// No user context.
				update.Message.Text = "/nocommand"
			} else if currentTimestamp-lastInteractionTimestamp > bot.contextTtl {
				// Outdated user context.
				delete(bot.userContext, userId)
				delete(bot.userTimestampContext, userId)
				update.Message.Text = "/outdatedcommand"
			} else {
				// There was a valid context. Update and use it.
				updatedCommand := storedCommand + " " + currentText
				bot.userContext[userId] = updatedCommand
				bot.userTimestampContext[userId] = currentTimestamp
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
			bot.userContext[userId] = currentText
			bot.userTimestampContext[userId] = currentTimestamp
		}

		// Create a new MessageConfig. We don't have text yet, so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		var clearContext bool

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /createalert, /sayhi and /status."
		case "myalerts":
			msg, clearContext = bot.Handler.HandleMyAlerts(update)
		case "createalert":
			msg, clearContext = bot.Handler.HandleCreateAlert(update)
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

		if clearContext {
			delete(bot.userContext, userId)
			delete(bot.userTimestampContext, userId)
		}

		if _, err := bot.BotApi.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func (bot *Bot) periodicContextCleanup() {
	// Clean all context expired values to avoid massive growth of junk data.
	time.Sleep(time.Hour)
	currentTimestamp := time.Now().Unix()
	for userId, contextTimestamp := range bot.userTimestampContext {
		if currentTimestamp-contextTimestamp > bot.contextTtl {
			delete(bot.userContext, userId)
			delete(bot.userTimestampContext, userId)
		}
	}
}
