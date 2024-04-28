package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"guiv1/models/blackbox"
	"guiv1/models/huggingface"
	"guiv1/models/merlin"
	"guiv1/models/tuneapp"
	"guiv1/models/youai"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var (
	f_               = misc.Funcs{}
	AppSettings, err = f_.LoadSettings()
	BlackBox_        = blackbox.NewBlackboxClient()
	tuneclient       = tuneapp.TuneClient{}
	c, errx          = tuneclient.GetConversations()
	tune_chat_id     = c[0]["conversation_id"].(string)
)

func YouAI(message string, chatlog *widget.Entry) error {
	message_content := message
	YouAI := youai.YouAIClient{}
	err, resp := YouAI.SendMessage(message_content, true)
	if err != nil {
		//http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return nil
	}
	reader := bufio.NewReader(resp.Body)
	chatlog.SetText(chatlog.Text + "YouAI" + ": " + "" + "\n ")
	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// w.Write([]byte("\n[DONE]"))
				// w.(http.Flusher).Flush()
				break
			}
			break
		}
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, `{"youChatToken": "`) {
				jsondata := make(map[string]interface{})
				err := json.Unmarshal([]byte(line), &jsondata)
				if err != nil {
					return nil
				}
				// fmt.Print(jsondata["youChatToken"])

				chatlog.SetText(chatlog.Text + jsondata["youChatToken"].(string))
				// w.Write([]byte(jsondata["youChatToken"].(string)))
				// w.(http.Flusher).Flush()
			}

		}
	}
	chatlog.SetText(chatlog.Text + "\n")
	return nil
}

var PrevMsgId = ""
var UsrPrevId = ""

func TuneAppAI(message string, model_ string, chatlog *widget.Entry) error {
	message_content := message

	settings_, err := f_.LoadSettings()
	if err != nil {
		log.Fatalf("Error loading settings: %v", err)
	}
	if settings_.TuneAppAccessToken == "" {
		tuneclient.NewChat(tuneclient.GetModels()[0])
	}

	if err != nil {
		log.Fatalf("Error getting conversations: %v", err)
	}

	resp, err := tuneclient.SendMessage(message_content, tune_chat_id, model_, false, true, PrevMsgId, UsrPrevId)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	reader := bufio.NewReader(resp.Body)
	chatlog.SetText(chatlog.Text + model_ + ": " + "" + "\n ")
	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// w.Write([]byte("\n[DONE]"))
				// w.(http.Flusher).Flush()
				break
			} else {
				break
			}
		}
		var response map[string]interface{}
		// decode line
		json.Unmarshal(line, &response)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
		}
		PrevMsgId_ := response["assistantMessageId"]
		UsrPrevId_ := response["userMessageId"]
		if PrevMsgId_ != nil {
			PrevMsgId = PrevMsgId_.(string)
			UsrPrevId = UsrPrevId_.(string)
			fmt.Println("Previous user message ID:", UsrPrevId)
			fmt.Println("Previous message ID:", PrevMsgId)
		}
		value_ := response["value"]
		if value_ != nil {
			chatlog.SetText(chatlog.Text + value_.(string))
			// w.Write([]byte(value_.(string)))
			// w.(http.Flusher).Flush()
		}

	}
	chatlog.SetText(chatlog.Text + "\n")

	return nil
}

func BlackBoxAI(message string, chatlog *widget.Entry) error {
	message_content := message

	reply := BlackBox_.SendMessage(message_content, true)
	chatlog.SetText(chatlog.Text + "BlackBox: " + "" + "\n ")
	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
	for {
		reader := bufio.NewReader(reply.Body)
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		chatlog.SetText(chatlog.Text + string(line))
	}
	chatlog.SetText(chatlog.Text + "\n")
	return nil
}

func BingAI(message string, chatlog *widget.Entry) error {
	message_content := message
	client := &http.Client{}
	var data = strings.NewReader(`{"message":"` + message_content + `"}`)
	req, err := http.NewRequest("POST", "http://"+AppSettings.BingHost+"/bing", data)
	if err != nil {
		chatlog.SetText(chatlog.Text + err.Error() + "\n")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		chatlog.SetText(chatlog.Text + err.Error() + "\n")
		return err
	}
	chatlog.SetText(chatlog.Text + "Bing: " + "" + "\n ")
	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading response: %v", err)
		}
		if string(line) == "[DONE]" {
			break
		}

		chatlog.SetText(chatlog.Text + string(line))
		chatlog.Wrapping = fyne.TextWrapWord
	}
	chatlog.SetText(chatlog.Text + "\n")
	return nil
}

var chatID = tuneapp.TuneClient{}.NewUuid()

func MerlinAI_(message string, chatlog *widget.Entry) error {
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer settingsFile.Close()

	data, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var result struct {
		MerlinAuthToken   string `json:"merlin_auth_token"`
		HuggingFaceCookie string `json:"huggingface_cookie"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	authToken := result.MerlinAuthToken

	m := merlin.NewMerlin(authToken, chatID)

	responseBody, err := m.Chat(message)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(responseBody)
	allfullcontent := ""
	UserSnippet := make([]interface{}, 0)
	AISnippet := make([]interface{}, 0)
	chatlog.SetText(chatlog.Text + "Merlin: " + "" + "\n ")

	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
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
		if scanner.Err() != nil {
			log.Println("Error reading response body:", scanner.Err())

		}
		allfullcontent += content
		chatlog.SetText(chatlog.Text + content)
		chatlog.Wrapping = fyne.TextWrapWord

		//fmt.Print(content)
		// w.Write([]byte(content))
		// w.(http.Flusher).Flush()
	}
	parentid := tuneapp.TuneClient{}.NewUuid()
	UserSnippet = append(UserSnippet, map[string]interface{}{
		"attachments": []interface{}{},
		"content":     merlin.OldMerlinMsg,
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
	merlin.ActiveThreadSnippet = append(merlin.ActiveThreadSnippet, &UserSnippet)
	merlin.ActiveThreadSnippet = append(merlin.ActiveThreadSnippet, &AISnippet)
	chatlog.SetText(chatlog.Text + "\n")
	// w.Write([]byte("\n[DONE]"))
	// w.(http.Flusher).Flush()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

var client = huggingface.NewHug()
var ChatCurrId = ""

func HuggingFaceAI(message string, model_ string, chatlog *widget.Entry) error {

	message_content := message

	if ChatCurrId == "" {
		ChatCurrId = client.ChangeModel(model_)
	}
	Id_ := client.GetMsgUID(ChatCurrId)
	err, r := client.SendMessage(message_content, ChatCurrId, Id_, true)
	if err != nil && r.Body == nil {
		return err
	}
	reader := bufio.NewReader(r.Body)
	chatlog.SetText(chatlog.Text + model_ + ": " + "" + "\n ")
	//chatlog.SetText(chatlog.Text + "=======================================================================================\n")
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event["type"] == "stream" {
			// w.Write([]byte(event["token"].(string)))
			// w.(http.Flusher).Flush()
			// chatlog.SetText(chatlog.Text + event["token"].(string))
			chatlog.SetText(chatlog.Text + event["token"].(string))
			chatlog.Wrapping = fyne.TextWrapWord
		}
	}
	// w.Write([]byte("\n[DONE]"))
	// w.(http.Flusher).Flush()
	chatlog.SetText(chatlog.Text + "\n")
	fmt.Print("\r\n")
	return nil
}
