package misc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	Auth "main/pkg/auth"
	commandLib "main/pkg/commands"
	MerlinAI "main/pkg/utils/merlin"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
)

var (
	authStatus = false
	sysInfo    = SystemInfo{}
)

type SystemInfo struct {
	GOARCH string
	GOOS   string
}

type Funcs struct {
	// Add any necessary functions or methods here
}

type Data struct {
	MerlinAuthToken  string `json:"merlin_auth_token"`
	HugginFaceCookie string `json:"huggingface_cookie"`
	BlackBoxCookie   string `json:"blackbox_cookie"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

func (f Funcs) LoadSettings() (Data, error) {
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		return Data{}, fmt.Errorf("error opening JSON file: %v", err)
	}
	//defer func(settingsFile *os.File) {
	//	_ = settingsFile.Close()
	//}(settingsFile)

	data, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return Data{}, fmt.Errorf("error reading JSON file: %v", err)
	}

	var result Data
	if err := json.Unmarshal(data, &result); err != nil {
		return Data{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return result, nil
}

func (f Funcs) Banner() string {
	banner := commandLib.Reset + `
    ███████╗███████╗     █████╗ ██╗     ██████╗██╗     ██╗
    ██╔════╝╚════██║    ██╔══██╗██║    ██╔════╝██║     ██║
    ███████╗    ██╔╝    ███████║██║    ██║     ██║     ██║
    ╚════██║   ██╔╝     ██╔══██║██║    ██║     ██║     ██║
    ███████║   ██║      ██║  ██║██║    ╚██████╗███████╗██║
    ╚══════╝   ╚═╝      ╚═╝  ╚═╝╚═╝     ╚═════╝╚══════╝╚═╝
    `
	coloredBanner := strings.ReplaceAll(banner, "█", commandLib.BoldPurple+"█"+commandLib.Reset)
	fmt.Println(coloredBanner)
	return ""
}

func (f Funcs) LoginForm() (Data, error) {
	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Login", func() {
			app.Stop()
		}).
		AddButton("Quit", func() {
			app.Stop()
			os.Exit(0)
		})

	form.SetButtonsAlign(tview.AlignCenter)

	form.SetBorder(true)
	form.SetTitle("Please Login to the Cli")
	form.SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		return Data{}, fmt.Errorf("error running form: %v", err)
	}
	username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	return Data{Username: username, Password: password}, nil
}

func (f Funcs) UpdateSettings(result Data) error {
	file, err := os.OpenFile("settings.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	datax, err := json.MarshalIndent(&result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	if err := ioutil.WriteFile("settings.json", datax, 0644); err != nil {
		return fmt.Errorf("error writing to JSON file: %v", err)
	}
	return nil
}

func (f Funcs) Authenticated() bool {
	return authStatus
}
func (f Funcs) SetAuthStatus(status bool) {
	authStatus = status
}
func (f Funcs) LoginRequest(username, password string) (int, error) {
	authService := Auth.AuthService{}
	response, err := authService.PerformLogin(username, password)
	if err != nil {
		return 0, err
	}
	if response.StatusCode == 200 {
		f.SetAuthStatus(true)
		return response.StatusCode, nil
	} else {
		f.SetAuthStatus(false)
		return response.StatusCode, fmt.Errorf("login failed with status code: %d", response.StatusCode)
	}
}

func (f Funcs) MerlinAI_(args []string, this commandLib.Command) error {
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer settingsFile.Close()

	data, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var result struct {
		MerlinAuthToken   string `json:"merlin_auth_token"`
		HuggingFaceCookie string `json:"huggingface_cookie"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	authToken := result.MerlinAuthToken
	chatID := "43ac5495-e1e1-4a68-9115-" + this.Name
	m := MerlinAI.NewMerlin(authToken, chatID)
	message := strings.Join(args[1:], " ")

	responseBody, err := m.Chat(message)
	if err != nil {
		return fmt.Errorf("error chatting with Merlin: %v", err)
	}

	if err := m.StreamContent(responseBody); err != nil {
		return fmt.Errorf("error streaming content: %v", err)
	}
	return nil
}

func (f Funcs) isUnix() bool {
	return sysInfo.GOOS == "linux" || sysInfo.GOOS == "darwin"
}

func (f Funcs) isMacOs() bool {
	return sysInfo.GOOS == "darwin"
}

func (f Funcs) OpenUrl(url string) {
	var cmdName string

	switch {
	case f.isUnix():
		cmdName = "xdg-open"
	case f.isMacOs():
		cmdName = "open"
	default:
		cmdName = "start"
	}

	cmdArgs := []string{"cmd.exe", "/C", cmdName, url}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error occurred starting the command '%s'. Error:\n%s", strings.Join(cmdArgs, " "), err.Error())
		fmt.Println("\nOutput:", string(output[:]))
	} else {
		// fmt.Printf("Successfully started command '%s'\n", strings.Join(cmdArgs, " "))
	}
}

func (f Funcs) ImageURLToBase64(url string) string {
	// Fetch the image
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Read the image data
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// Encode image data to base64
	base64String := base64.StdEncoding.EncodeToString(imageData)

	return base64String
}

func (f Funcs) SettingsPage() {
	settings_file, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	data, err := ioutil.ReadAll(settings_file)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	type Data struct {
		Key1     string `json:"merlin_auth_token"`
		Key2     string `json:"huggingface_cookie"`
		Key3     string `json:"blackbox_cookie"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var result Data
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	defer settings_file.Close()
	app := tview.NewApplication()
	pages := tview.NewPages()

	form := tview.NewForm()
	form.Box.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	form.SetTitle("Settings")
	form.SetBorder(true)
	form.SetTitleAlign(tview.AlignCenter)
	form.SetButtonsAlign(tview.AlignCenter)
	form.AddInputField("Merlin Auth Token", result.Key1, 50, nil, nil)
	form.AddInputField("Hugging Face Cookie", result.Key2, 50, nil, nil)
	form.AddInputField("Blackbox Cookie", result.Key3, 50, nil, nil)
	form.AddInputField("Username", result.Username, 50, nil, nil)
	form.AddInputField("Password", result.Password, 50, nil, nil)
	form.AddButton("Save", func() {
		// Save the updated settings to the JSON file
		updatedData := Data{
			Key1:     form.GetFormItemByLabel("Merlin Auth Token").(*tview.InputField).GetText(),
			Key2:     form.GetFormItemByLabel("Hugging Face Cookie").(*tview.InputField).GetText(),
			Key3:     form.GetFormItemByLabel("Blackbox Cookie").(*tview.InputField).GetText(),
			Username: form.GetFormItemByLabel("Username").(*tview.InputField).GetText(),
			Password: form.GetFormItemByLabel("Password").(*tview.InputField).GetText(),
		}
		// Write the updated settings to the JSON file
		updatedJSON, err := json.MarshalIndent(updatedData, "", "  ")
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
		// Write the updated settings to the JSON file
		app.Stop()
	})
	form.AddButton("Cancel", func() {
		// force_clear()
		app.Stop()
	})
	// clear the form and reset the page
	form.AddButton("Reset Cookies", func() {
		form.GetFormItemByLabel("Merlin Auth Token").(*tview.InputField).SetText("")
		form.GetFormItemByLabel("Hugging Face Cookie").(*tview.InputField).SetText("")
		form.GetFormItemByLabel("Blackbox Cookie").(*tview.InputField).SetText("")
	})
	pages.AddPage(fmt.Sprintf("settings_%d", 1), form, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
