package gpt4

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GPT4Client struct {
	Chatid   string
	Messages []struct {
		ID        string `json:"id"`
		Role      string `json:"role"`
		Content   string `json:"content"`
		Who       string `json:"who"`
		Timestamp int64  `json:"timestamp"`
	}
}

func (c GPT4Client) SendMessage(NewMessage string) http.Response {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"botId":"default","customId":null,"session":"N/A","chatId":"%s","contextId":100,"messages":[],"newMessage":"%s","newFileId":null,"stream":true}`, c.Chatid, NewMessage))
	req, err := http.NewRequest("POST", "https://chatgpt4online.org/wp-json/mwai-ui/v1/chats/submit", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://chatgpt4online.org/")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("X-WP-Nonce", "2de9874226")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return *resp
}
