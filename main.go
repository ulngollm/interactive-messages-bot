package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
)

var bot *tele.Bot
var webAppURL string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}

	botToken := os.Getenv("TOKEN")
	pref := tele.Settings{
		Token:     botToken,
		ParseMode: tele.ModeMarkdownV2,
		Poller:    &tele.LongPoller{Timeout: 1 * time.Second},
	}
	bot, err = tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}
	webAppURL = os.Getenv("WEB_APP_URL")
}

func main() {
	bot.Handle("/start", openAppButton)
	bot.Handle("/help", help)
	bot.Handle("/send", selectChat)

	bot.Handle(tele.OnCallback, handler)
	bot.Handle(tele.OnWebApp, webAppCallback)
	//todo отправить сообщение в определенный чат. Сначала этот чат нужно запросить
	bot.Handle(tele.OnChatShared, sendToChat)

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

func openAppButton(c tele.Context) error {
	button := tele.Btn{Text: "Open Message Generator", WebApp: &tele.WebApp{URL: webAppURL}}
	keyboard := &tele.ReplyMarkup{}
	keyboard.Reply(
		keyboard.Row(button),
	)
	return c.Send("Чтобы перейти к созданию сообщения, нажмите на кнопку и перейдите в приложение", keyboard)
}

func webAppCallback(c tele.Context) error {
	response := c.Message().WebAppData.Data
	var data WebAppMessageData
	if err := json.Unmarshal([]byte(response), &data); err != nil {
		log.Printf("webAppCallback.json unmarshal: %s", err.Error())
		return c.Send("Не получилось сгенерировать сообщение. Повторите, пожалуйста")
	}

	keyboard := &tele.ReplyMarkup{}
	for _, row := range data.ReplyMarkup.Keyboard {
		var btnRow []tele.InlineButton
		for _, btn := range row {
			if btn.Text == "" {
				continue
			}
			inlineBtn := tele.InlineButton{
				Text: btn.Text,
			}
			if btn.URL != "" {
				inlineBtn.URL = btn.URL
			} else if btn.CallbackData != "" {
				inlineBtn.Data = btn.CallbackData
			}
			btnRow = append(btnRow, inlineBtn)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, btnRow)
	}
	return c.Send(data.Text, keyboard)
}

func selectChat(c tele.Context) error {
	keyboard := c.Bot().NewMarkup()
	flag := true
	button := tele.Btn{
		Text: "выбрать чат",
		Chat: &tele.ReplyRecipient{
			Channel: true,
			//todo синтаксически фигня какая-то
			Bot:       &flag,
			BotMember: &flag,
		},
	}
	keyboard.Reply(
		keyboard.Row(button),
	)
	return c.Send("выберите чат", keyboard)
}

func sendToChat(c tele.Context) error {
	id := c.Message().ChatShared.ChatID
	send, err := c.Bot().Send(&tele.Chat{ID: id}, "message")
	if err != nil {
		//todo добавить id ошибки. Отдавать пользователю и писать в лог
		log.Printf("sendToChat: %s", err.Error())
		return c.Send("Сообщение не отправлено в чат. Ошибка")
	}
	return c.Send(fmt.Sprintf("Сообщение %s отправлено в чат ID %d", send.Text, id))
}

func help(c tele.Context) error {
	return c.Send("как пользоваться...")
}
