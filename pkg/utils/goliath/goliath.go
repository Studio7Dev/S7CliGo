package goliath

import (
	"CLI/pkg/misc"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
)

type GoliathClient struct {
}

func (c GoliathClient) SendMessage(message string, raw bool) (http.Response, error) {
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	var data = strings.NewReader(`{"messages":[{"role":"user","content":"` + message + `"}],"model":"rohan/goliath-120b-16k-gptq","stream":true,"temperature":0.8,"max_tokens":16300}`)
	req, err := http.NewRequest("POST", "https://proxy.tune.app/chat/completions", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", settings.GoliathAuthToken)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://studio.tune.app")
	req.Header.Set("referer", "https://studio.tune.app/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//defer resp.Body.Close()
	// bodyText, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s\n", bodyText)
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/event-stream" {
		return http.Response{
				StatusCode: 400,
			},
			fmt.Errorf("unexpected content type: %s", contentType)
	}

	if !raw {
		reader := bufio.NewReader(resp.Body)
		for {
			if err != nil {
				if err == io.EOF {
					break
				}
				return http.Response{
					StatusCode: 402,
				}, err
			}
			line, err := reader.ReadString('\n')
			if err != nil {
				return http.Response{}, err
			}
			if strings.HasPrefix(line, "data: ") {
				line = line[6:]
				jsondata := make(map[string]interface{})
				err := json.Unmarshal([]byte(line), &jsondata)
				if err != nil {
					return http.Response{}, err
				}
				if jsondata["choices"] != nil {
					content_ := jsondata["choices"]
					for _, choice := range content_.([]interface{}) {
						choice = choice.(map[string]interface{})["delta"]
						content := choice.(map[string]interface{})["content"]
						if content != nil {
							fmt.Print(content.(string))
						}
					}
				}
			}

		}

	} else {
		return *resp, nil
	}

	return http.Response{}, nil
}
