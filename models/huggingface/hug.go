package huggingface

import (
	"bufio"
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type ChatClient struct {
	Client *http.Client
}

func NewHug() *ChatClient {
	return &ChatClient{
		Client: &http.Client{},
	}
}

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
	cookie        = misc.CookieUtil{}.ReadCookiesFile(settings.HugginFaceCookie)
)

func (s *ChatClient) FindRandomUUID(text string) string {
	uuids := make([]string, 0)

	guidPattern := `[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[a-f0-9]{4}-[a-f0-9]{12}`
	re := regexp.MustCompile(guidPattern)

	foundGUIDs := re.FindAllString(text, -1)

	for _, g := range foundGUIDs {
		uuids = append(uuids, g)
	}

	if len(uuids) > 0 {
		rand.Seed(time.Now().UnixNano())
		return uuids[rand.Intn(len(uuids))]
	}

	return ""
}

func (s *ChatClient) GetMsgUID(chatID string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://huggingface.co/chat/conversation/"+chatID+"/__data.json?x-sveltekit-invalidated=11", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", "https://huggingface.co/chat/conversation/6608a05392dfb775db102588")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return s.FindRandomUUID(string(bodyText))
}

func (s *ChatClient) ChangeModel(model string) string {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"model":"%s","preprompt":""}`, model))
	req, err := http.NewRequest("POST", "https://huggingface.co/chat/conversation", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", cookie)
	req.Header.Set("origin", "https://huggingface.co")
	req.Header.Set("referer", "https://huggingface.co/chat/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var event map[string]interface{}
	if err := json.Unmarshal([]byte(bodyText), &event); err != nil {
		return "err"
	}
	convId := event["conversationId"]
	return string(fmt.Sprint(convId))
}

func (c *ChatClient) SendMessage(message string, convId string, Id string, raw bool) (error, http.Response) {
	data := strings.NewReader(fmt.Sprintf(`{"inputs":"%s","id":"%s","is_retry":false,"is_continue":false,"web_search":false,"files":[]}`, message, Id))
	req, err := http.NewRequest("POST", "https://huggingface.co/chat/conversation/"+convId, data)
	if err != nil {
		return err, http.Response{}
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", cookie)
	req.Header.Set("origin", "https://huggingface.co")
	req.Header.Set("referer", "https://huggingface.co/chat/conversation/6608a05392dfb775db102588")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err, http.Response{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
		}
	}
	// defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/event-stream" {
		return fmt.Errorf("unexpected content type: %s", contentType), *resp
	}

	if !raw {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return err, http.Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header,
				}
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				return err, http.Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header,
				}
			}
			if event["type"] == "stream" {
				fmt.Print(event["token"])
			}
		}

		fmt.Print("\r\n")
	}
	if raw {
		return nil, *resp
	}

	return nil, *resp
}

func (c ChatClient) GetModels() ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://huggingface.co/chat/models", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("cookie", "_ga=GA1.1.695792381.1711492971; __stripe_mid=1f7b2183-e484-497b-bb6a-1df73c8ad9ddbeaca1; token=wVUCLTYNPSmXIMRTElPuYAYAvnaXclCAKFbmrYVUZpWgLJWSzQrOvYLJEieNktjCueFdxsTcTXEUAjYxkShOoCAMoIqEZMhyhTBIsAfGNPvyLASsGXNQBursRrzoNYAh; token=wVUCLTYNPSmXIMRTElPuYAYAvnaXclCAKFbmrYVUZpWgLJWSzQrOvYLJEieNktjCueFdxsTcTXEUAjYxkShOoCAMoIqEZMhyhTBIsAfGNPvyLASsGXNQBursRrzoNYAh; hf-chat=ee19994f-12dd-4bd1-bb55-e8324018495a; _ga_8Q63TH4CSL=GS1.1.1713208077.22.0.1713208077.60.0.0; _ga_R4JMGZWPD9=GS1.1.1713214248.1.0.1713214255.0.0.0")
	req.Header.Set("if-none-match", `W/"1rdws4r"`)
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// doc, err := html.Parse(strings.NewReader(string(bodyText)))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var traverse func(*html.Node)
	// traverse = func(n *html.Node) {
	// 	if n.Type == html.ElementNode {
	// 		fmt.Println(n.Data) // Print the name of the HTML element
	// 	}
	// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 		traverse(c)
	// 	}
	// }

	// Traverse the HTML document
	// traverse(doc)

	doc, err := html.Parse(strings.NewReader(string(bodyText)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil, err
	}

	// Find all "dt" tags and print their text
	// make a string list
	models := []string{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "dt" {
			model_name := strings.Split(n.FirstChild.Data, " ")[0]
			models = append(models, model_name)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	// get elements by the html tag "dd"
	// extract the text content of each "dd" element and append it to the models slice
	var g func(*html.Node)
	g = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "dd" {
			model_name := strings.TrimSpace(n.FirstChild.Data)
			models = append(models, model_name)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			g(c)
		}
	}

	return models, nil
}
