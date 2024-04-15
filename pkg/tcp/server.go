package tcpserver

import (
	"CLI/pkg/misc"
	BlackBox "CLI/pkg/utils/blackbox"
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
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/labstack/gommon/random"
)

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
)

func handleCommand(conn net.Conn, message string) {
	args := strings.Split(message, " ")
	if args[0] == "exit" || args[0] == "quit" {
		fmt.Println("Closing connection...")
		conn.Close()
		return
	}
	if args[0] == "ai" {
		Aname := args[1]
		message_content := strings.Join(args[2:], " ")
		if Aname == "blackbox" {
			BlackBox_ := BlackBox.NewBlackboxClient()
			reply := BlackBox_.SendMessage(message_content, true)
			for {
				reader := bufio.NewReader(reply.Body)
				line, err := reader.ReadString('\n')
				if err != nil {
					//conn.Write([]byte("Error: " + err.Error() + "\n"))
					if err == io.EOF {
						conn.Write([]byte("\n\n[DONE]"))
					}
					break
				}
				conn.Write([]byte(line))
			}
		}
		if Aname == "hug" {
			model_ := args[2]
			client := HugginFace.NewHug()
			cookie := settings.HugginFaceCookie
			message_content := strings.Join(args[3:], " ")
			ChatId := client.ChangeModel(model_, cookie)
			Id_ := client.GetMsgUID(ChatId, cookie)
			err, r := client.SendMessage(message_content, ChatId, Id_, cookie, true)
			if err != nil && r.Body == nil {
				conn.Write([]byte("Error: " + err.Error() + "\n"))
			}
			reader := bufio.NewReader(r.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						conn.Write([]byte("\n\n[DONE]"))
						break
					}

				}
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				var event map[string]interface{}
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					conn.Write([]byte("Error: " + err.Error() + "\n"))
					continue
				}
				if event["type"] == "stream" {
					conn.Write([]byte(event["token"].(string)))
				}
			}
			fmt.Print("\r\n")
		}
		if Aname == "bing" {
			message_content := strings.Join(args[2:], " ")
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
				conn.Write([]byte(msg.Text))
				if msg.Error != nil {
					log.Printf("Error: %v", msg.Error)
				}
			}
			conn.Write([]byte("\n\n[DONE]"))
		}
		if Aname == "merlin" {
			authToken := settings.MerlinAuthToken
			chatID := "43ac5495-e1e1-4a68-9115-" + random.String(8)
			m := MerlinAI.NewMerlin(authToken, chatID)
			message_content := strings.Join(args[2:], " ")

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
					conn.Write([]byte("Error: " + scanner.Err().Error() + "\n"))
				}
				conn.Write([]byte(content))
			}
			conn.Write([]byte("\n\n[DONE]"))
			if err := scanner.Err(); err != nil {
				conn.Write([]byte("Error: " + err.Error() + "\n"))
			}
			fmt.Print("\r\n")

		}
		if Aname == "youai" {
			message_content := strings.Join(args[2:], " ")
			YouAI := youai.YouAIClient{}
			err, resp := YouAI.SendMessage(message_content, true)
			if err != nil {
				conn.Write([]byte("Error: " + err.Error() + "\n"))
			}
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						conn.Write([]byte("\n\n[DONE]"))
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
							log.Println("Error parsing JSON:", err)
							continue
						}
						// fmt.Print(jsondata["youChatToken"])
						conn.Write([]byte(jsondata["youChatToken"].(string)))
					}

				}
			}
		}

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			break
		}

		// Remove newline character from message
		message = strings.TrimSpace(message)

		fmt.Printf("Received command: %s\n", message)
		handleCommand(conn, message)
		// Process the received command here
		// ...
	}
}

var (
	stop = make(chan struct{})
)

func NewServer() {
	listen, err := net.Listen("tcp4", settings.TcpHost)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	fmt.Println("Listening on " + settings.TcpHost)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		sig := <-sigChan
		fmt.Printf("Caught signal %s: terminating\n", sig)
		close(stop)
	}()

	// Accept incoming connections
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				close(stop)
				return
			default:
				conn, err := listen.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err.Error())
					continue
				}

				// Handle each connection in a new goroutine
				go handleConnection(conn)
			}
		}
	}()

	// Wait for the server to finish
	wg.Wait()
}
