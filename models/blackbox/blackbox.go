package blackbox

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type BlackboxClient struct {
	client *http.Client
}

func NewBlackboxClient() *BlackboxClient {
	return &BlackboxClient{
		client: &http.Client{},
	}
}

func (bb *BlackboxClient) GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	id := make([]byte, length)
	for i := range id {
		id[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(id)
}

func (bb *BlackboxClient) GenerateUserID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", bb.GenerateID(8), bb.GenerateID(4), bb.GenerateID(4), bb.GenerateID(4), bb.GenerateID(12))
}

var userID = NewBlackboxClient().GenerateUserID()

func (bb *BlackboxClient) SendMessage(content string, raw bool) http.Response {

	messageID := bb.GenerateID(7)
	data := strings.NewReader(fmt.Sprintf(`{"webSearchMode":"false","messages":[{"id":"%s","role":"user","content":"%s"}],"id":"%s","previewToken":null,"userId":"%s","codeModelMode":true,"agentMode":{},"trendingAgentMode":{},"isMicMode":false,"isChromeExt":false,"githubToken":null}`, messageID, content, messageID, userID))
	req, err := http.NewRequest("POST", "https://www.blackbox.ai/api/chat", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://www.blackbox.ai")
	req.Header.Set("referer", "https://www.blackbox.ai/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := bb.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if !raw {
		for {
			reader := bufio.NewReader(resp.Body)
			line, err := reader.ReadString('\n')
			if err != nil {
				return http.Response{}
			}
			fmt.Print(line)
		}
	}
	if raw {
		return *resp
	}
	return *resp
}
