package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"guiv1/misc"
	"guiv1/models/handler"
	"guiv1/models/huggingface"
	"guiv1/models/tuneapp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	f_                   = misc.Funcs{}
	hugmodels, err       = huggingface.NewHug().GetModels()
	tuneclient           = tuneapp.TuneClient{}
	tunemodels           = tuneclient.GetModels()
	AppSettings, setterr = f_.LoadSettings()
	CurrentAIProvider    = "merlin"
	CurrentHugModel      = AppSettings.CurrentHugModel
	CurrentTuneAppModel  = AppSettings.CurrentTuneModel
)

type ChatApp struct {
	app     fyne.App
	win     fyne.Window
	input   *widget.Entry
	chatLog *widget.Entry
}

func NewChatApp() *ChatApp {
	// go func() {
	// 	httpserver.NewHttpServer()
	// }()
	a := app.New()

	w := a.NewWindow("S7 Gui V1")
	// borderless window
	//w.SetFullScreen(true)
	w.Resize(fyne.NewSize(1000, 700))
	if AppSettings.DarkMode {
		a.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.Settings().SetTheme(theme.LightTheme())
	}

	input := widget.NewEntry()
	input.PlaceHolder = "Type a message..."

	chatLog := widget.NewMultiLineEntry()
	chatLog.Wrapping = fyne.TextWrapWord

	chatLog.Resize(fyne.NewSize(1000, 500))
	//chatLog.Disable()
	chatLog.TextStyle.Monospace = true
	chatLog.TextStyle.Symbol = true

	chatLog.OnChanged = func(s string) {
		chatLog.CursorRow = len(chatLog.Text) - 1

	}

	SendBtn := widget.NewButton("Send", func() {
		text := input.Text
		if text != "" {
			chatLog.SetText(chatLog.Text + "You: " + text + "\n")
			input.SetText("")
			getAIResponse(text, chatLog)

		}
	})
	ChangeModelBtn := widget.NewButton("Providers", func() {
		ModelMenuModal(w, &ChatApp{a, w, input, chatLog})
	})
	ClearBtn := widget.NewButton("Clear", func() {
		chatLog.SetText("")
	})
	ExitBtn := widget.NewButton("Exit", func() {
		a.Quit()
	})
	SettingsBtn := widget.NewButton("Settings", func() {
		showSettingsModal(w, &ChatApp{a, w, input, chatLog})
	})
	Container_ := container.NewBorder(
		nil,
		container.NewGridWithColumns(2, input, SendBtn),
		nil,
		container.NewGridWithRows(10, ExitBtn, ClearBtn, SettingsBtn, ChangeModelBtn),
		container.NewAdaptiveGrid(1, chatLog),
	)
	//menu := container.NewGridWithColumns(2, ExitBtn, ClearBtn)
	//cont := container.NewBorder(menu, container.NewGridWithColumns(2, input, SendBtn), nil, nil, Container_)
	w.SetContent(Container_)

	return &ChatApp{a, w, input, chatLog}
}

func getAIResponse(input string, chatlog *widget.Entry) string {
	fmt.Println("Current model:", CurrentAIProvider)
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
		// Implement Black Box AI logic here
		handler.BlackBoxAI(input, chatlog)
	}
	if CurrentAIProvider == "tune-app" {
		handler.TuneAppAI(input, CurrentTuneAppModel, chatlog)
	}
	if CurrentAIProvider == "youai" {
		handler.YouAI(input, chatlog)
	}
	// TO DO: implement AI logic here
	// for now, just return a simple response
	return "I'm not sure I understand. Can you please rephrase?"
}
func showSettingsModal(w fyne.Window, a *ChatApp) {
	if setterr != nil {
		log.Println("Error loading settings:", setterr)
	}
	// Create a settings modal
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
	) // Len 6

	popup := widget.NewModalPopUp(settings, w.Canvas())
	popup.Resize(fyne.NewSize(400, 300))
	SaveBtn := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-1].(*widget.Button)

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
	CurrentHugModel_Dropdown.Selected = CurrentHugModel
	CurrentTuneModel_Dropdown := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-21].(*widget.Select)
	CurrentTuneModel_Dropdown.Selected = CurrentTuneAppModel
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
	notification := container.NewVBox(
		widget.NewLabel(title),
		widget.NewLabel(message),
		widget.NewButton("OK", nil),
	)
	popup := widget.NewModalPopUp(notification, w.Canvas())
	OKBtn := popup.Content.(*fyne.Container).Objects[len(popup.Content.(*fyne.Container).Objects)-1].(*widget.Button)
	OKBtn.OnTapped = func() {
		popup.Hide()
	}
	popup.Resize(fyne.NewSize(300, 150))
	popup.Show()
}

func ModelMenuModal(w fyne.Window, a *ChatApp) {
	// Create a model selection modal
	modelMenu := container.NewVBox(
		widget.NewLabel("Select Model Provider"),
		widget.NewButton("Merlin", func() {
			CurrentAIProvider = "merlin"

		}),
		widget.NewButton("Bing", func() {
			CurrentAIProvider = "bing"

		}),
		widget.NewButton("Hugging Face", func() {
			CurrentAIProvider = "hugging-face"

		}),
		widget.NewButton("Black Box", func() {
			CurrentAIProvider = "black-box"
		}),
		widget.NewButton("Tune App", func() {
			CurrentAIProvider = "tune-app"
		}),
		widget.NewButton("YouAI", func() {
			CurrentAIProvider = "youai"
		}),
		// widget.NewButton("Cancel", func() {
		// 	CurrentAIProvider = CurrentAIProvider
		// }),
	)

	popup := widget.NewModalPopUp(modelMenu, w.Canvas())
	popup.Resize(fyne.NewSize(300, 200))
	Providers := []string{"merlin", "bing", "hugging-face", "black-box", "tune-app", "youai"}
	for i, btn := range popup.Content.(*fyne.Container).Objects {
		// skip the first label and add an OnTapped handler to each button
		if _, ok := btn.(*widget.Label); !ok {
			btn.(*widget.Button).OnTapped = func() {
				CurrentAIProvider = Providers[i-1]
				popup.Hide()

			}
		}
	}

	popup.Show()
}

func main() {
	app := NewChatApp()
	app.win.ShowAndRun()
}
