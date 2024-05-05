package dalle3

import (
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/html"
)

var (
	f_              = misc.Funcs{}
	cu              = misc.CookieUtil{}
	cookie          = cu.ReadCookiesFile("./data/chatgatecookies.json")
	cookie_raw, err = cu.ReadCookiesFileRaw("./data/chatgatecookies.json")
	mwai_session_id = cu.GetCValue("mwai_session_id", cookie_raw)
	icns            = misc.IconUtil{}
	iu              = misc.ImageUtil{}
	ApiNonce        = GetNonce()
)

func GetNonce() string {
	nonce := ""
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://chatgate.ai/dalle-e-3/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("cookie", cookie)
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://chatgate.ai/")
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
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return ""
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "id" {
					if a.Val == "jetpack-instant-search-js-before" {
						data := n.FirstChild.Data
						data2 := strings.Split(strings.Split(data, `var JetpackInstantSearchOptions=JSON.parse(decodeURIComponent("`)[1], `"));`)[0]
						urldecoded, err := url.QueryUnescape(data2)
						if err != nil {
							log.Fatal(err)
						}
						var options map[string]interface{}
						err = json.Unmarshal([]byte(urldecoded), &options)
						if err != nil {
							log.Fatal(err)
						}
						nonce = options["apiNonce"].(string)

					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nonce
}

func DalleImage(prompt string) (string, error) {

	client := &http.Client{}
	var data = strings.NewReader(`{"botId":"chatbot-5rwkvr","customId":null,"session":"` + mwai_session_id + `","chatId":"mc8a16m09rg","contextId":` + "216" + `,"messages":[],"newMessage":"` + prompt + `","newFileId":null,"stream":false}`)
	req, err := http.NewRequest("POST", "https://chatgate.ai/wp-json/mwai-ui/v1/chats/submit", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", cookie)
	req.Header.Set("origin", "https://chatgate.ai")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://chatgate.ai/dalle-e-3/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("x-wp-nonce", ApiNonce)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	var json_event map[string]interface{}

	err = json.Unmarshal(bodyText, &json_event)
	if err != nil {
		log.Fatal(err)
	}
	if json_event["success"] == nil {
		return "", fmt.Errorf("DALL-E API response did not contain a 'success' field")
	}
	if json_event["success"].(bool) {
		image_url := json_event["images"].([]interface{})[0].(string)
		log.Println("Generated image URL:", image_url)
		return image_url, nil
	} else {
		log.Println("Failed to generate image:", strings.Split(json_event["message"].(string), ",")[0])
		return "", fmt.Errorf("%s", strings.Split(json_event["message"].(string), ",")[0])
	}
}

func DalleWindow(a fyne.App, w fyne.Window) {
	infinite_progress := widget.NewProgressBarInfinite()

	PromptEntry := widget.NewEntry()
	PromptEntry.PlaceHolder = "Enter prompt here"
	PromptSubmit := widget.NewButton("Generate", nil)
	PromptContainer := container.NewGridWithColumns(
		2,
		PromptEntry,
		PromptSubmit,
	)
	Grid := container.NewCenter()

	title_label := widget.NewRichTextFromMarkdown("# Image Generation ( DALLE 3 ) " + misc.InvisFill + misc.InvisFill)
	TopBorder := container.NewBorder(
		container.NewHBox(
			title_label,
			widget.NewSeparator(),
			widget.NewToolbar(
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(icns.Icons8("256", "cancel--v1.png", ""), nil),
			),
		),
		Grid,
		nil,
		nil,
		nil,
	)
	Images_ := container.NewBorder(
		TopBorder,
		container.NewVBox(PromptContainer, infinite_progress),
		nil,
		nil,
		nil,
	)
	Grid.Refresh()
	Images_.Refresh()

	content := container.NewBorder(
		Images_,
		nil,
		nil,
		nil,
		nil,
	)

	modal := widget.NewModalPopUp(content, w.Canvas())

	TopBorder.Objects[0].(*fyne.Container).Objects[2].(*widget.Toolbar).Items[3].(*widget.ToolbarAction).OnActivated = func() {
		modal.Hide()
	}
	infinite_progress.Hidden = true
	PromptSubmit.OnTapped = func() {
		infinite_progress.Hidden = false
		if err != nil {
			log.Fatalf("Error creating animated GIF: %v", err)
		}
		Grid.RemoveAll()

		PromptEntry.Text = ""
		PromptEntry.Refresh()
		text := PromptEntry.Text
		Gen, err := DalleImage(text)
		if err != nil {

		}
		Grid.RemoveAll()
		infinite_progress.Hidden = true
		if Gen != "" {
			image_url := Gen
			Grid.Add(iu.NewCanvasImageUri(512, 512, image_url))
			Grid.Refresh()
			Images_.Refresh()
		} else {
			log.Println("Error generating image:", err)
			misc.Funcs{}.NotificationModal(w, &misc.ChatApp{}, "Error", string(err.Error()))
		}
	}
	PromptEntry.OnSubmitted = func(text string) {
		infinite_progress.Hidden = false
		if err != nil {
			log.Fatalf("Error creating animated GIF: %v", err)
		}
		Grid.RemoveAll()

		PromptEntry.Text = ""
		PromptEntry.Refresh()
		Gen, err := DalleImage(text)
		if err != nil {

		}
		Grid.RemoveAll()
		infinite_progress.Hidden = true
		if Gen != "" {
			image_url := Gen
			Grid.Add(iu.NewCanvasImageUri(512, 512, image_url))
			Grid.Refresh()
			Images_.Refresh()
		} else {

			log.Println("Error generating image:", err)
			misc.Funcs{}.NotificationModal(w, &misc.ChatApp{}, "Error", string(err.Error()))
		}
	}
	content.Refresh()
	modal.Resize(fyne.NewSize(700, 600))
	modal.Canvas.SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEscape {
			modal.Hide()
		}
	})
	modal.Show()
}
