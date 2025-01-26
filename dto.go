package main

type WebAppMessageData struct {
	Text        string `json:"text"`
	ReplyMarkup struct {
		Keyboard [][]struct {
			Text         string `json:"text"`
			URL          string `json:"url"`
			CallbackData string `json:"callback_data"`
		} `json:"inline_keyboard"`
	} `json:"reply_markup"`
}
