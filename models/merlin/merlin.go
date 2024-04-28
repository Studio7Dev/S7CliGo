package merlin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"guiv1/models/tuneapp"
	"io"
	"log"
	"net/http"
)

type Merlin struct {
	AuthToken string
	ChatID    string
	Client    *http.Client
}

func NewMerlin(authToken, chatID string) *Merlin {
	return &Merlin{
		AuthToken: authToken,
		ChatID:    chatID,
		Client:    &http.Client{},
	}
}

var ActiveThreadSnippet []interface{}
var OldMerlinMsg = ""

func (m *Merlin) Chat(message string) (io.Reader, error) {
	OldMerlinMsg = message
	fmt.Println("ChatID:", m.ChatID)
	data := map[string]interface{}{
		"action": map[string]interface{}{
			"message": map[string]interface{}{
				"attachments": []interface{}{},
				"content":     message,
				"metadata": map[string]string{
					"context": "",
				},
				"parentId": "root",
				"role":     "user",
			},
			"type": "NEW",
		},
		"activeThreadSnippet": ActiveThreadSnippet,
		"chatId":              m.ChatID,
		"language":            "AUTO",
		"metadata":            nil,
		"mode":                "VANILLA_CHAT",
		"model":               "GPT 3",
		"personaConfig":       map[string]interface{}{},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://uam.getmerlin.in/thread/unified?customJWT=true&version=1.1", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.AuthToken)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("x-merlin-version", "extension-null")

	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (m *Merlin) RawStream(message string) (http.Response, error) {

	return http.Response{}, nil
}

func (m *Merlin) StreamContent(responseBody io.Reader) error {
	scanner := bufio.NewScanner(responseBody)
	allfullcontent := ""
	UserSnippet := make([]interface{}, 0)
	AISnippet := make([]interface{}, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0:6] != "data: " {
			continue
		}
		line = line[6:]

		var event map[string]interface{}
		err := json.Unmarshal([]byte(line), &event)
		if err != nil {
			log.Println("Error parsing event JSON:", err)
			continue
		}

		data, ok := event["data"].(map[string]interface{})
		if !ok {
			log.Println("Error: data field not found in event")
			continue
		}

		content, ok := data["content"].(string)
		if !ok {
			log.Println("Error: content field not found in event data")
			continue
		}
		allfullcontent += content
		fmt.Print(content)
	}
	parentid := tuneapp.TuneClient{}.NewUuid()
	UserSnippet = append(UserSnippet, map[string]interface{}{
		"attachments": []interface{}{},
		"content":     OldMerlinMsg,
		"id":          parentid,
		"metadata": []interface{}{
			map[string]interface{}{
				"key":   "context",
				"value": "This is the context for the user message.",
			},
		},
		"parentId":       "root",
		"role":           "user",
		"status":         "SUCCESS",
		"activeChildIdx": 0,
		"totalChildren":  1,
		"idx":            0,
		"totSiblings":    1,
	})
	ai_id := tuneapp.TuneClient{}.NewUuid()
	AISnippet = append(AISnippet, map[string]interface{}{
		"content":        allfullcontent,
		"id":             ai_id,
		"parentId":       parentid,
		"role":           "assistant",
		"status":         "SUCCESS",
		"activeChildIdx": 0,
		"totalChildren":  0,
		"idx":            0,
		"totSiblings":    1,
	})
	ActiveThreadSnippet = append(ActiveThreadSnippet, &UserSnippet)
	ActiveThreadSnippet = append(ActiveThreadSnippet, &AISnippet)
	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Print("\r\n")
	return nil
}
