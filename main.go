package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"guiv1/misc"
	"guiv1/models/handler"
	"guiv1/models/huggingface"
	"guiv1/models/tuneapp"
	httpserver "guiv1/web"

	"guiv1/windows/hastebin"
	"guiv1/windows/mailraid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

var (
	f_                   = misc.Funcs{}
	icns                 = misc.IconUtil{}
	Chap                 = &misc.ChatApp{}
	hugmodels, err       = huggingface.NewHug().GetModels()
	cliperr              = clipboard.Init()
	tuneclient           = tuneapp.TuneClient{}
	tunemodels           = tuneclient.GetModels()
	AppSettings, setterr = f_.LoadSettings()
	CurrentAIProvider    = "merlin"
	CurrentHugModel      = AppSettings.CurrentHugModel
	CurrentTuneAppModel  = AppSettings.CurrentTuneModel
	HttpServerStatus     = "Stopped"
)

type ChatApp struct {
	app     fyne.App
	win     fyne.Window
	input   *widget.Entry
	chatLog *widget.Entry
}

func NewChatApp() *ChatApp {

	web_ := httpserver.Web{}
	web_.Init_Routes()
	a := app.New()

	w := a.NewWindow("S7 Gui V1")

	w.Resize(fyne.NewSize(1200, 800))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetIcon(icns.Icon("appicon"))
	w.SetTitle(fmt.Sprintf("S7 Gui V1 - Current AI Provider: %s - HTTP API: %s - MSG | > ", CurrentAIProvider, HttpServerStatus))
	if AppSettings.DarkMode {
		a.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.Settings().SetTheme(theme.LightTheme())
	}

	messagegrid := container.NewVBox()
	messagegrid.Layout = layout.NewVBoxLayout()

	messagegrid.Resize(fyne.NewSize(150, 20))

	chatLog := widget.NewMultiLineEntry()
	chatLog.Wrapping = fyne.TextWrapWord

	chatLog.Resize(fyne.NewSize(800, 600))
	chatLog.TextStyle.Monospace = true
	chatLog.TextStyle.Symbol = true

	scroll := container.NewVScroll(messagegrid)
	scroll.SetMinSize(fyne.NewSize(150, 700))
	scroll.Direction = container.ScrollBoth

	chatLog.OnChanged = func(s string) {
		chatLog.CursorRow = len(chatLog.Text) - 1
		scroll.ScrollToBottom()

	}
	scrollborder := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		scroll)
	input := widget.NewEntry()
	input.PlaceHolder = "Type a message..."
	input.OnSubmitted = func(s string) {
		if s != "" {
			messagegrid.Add(NewUserMessageElement(s))

			input.SetText("")
			getAIResponse(s, chatLog, w)

		}
	}
	SendBtn := widget.NewButton("Send", func() {
		text := input.Text
		if text != "" {
			messagegrid.Add(NewUserMessageElement(text))

			input.SetText("")
			getAIResponse(text, chatLog, w)

		}
	})
	MContainer := container.NewBorder(
		nil,
		nil,
		scrollborder,
		nil,
		chatLog,
	)
	MContainer.Resize(fyne.NewSize(900, 600))
	stop_ := 0
	httpServerStartBtn := widget.NewToolbarAction(icns.Icons8("256", "start.png", "fluency"), func() {
		httpserver.ProcStop = make(chan struct{})
		httpserver.WaitGroup.Add(1)
		go web_.NewHttpServer()

		HttpServerStatus = "Running"
		w.SetTitle(fmt.Sprintf("S7 Gui V1 - Current AI Provider: %s - HTTP API: %s - MSG | > ", CurrentAIProvider, HttpServerStatus))
		stop_ = 1
		f_.NotificationModal(w, &misc.ChatApp{a, w, input, chatLog}, "HTTP Server Started", "The HTTP server is now running on "+AppSettings.Httphost)
	})
	httpserverStopBtn := widget.NewToolbarAction(icns.Icons8("256", "stop.png", "fluency"), func() {
		if stop_ == 0 {
			f_.NotificationModal(w, &misc.ChatApp{a, w, input, chatLog}, "HTTP Server Error", "The HTTP server is not running ")
			return
		}
		// <-httpserver.ProcStop
		close(httpserver.ProcStop)
		// httpserver.WaitGroup.Done()
		HttpServerStatus = "Stopped"
		w.SetTitle(fmt.Sprintf("S7 Gui V1 - Current AI Provider: %s - HTTP API: %s - MSG | > ", CurrentAIProvider, HttpServerStatus))
		f_.NotificationModal(w, &misc.ChatApp{a, w, input, chatLog}, "HTTP Server Stopped", "The HTTP server has been stopped.")
		stop_ = 0
	})
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(icns.Icons8("256", "trash--v1.png", ""), func() {
			chatLog.SetText("")
			messagegrid.Objects = nil
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(icns.Icons8("256", "source-code.png", ""), func() {
			CodeModal(w, &ChatApp{a, w, input, chatLog})
		}),
		widget.NewToolbarAction(icns.Icons8("256", "copy--v1.png", ""), func() {
			//TexttoCopy := strings.Split(chatLog.Text, "Merlin:")[len(strings.Split(chatLog.Text, "Merlin: "))-1]
			clipboard.Write(clipboard.FmtText, []byte(chatLog.Text))
		}),
		widget.NewToolbarAction(icns.Icons8("256", "chatgpt.png", "nolan"), func() {
			ModelMenuModal(w, &ChatApp{a, w, input, chatLog})
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		httpServerStartBtn,
		widget.NewToolbarSeparator(),
		httpserverStopBtn,
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(icns.Icons8("256", "bomb-emoji.png", "emoji"), func() {
			mailraid.RaidWindow(a, w)
		}),
		widget.NewToolbarAction(icns.Icons8("256", "maintenance.png", ""), func() {
			ToolMenu(w, a)
		}),
		widget.NewToolbarAction(icns.Icons8("256", "help--v1.png", ""), func() {
			log.Println("Display help")
			HelpModalMarkdown(w, &ChatApp{a, w, input, chatLog})
		}),
		widget.NewToolbarAction(icns.Icons8("256", "services--v1.png", ""), func() {
			showSettingsModal(w, &ChatApp{a, w, input, chatLog})
		}),
		widget.NewToolbarAction(icns.Icons8("256", "fire-exit.png", ""), func() {
			w.Close()
		}),
	)
	// add tooltips to toolbar actions
	toolbar.Resize(fyne.NewSize(900, 100))
	SendBtn.SetIcon(icns.Icons8("256", "sent--v2.png", ""))
	Container_ := container.NewBorder(
		toolbar,
		container.NewGridWithColumns(2, input, SendBtn),
		nil,
		nil,
		MContainer,
	)
	w.SetContent(Container_)

	return &ChatApp{a, w, input, chatLog}
}

func HelpModalMarkdown(w fyne.Window, app *ChatApp) {

	resp, err := http.Get("https://pastebin.com/raw/ULD756dK")
	if err != nil {
		log.Println("Error fetching README.md:", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading README.md:", err)
		return
	}
	richmd := widget.NewRichTextFromMarkdown(string(body))
	scroll := container.NewVScroll(richmd)
	scroll.SetMinSize(fyne.NewSize(600, 500))
	CloseBtn := widget.NewButton("Close", nil)
	x := container.NewVBox(
		widget.NewLabel("Readme.md"),
		widget.NewSeparator(),
		scroll,
		container.NewHBox(
			layout.NewSpacer(),
			CloseBtn,
			layout.NewSpacer(),
		),
	)

	modal := widget.NewModalPopUp(x, w.Canvas())
	CloseBtn.OnTapped = func() {
		modal.Hide()
	}
	modal.Resize(fyne.NewSize(700, 600))
	modal.Show()
}

func CodeModal(w fyne.Window, App *ChatApp) {
	richcodemd := widget.NewRichTextFromMarkdown(App.chatLog.Text)
	richcodemd.Resize(fyne.NewSize(800, 550))
	scroll := container.NewVScroll(richcodemd)
	scroll.SetMinSize(fyne.NewSize(700, 600))
	copybtn := widget.NewButton("Copy", func() {
		clipboard.Write(clipboard.FmtText, []byte(App.chatLog.Text))
	})
	copybtn.SetIcon(icns.Icons8("256", "copy--v1.png", ""))
	OKBtn := widget.NewButton("OK", nil)
	OKBtn.SetIcon(icns.Icons8("256", "checkmark.png", ""))
	notification := container.NewVBox(
		scroll,
		container.NewHBox(
			layout.NewSpacer(),
			copybtn,
			OKBtn,
			layout.NewSpacer(),
		),
	)
	popup := widget.NewModalPopUp(notification, w.Canvas())
	// OKBtn := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-1].(*widget.Button)
	OKBtn.OnTapped = func() {
		popup.Hide()
	}

	popup.Resize(fyne.NewSize(900, 600))
	popup.Show()
}

func NewUserMessageElement(message string) *widget.RichText {
	userMessage := widget.NewRichTextFromMarkdown(fmt.Sprintf("**You:** %s", message))
	userMessage.Wrapping = fyne.TextWrapWord

	return userMessage
}

func getAIResponse(input string, chatlog *widget.Entry, w fyne.Window) string {
	fmt.Println("Current model:", CurrentAIProvider)
	w.SetTitle(fmt.Sprintf("S7 Gui V1 - Current AI Provider: %s - HTTP API: %s - MSG | > %s", CurrentAIProvider, HttpServerStatus, input))
	if CurrentAIProvider == "merlin" {
		handler.MerlinAI_(input, chatlog)
	}
	if CurrentAIProvider == "bing" {
		handler.BingAI(input, chatlog)

	}
	if CurrentAIProvider == "hugging-face" {
		handler.HuggingFaceAI(input, CurrentHugModel, chatlog)
	}
	if CurrentAIProvider == "black-box" {
		handler.BlackBoxAI(input, chatlog)
	}
	if CurrentAIProvider == "tune-app" {
		handler.TuneAppAI(input, CurrentTuneAppModel, chatlog)
	}
	if CurrentAIProvider == "youai" {
		handler.YouAI(input, chatlog)
	}
	if CurrentAIProvider == "openaigpt4" {
		handler.Gpt4AI(input, chatlog)
	}
	return "I'm not sure I understand. Can you please rephrase?"
}
func showSettingsModal(w fyne.Window, a *ChatApp) {
	if setterr != nil {
		log.Println("Error loading settings:", setterr)
	}
	settings := container.NewVBox(
		widget.NewLabel("Settings"),
		widget.NewTextGridFromString("TuneAPP Model:"),
		widget.NewSelect(tunemodels, func(model string) {
			CurrentTuneAppModel = model
		}),
		widget.NewTextGridFromString("HuggingFace Model:"),
		widget.NewSelect(hugmodels, func(model string) {
			CurrentHugModel = model
		}),
		widget.NewTextGridFromString("Bing Host:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("TCP Host:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("HTTP Host:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("YouAICookie:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("BingCookie:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("Hugging Face Cookie:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("TuneApp Access Token:"),
		widget.NewEntry(),
		widget.NewTextGridFromString("Merlin Auth Token:"),
		widget.NewEntry(),
		widget.NewCheck(("Dark Mode"), nil),
		widget.NewButton("Save", nil),
	)

	popup := widget.NewModalPopUp(settings, w.Canvas())
	popup.Resize(fyne.NewSize(400, 300))
	SaveBtn := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-1].(*widget.Button)
	title := popup.Content.(*fyne.Container).Objects[0].(*widget.Label)
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	DarkMode := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-2].(*widget.Check)
	DarkMode.OnChanged = func(b bool) {
		if b {
			a.app.Settings().SetTheme(theme.DarkTheme())
		} else {
			a.app.Settings().SetTheme(theme.LightTheme())
		}
	}
	MerlinAuthToken := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-3].(*widget.Entry)
	MerlinAuthToken.OnChanged = func(s string) {
		AppSettings.MerlinAuthToken = s
	}
	TuneAppAccessToken := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-5].(*widget.Entry)
	TuneAppAccessToken.OnChanged = func(s string) {
		AppSettings.TuneAppAccessToken = s
	}
	HugginFaceCookie := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-7].(*widget.Entry)
	HugginFaceCookie.OnChanged = func(s string) {
		AppSettings.HugginFaceCookie = s
	}
	BingCookie := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-9].(*widget.Entry)
	BingCookie.OnChanged = func(s string) {
		AppSettings.BingCookie = s
	}
	YouAICookie := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-11].(*widget.Entry)
	YouAICookie.OnChanged = func(s string) {
		AppSettings.YouAICookie = s
	}
	HTTPHost := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-13].(*widget.Entry)
	HTTPHost.OnChanged = func(s string) {
		AppSettings.Httphost = s
	}
	TCPHost := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-15].(*widget.Entry)
	TCPHost.OnChanged = func(s string) {
		AppSettings.TcpHost = s
	}
	BingHost := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-17].(*widget.Entry)
	BingHost.OnChanged = func(s string) {
		AppSettings.BingHost = s
	}
	CurrentHugModel_Dropdown := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-19].(*widget.Select)
	CurrentHugModel_Dropdown.OnChanged = func(model string) {
		handler.ChatCurrId = ""
		CurrentHugModel = model
	}
	CurrentHugModel_Dropdown.Selected = CurrentHugModel
	CurrentTuneModel_Dropdown := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-21].(*widget.Select)
	CurrentTuneModel_Dropdown.OnChanged = func(model string) {
		CurrentTuneAppModel = model
	}
	CurrentTuneModel_Dropdown.Selected = CurrentTuneAppModel
	SaveBtn.SetIcon(icns.Icons8("256", "save--v1.png", ""))
	SaveBtn.OnTapped = func() {
		updatedSettings := misc.Data{
			CurrentTuneModel:   CurrentTuneModel_Dropdown.Selected,
			BingHost:           AppSettings.BingHost,
			CurrentHugModel:    CurrentHugModel_Dropdown.Selected,
			DarkMode:           DarkMode.Checked,
			HugginFaceCookie:   AppSettings.HugginFaceCookie,
			MerlinAuthToken:    AppSettings.MerlinAuthToken,
			TuneAppAccessToken: AppSettings.TuneAppAccessToken,
			YouAICookie:        AppSettings.YouAICookie,
			BingCookie:         AppSettings.BingCookie,
			Username:           AppSettings.Username,
			Password:           AppSettings.Password,
			TcpHost:            AppSettings.TcpHost,
			Httphost:           AppSettings.Httphost,
		}
		updatedJSON, err := json.MarshalIndent(updatedSettings, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling updated settings:", err)
			return
		}
		err = ioutil.WriteFile("settings.json", updatedJSON, 0644)
		if err != nil {
			fmt.Println("Error writing JSON file:", err)
			return
		}
		fmt.Println("Settings file updated successfully")
		popup.Hide()
		NotificationModal(w, a, "Settings Saved", "The settings have been saved successfully.")
	}
	DarkMode.Checked = AppSettings.DarkMode
	HugginFaceCookie.Text = AppSettings.HugginFaceCookie
	MerlinAuthToken.Text = AppSettings.MerlinAuthToken
	BingCookie.Text = AppSettings.BingCookie
	YouAICookie.Text = AppSettings.YouAICookie
	TuneAppAccessToken.Text = AppSettings.TuneAppAccessToken
	TCPHost.Text = AppSettings.TcpHost
	HTTPHost.Text = AppSettings.Httphost
	BingHost.Text = AppSettings.BingHost
	popup.Show()
}

func NotificationModal(w fyne.Window, a *ChatApp, title string, message string) {
	title_element := widget.NewRichTextFromMarkdown(fmt.Sprintf("# %s", title))

	notification := container.NewBorder(
		title_element,
		container.NewVBox(
			widget.NewRichTextFromMarkdown(fmt.Sprintf("##  %s", message)),
			widget.NewButton("OK", nil),
		),
		nil,
		nil,
	)
	popup := widget.NewModalPopUp(notification, w.Canvas())
	popup.CreateRenderer().Layout(notification.Size())
	OKBtn := notification.Objects[1].(*fyne.Container).Objects[1].(*widget.Button)
	OKBtn.SetIcon(icns.Icons8("256", "checkmark--v1.png", ""))
	OKBtn.OnTapped = func() {
		popup.Hide()
	}
	popup.Resize(fyne.NewSize(300, 150))
	popup.Show()
}

func ModelMenuModal(w fyne.Window, a *ChatApp) {
	merlin_btn := widget.NewButton("Merlin", func() {
		CurrentAIProvider = "merlin"

	})
	merlin_btn.SetIcon(icns.Icon("merlin"))
	bing_btn := widget.NewButton("Bing", func() {
		CurrentAIProvider = "bing"

	})
	bing_btn.SetIcon(icns.Icons8("90", "bing--v1.png", "fluency"))
	hugging_face_btn := widget.NewButton("Hugging Face", func() {
		CurrentAIProvider = "hugging-face"
	})
	hugging_face_btn.SetIcon(icns.Icon("huggingface"))
	blackbox_btn := widget.NewButton("Black Box", func() {
		CurrentAIProvider = "black-box"
	})
	blackbox_btn.SetIcon(icns.Icon("blackbox"))
	tuneapp_btn := widget.NewButton("Tune App", func() {
		CurrentAIProvider = "tune-app"
	})
	tuneapp_btn.SetIcon(icns.Icon("tuneapp"))
	youai_btn := widget.NewButton("YouAI", func() {
		CurrentAIProvider = "youai"
	})
	youai_btn.SetIcon(icns.Icon("youai"))
	gpt4_btn := widget.NewButton("GPT-4", func() {
		CurrentAIProvider = "openaigpt4"
	})
	gpt4_btn.SetIcon(icns.Icons8("256", "chatgpt.png", ""))
	title_ := widget.NewRichTextFromMarkdown("# Select AI Provider:" + misc.InvisFill)
	modelMenu := container.NewVBox(
		container.NewHBox(
			title_,
			widget.NewSeparator(),
			widget.NewToolbar(
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(icns.Icons8("256", "cancel--v1.png", ""), nil),
			),
		),
		merlin_btn,
		bing_btn,
		hugging_face_btn,
		blackbox_btn,
		tuneapp_btn,
		youai_btn,
		gpt4_btn,
	)

	popup := widget.NewModalPopUp(modelMenu, w.Canvas())
	popup.Resize(fyne.NewSize(300, 200))
	toolbar_ := popup.Content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Toolbar)
	toolbar_.Resize(fyne.NewSize(300, 100))
	toolbar_.Refresh()
	cancelbtn := popup.Content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Toolbar).Items[3].(*widget.ToolbarAction)
	cancelbtn.OnActivated = func() {
		popup.Hide()
	}
	cancelbtn.SetIcon(icns.Icons8("256", "cancel--v1.png", ""))
	cancelbtn.ToolbarObject().Resize(fyne.NewSize(100, 100))
	cancelbtn.ToolbarObject().Refresh()
	toolbar_.Refresh()
	Providers := []string{"merlin", "bing", "hugging-face", "black-box", "tune-app", "youai", "openaigpt4"}
	for i, btn := range popup.Content.(*fyne.Container).Objects {
		if _, ok := btn.(*fyne.Container); !ok {
			btn.(*widget.Button).OnTapped = func() {
				CurrentAIProvider = Providers[i-1]
				w.SetTitle(fmt.Sprintf("S7 Gui V1 - Current AI Provider: %s - HTTP API: %s - MSG | > ", CurrentAIProvider, HttpServerStatus))
				popup.Hide()

			}
		}
	}

	popup.Show()
}

func ToolMenu(w fyne.Window, a fyne.App) {
	title_ := widget.NewRichTextFromMarkdown("# Select a Tool:" + misc.InvisFill)
	newHaste := widget.NewButton("New Hastebin", func() {
		hastebin.NewHastebin(w, a)
	})
	newHaste.SetIcon(icns.IconFromBytes("hastebin", icns.IconByteLoader("hb_400x400", "")))
	modelMenu := container.NewVBox(
		container.NewHBox(
			title_,
			widget.NewSeparator(),
			widget.NewToolbar(
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarSpacer(),
				widget.NewToolbarAction(icns.Icons8("256", "cancel--v1.png", ""), nil),
			),
		),
		newHaste,
	)

	popup := widget.NewModalPopUp(modelMenu, w.Canvas())
	popup.Resize(fyne.NewSize(300, 200))
	toolbar_ := popup.Content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Toolbar)
	toolbar_.Resize(fyne.NewSize(300, 100))
	toolbar_.Refresh()
	cancelbtn := popup.Content.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Toolbar).Items[3].(*widget.ToolbarAction)
	cancelbtn.OnActivated = func() {
		popup.Hide()
	}
	cancelbtn.SetIcon(icns.Icons8("256", "cancel--v1.png", ""))
	cancelbtn.ToolbarObject().Resize(fyne.NewSize(100, 100))
	cancelbtn.ToolbarObject().Refresh()
	toolbar_.Refresh()

	popup.Show()
}
func main() {
	app := NewChatApp()
	app.win.ShowAndRun()
}
