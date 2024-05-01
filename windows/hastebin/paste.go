package hastebin

import (
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"io"
	"log"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

func NewHastebin(w fyne.Window, app fyne.App) {

	richmd := widget.NewMultiLineEntry()
	scroll := container.NewVScroll(richmd)
	scroll.SetMinSize(fyne.NewSize(600, 500))
	CloseBtn := widget.NewButton("Close", nil)
	PasteBtn := widget.NewButton("Paste", nil)
	x := container.NewVBox(
		widget.NewLabel("New Paste"),
		widget.NewSeparator(),
		scroll,
		container.NewHBox(
			layout.NewSpacer(),
			CloseBtn,
			PasteBtn,
			layout.NewSpacer(),
		),
	)

	modal := widget.NewModalPopUp(x, w.Canvas())
	CloseBtn.OnTapped = func() {
		modal.Hide()
	}
	PasteBtn.OnTapped = func() {
		// Get the text from the multi-line entry
		pasteText := richmd.Text
		if pasteText == "" {
			misc.Funcs{}.NotificationModal(w, &misc.ChatApp{}, "Error", "Please enter some text to paste.")
			return
		}
		client := &http.Client{}
		var data = strings.NewReader(pasteText)
		req, err := http.NewRequest("POST", "https://hastebin.skyra.pw/documents", data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("accept", "application/json")
		req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
		req.Header.Set("content-type", "text/plain")
		req.Header.Set("origin", "https://hastebin.skyra.pw")
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("referer", "https://hastebin.skyra.pw/")
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
		var event map[string]interface{}
		err = json.Unmarshal(bodyText, &event)
		uuid := event["key"].(string)
		misc.Funcs{}.NotificationModal(w, &misc.ChatApp{}, "Paste Uploaded", fmt.Sprintf("Your paste has been uploaded to: https://hastebin.skyra.pw/%s, Link copied to clipboard.", uuid))
		clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("https://hastebin.skyra.pw/%s", uuid)))
		richmd.Text = ""
		// modal.Hide()

	}
	modal.Resize(fyne.NewSize(700, 600))
	modal.Show()
}
