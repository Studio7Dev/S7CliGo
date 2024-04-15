package httpserver

import (
	"CLI/pkg/misc"
	BlackBox "CLI/pkg/utils/blackbox"
	"CLI/pkg/utils/goliath"
	HugginFace "CLI/pkg/utils/huggingface"
	MerlinAI "CLI/pkg/utils/merlin"
	"CLI/pkg/utils/sydney"
	"CLI/pkg/utils/util"
	"CLI/pkg/utils/youai"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/gommon/random"
)

var (
	f_                = misc.Funcs{}
	settings, setterr = f_.LoadSettings()
)

func MerlinChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		// Process the request body
		message_content := bodyJson["message"].(string)
		authToken := settings.MerlinAuthToken
		chatID := "43ac5495-e1e1-4a68-9115-" + random.String(8)
		m := MerlinAI.NewMerlin(authToken, chatID)

		responseBody, err := m.Chat(message_content)
		if err != nil {
			log.Fatalf("Error chatting with Merlin: %v", err)
		}

		scanner := bufio.NewScanner(responseBody)
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

			w.Write([]byte(content))
			w.(http.Flusher).Flush()
			//conn.Write([]byte(content))
		}
		w.Write([]byte("\n[DONE]"))
		w.(http.Flusher).Flush()
		//conn.Write([]byte("\n[DONE]"))
		if err := scanner.Err(); err != nil {
			return
		}
		fmt.Print("\r\n")
	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}
}

func HugChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		// Process the request body
		model_ := bodyJson["model"].(string)
		message_content := bodyJson["message"].(string)

		client := HugginFace.NewHug()
		cookie := settings.HugginFaceCookie

		ChatId := client.ChangeModel(model_, cookie)
		Id_ := client.GetMsgUID(ChatId, cookie)
		err, r := client.SendMessage(message_content, ChatId, Id_, cookie, true)
		if err != nil && r.Body == nil {
			w.Write([]byte("Error: " + err.Error() + "\n"))
		}
		reader := bufio.NewReader(r.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					//conn.Write([]byte("\n[DONE]"))
					break
				}

			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				//conn.Write([]byte("Error: " + err.Error() + "\n"))
				continue
			}
			if event["type"] == "stream" {
				//conn.Write([]byte(event["token"].(string)))
				w.Write([]byte(event["token"].(string)))
				w.(http.Flusher).Flush()
			}
		}
		w.Write([]byte("\n[DONE]"))
		w.(http.Flusher).Flush()
		fmt.Print("\r\n")

	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}
}

func BingChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		message_content := bodyJson["message"].(string)
		cookies, err := util.ReadCookiesFile()
		if err != nil {
			log.Fatalf("Error reading cookies file: %v", err)
		}
		sydney_ := sydney.NewSydney(sydney.Options{
			Debug:                 false,
			Cookies:               cookies,
			Proxy:                 "",
			ConversationStyle:     "",
			Locale:                "en-US",
			WssDomain:             "",
			CreateConversationURL: "",
			NoSearch:              false,
			GPT4Turbo:             true,
		})
		ch, err := sydney_.AskStream(sydney.AskStreamOptions{
			StopCtx:        context.TODO(),
			Prompt:         message_content,
			WebpageContext: "",
			ImageURL:       "",
		})
		if err != nil {
			log.Fatalf("Error creating Sydney instance: %v", err)
		}
		for msg := range ch {
			w.Write([]byte(msg.Text))
			w.(http.Flusher).Flush()
			if msg.Error != nil {
				log.Printf("Error: %v", msg.Error)
			}
		}
		w.Write([]byte("\n[DONE]"))
		w.(http.Flusher).Flush()
	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}

}
func BlackBoxChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		// Process the request body
		message_content := bodyJson["message"].(string)
		BlackBox_ := BlackBox.NewBlackboxClient()
		reply := BlackBox_.SendMessage(message_content, true)
		for {
			reader := bufio.NewReader(reply.Body)
			line, err := reader.ReadString('\n')
			if err != nil {
				//conn.Write([]byte("Error: " + err.Error() + "\n"))
				if err == io.EOF {
					w.Write([]byte("\n[DONE]"))
					w.(http.Flusher).Flush()
				}
				break
			}
			w.Write([]byte(line))
			w.(http.Flusher).Flush()
		}
	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}
}
func YouAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		message_content := bodyJson["message"].(string)
		YouAI := youai.YouAIClient{}
		err, resp := YouAI.SendMessage(message_content, true)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					w.Write([]byte("\n[DONE]"))
					w.(http.Flusher).Flush()
					break
				}
				break
			}
			// starts with
			if strings.HasPrefix(line, "data: ") {
				// remove "data:" prefix
				line = strings.TrimPrefix(line, "data: ")
				// remove leading/trailing whitespace
				line = strings.TrimSpace(line)
				// check if line is empty
				if len(line) == 0 {
					continue
				}
				if strings.HasPrefix(line, `{"youChatToken": "`) {
					jsondata := make(map[string]interface{})
					err := json.Unmarshal([]byte(line), &jsondata)
					if err != nil {
						http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
					}
					// fmt.Print(jsondata["youChatToken"])
					w.Write([]byte(jsondata["youChatToken"].(string)))
					w.(http.Flusher).Flush()
				}

			}
		}
	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}
}

func GoliathAIChat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		message_content := bodyJson["message"].(string)
		client := goliath.GoliathClient{}
		resp, err := client.SendMessage(message_content, true)
		if err != nil {
			fmt.Println("Error sending message to Goliath:", err)
			return
		}
		reader := bufio.NewReader(resp.Body)
		for {
			if err != nil {
				if err == io.EOF {
					break
				}
				w.Write([]byte("[DONE]"))
				w.(http.Flusher).Flush()
			}
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					w.Write([]byte("\n[DONE]"))
					w.(http.Flusher).Flush()
					break
				}
				http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if strings.HasPrefix(line, "data: ") {
				line = line[6:]
				jsondata := make(map[string]interface{})
				err := json.Unmarshal([]byte(line), &jsondata)
				if err != nil {

					continue
				}
				if jsondata["choices"] != nil {
					content_ := jsondata["choices"]
					for _, choice := range content_.([]interface{}) {
						choice = choice.(map[string]interface{})["delta"]
						content := choice.(map[string]interface{})["content"]
						if content != nil {
							w.Write([]byte(content.(string)))
							w.(http.Flusher).Flush()
						}
					}
				}
			}

		}
	} else {
		http.Error(w, "Only POST requests supported", http.StatusNotImplemented)
		return
	}

}

func NewHttpServer() {
	fmt.Println("Starting HTTP server...")
	if setterr != nil {
		log.Printf("Error loading settings: %v", setterr)
		return
	}
	// Register the handler function for the path "/process".
	http.HandleFunc("/merlin", MerlinChat)
	http.HandleFunc("/hug", HugChat)
	http.HandleFunc("/bing", BingChat)
	http.HandleFunc("/blackbox", BlackBoxChat)
	http.HandleFunc("/youai", YouAIChat)
	http.HandleFunc("/goliath", GoliathAIChat)
	fmt.Println("HTTP server started on:", settings.Httphost)
	// Start the web server on port 8080 and handle requests to localhost.
	if err := http.ListenAndServe(settings.Httphost, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
