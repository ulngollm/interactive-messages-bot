package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
)

var bot *tele.Bot

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}

	botToken := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:     botToken,
		ParseMode: tele.ModeMarkdown,
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
}

func main() {
	bot.Handle(tele.OnCallback, handler)

	log.Println("bot is starting...")
	bot.Start()
}

func handler(c tele.Context) error {
	if err := c.Respond(&tele.CallbackResponse{
		CallbackID: c.Callback().ID,
		Text:       c.Callback().Data,
		ShowAlert:  true,
	}); err != nil {
		log.Printf("handler.Respond: %s", err.Error())
		log.Printf("callback: %v", map[string]interface{}{
			"respondErr":   err,
			"callbackID":   c.Callback().ID,
			"callbackData": c.Callback().Data,
			"messageText":  c.Callback().Message.Text,
		})
	}
	return nil
}
