package tuneapp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TuneClient struct {
}

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
	cookie_util   = misc.CookieUtil{}
	// cookie        = misc.CookieUtil{}.ReadCookiesFile(settings.TuneAppCookie)
	AccessToken                   = settings.TuneAppAccessToken
	cookie, AccessToken_New, errx = TuneClient{}.GetCookieAuto()
)

func (c TuneClient) NewUuid() string {
	uuid_base, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uuid_base.String()
}

func (c TuneClient) NewChat(model_name string) string {
	if AccessToken == "" {
		AccessToken = AccessToken_New
	}
	// Initialize a new HTTP client
	if err != nil {
		log.Fatal(err)
	}
	chat_id := c.NewUuid()
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://chat.tune.app/api/new?conversation_id="+chat_id+"&model="+model_name+"&currency=USD", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("authorization", AccessToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", "https://chat.tune.app/")
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
	jsondata := make(map[string]interface{})
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &jsondata)
	if err != nil {
		log.Fatal(err)
	}
	// Process the JSON response
	successStatus, ok := jsondata["success"].(bool)
	//fmt.Println("Response: ", jsondata["success"])
	if !ok {
		log.Fatal("Failed to get success status from response")
	}
	if !successStatus {
		fmt.Println("Error creating new chat")
		return "err"
	}
	return chat_id
}

func (c TuneClient) SendMessage(message string, chat_id string, model string, internet bool, raw bool) (http.Response, error) {
	if AccessToken == "" {
		AccessToken = AccessToken_New
	}
	client := &http.Client{}
	var data = strings.NewReader(`{"query":"` + message + `","conversation_id":"` + chat_id + `","model_id":"` + model + `","browseWeb":` + strconv.FormatBool(internet) + `,"attachement":"","attachment_name":"","messageId":"` + c.NewUuid() + `","prevMessageId":"` + c.NewUuid() + `"}`)
	req, err := http.NewRequest("POST", "https://chat.tune.app/api/prompt", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("authorization", AccessToken)
	req.Header.Set("content-type", "text/plain;charset=UTF-8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("origin", "https://chat.tune.app")
	req.Header.Set("referer", "https://chat.tune.app/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	//fmt.Println(resp.StatusCode)
	// fmt.Println(resp.Header)
	if err != nil {
		log.Fatal(err)
	}
	if !raw {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return http.Response{
						StatusCode: 500,
					}, nil
				}
			}
			var response map[string]interface{}
			// decode line
			json.Unmarshal(line, &response)
			if err != nil {
				fmt.Println("Error decoding JSON response:", err)
			}
			value_ := response["value"]
			if value_ != nil {
				fmt.Print(value_)
			}

		}
	} else {
		return *resp, nil
	}
	return *resp, err
}

func (c TuneClient) GetCookieAuto() (string, string, error) {

	// http request https://chat.tune.app/
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://chat.tune.app/", nil)
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	if err != nil {
		return "nil", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "nil", "", err
	}
	var cookieFields []string
	for _, field := range resp.Header.Values("set-cookie") {
		cookieFields = append(cookieFields, strings.Split(field, ";")[0])
	}
	initial_cookie := strings.Join(cookieFields, "; ")
	client2 := &http.Client{}
	var data = strings.NewReader(`{"utms":{}}`)
	req2, err := http.NewRequest("POST", "https://chat.tune.app/api/guestLogin", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", initial_cookie)
	req.Header.Set("origin", "https://chat.tune.app")
	req.Header.Set("referer", "https://chat.tune.app/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp2, err := client2.Do(req2)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response map[string]interface{}
	// decode line
	json.Unmarshal(bodyText2, &response)
	accessToken, ok := response["accessToken"].(string)
	if !ok {
		return "nil", "", fmt.Errorf("failed to get access token from response")
	}

	newCookies := strings.Join(cookieFields, "; ")

	//fmt.Println("New cookies:", newCookies)
	return newCookies, accessToken, nil
}

func (c TuneClient) GetConversations() ([]map[string]interface{}, error) {
	if AccessToken == "" {
		AccessToken = AccessToken_New
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://chat.tune.app/api/conversations", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("authorization", AccessToken)
	req.Header.Set("cookie", cookie)
	req.Header.Set("referer", "https://chat.tune.app/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	jsondata := make(map[string]interface{})
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &jsondata)
	if err != nil {
		log.Fatal(err)
	}
	conv := jsondata["conversations"].([]interface{})
	conversations := make([]map[string]interface{}, len(conv))
	for i, v := range conv {
		conversations[i] = v.(map[string]interface{})
	}
	return conversations, nil

}

func (c TuneClient) DeleteConversation(conversationId string) error {
	if AccessToken == "" {
		AccessToken = AccessToken_New
	}
	client := &http.Client{}
	var data = strings.NewReader(`{"conversation_id":"` + conversationId + `"}`)
	req, err := http.NewRequest("PUT", "https://chat.tune.app/api/deleteConversation", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("authorization", AccessToken)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", cookie)
	req.Header.Set("origin", "https://chat.tune.app")
	req.Header.Set("referer", "https://chat.tune.app/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// json
	json_data := make(map[string]interface{})
	err = json.Unmarshal(bodyText, &json_data)
	if err != nil {
		return err
	}
	successStatus := json_data["success"].(bool)
	if successStatus {
		return nil
	}
	return err
}
func (c TuneClient) GetModels() []string {

	client := &http.Client{}
	var data = strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "https://studio.tune.app/tune.Studio/ListPublicModels", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")
	req.Header.Set("authorization", "Bearer null")
	req.Header.Set("connect-protocol-version", "1")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "_ga_EN6NJLMZ17=GS1.1.1714074056.1.1.1714074062.0.0.0")
	req.Header.Set("origin", "https://studio.tune.app")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://studio.tune.app/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
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

	var json_data map[string]interface{}
	err = json.Unmarshal(bodyText, &json_data)
	if err != nil {
		log.Fatal(err)
	}
	models := json_data["models"].([]interface{})
	models_list := make([]string, len(models))
	for i, model := range models {
		models_list[i] = model.(map[string]interface{})["uri"].(string)
	}
	return models_list
}
