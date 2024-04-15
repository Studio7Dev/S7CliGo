package youai

import (
	"CLI/pkg/misc"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type YouAIClient struct {
}

type Params struct {
	Query              string
	Page               int
	Count              int
	SafeSearch         string
	ResponseFilters    []string
	Domain             string
	Personalization    bool
	QueryTraceID       string
	ChatID             string
	ConversationTurnID string
	PastChatLength     int
	SelectedChatMode   string
	UseTracing         bool
	TraceID            string
	Chat               []interface{}
}

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
)

func (c YouAIClient) NewUuid() string {
	uuid_base, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uuid_base.String()
}

func (c YouAIClient) ParamBuilder(message string) string {
	params := Params{
		Query:              message,
		Page:               1,
		Count:              10,
		SafeSearch:         "Off",
		ResponseFilters:    []string{"WebPages", "TimeZone", "Computation", "RelatedSearches"},
		Domain:             "youchat",
		Personalization:    true,
		QueryTraceID:       c.NewUuid(),
		ChatID:             c.NewUuid(),
		ConversationTurnID: c.NewUuid(),
		PastChatLength:     0,
		SelectedChatMode:   "default",
		UseTracing:         true,
		Chat:               []interface{}{},
	}
	baseurl := "https://you.com/api/streamingSearch"
	prompt_parms := ""
	prompt_parms += "?q=" + url.QueryEscape(params.Query)
	prompt_parms += "&page=" + url.QueryEscape(strconv.Itoa(params.Page))
	prompt_parms += "&count=" + url.QueryEscape(strconv.Itoa(params.Count))
	prompt_parms += "&safeSearch=" + url.QueryEscape(params.SafeSearch)
	prompt_parms += "&domain=" + url.QueryEscape(params.Domain)
	prompt_parms += "&queryTraceId=" + url.QueryEscape(params.QueryTraceID)
	prompt_parms += "&chatId=" + url.QueryEscape(params.ChatID)
	prompt_parms += "&conversationTurnId=" + url.QueryEscape(params.ConversationTurnID)
	prompt_parms += "&pastChatLength=" + strconv.Itoa(params.PastChatLength)
	prompt_parms += "&selectedChatMode=" + url.QueryEscape(params.SelectedChatMode)
	// urlencode
	urlencoded := prompt_parms
	return baseurl + urlencoded

}

func (c YouAIClient) SendMessage(message string, raw bool) (error, http.Response) {
	url := c.ParamBuilder(message)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "text/event-stream")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", settings.YouAICookie)
	req.Header.Set("referer", "https://you.com/search?q=hi&fromSearchBar=true&tbm=youchat&chatMode=default")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-full-version", `"123.0.6312.122"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-platform-version", `"10.0.0"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// defer resp.Body.Close()
	// bodyText, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", bodyText)

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/event-stream;charset=utf-8" {
		return nil, http.Response{
			StatusCode: 401,
		}
	}
	if raw {
		return nil, *resp
	} else {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return err, http.Response{
					StatusCode: 402,
				}
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
						return err, http.Response{
							StatusCode: 402,
						}
					}
					fmt.Print(jsondata["youChatToken"])
				}

			}
		}

		fmt.Print("\r\n")
		return nil, *resp
	}
	return nil, http.Response{
		StatusCode: 403,
	}
}
