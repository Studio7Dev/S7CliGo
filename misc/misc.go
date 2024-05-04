package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type Init_ struct {
}

type Funcs struct {
}
type Data struct {
	DarkMode           bool   `json:"dark_mode"`
	BingHost           string `json:"bing_host"`
	CurrentTuneModel   string `json:"current_tune_model"`
	CurrentHugModel    string `json:"current_hug_model"`
	GoliathAuthToken   string `json:"goliath_auth_token"`
	MerlinAuthToken    string `json:"merlin_auth_token"`
	TuneAppAccessToken string `json:"tuneapp_auth_token"`
	HugginFaceCookie   string `json:"huggingface_cookie"`
	BingCookie         string `json:"bing_cookie"`
	BlackBoxCookie     string `json:"blackbox_cookie"`
	YouAICookie        string `json:"youai_cookie"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	TcpHost            string `json:"tcphost"`
	Httphost           string `json:"httphost"`
}

var InvisFill = "‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎ ‎"

func (f Funcs) LoadSettings() (Data, error) {
	settingsFile, err := os.Open("settings.json")
	if err != nil {
		return Data{}, fmt.Errorf("error opening JSON file: %v", err)
	}

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

func (f Funcs) DownloadImage(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create(filepath)
	defer out.Close()
	opts := &jpeg.Options{
		Quality: 100,
	}
	// save as png

	err = jpeg.Encode(out, img, opts)
	if err != nil {
		return err
	}

	return nil
}

type CookieUtil struct{}
type FileCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var initWithPath = sync.OnceFunc(func() {
	if runtime.GOOS == "darwin" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return
		}
		full := filepath.Join(dir, "SydneyQt")
		err = os.MkdirAll(full, 0750)
		if err != nil {
			return
		}
		withPathBaseDir = full
	}
})
var withPathBaseDir string

func (cu CookieUtil) GetCValue(name string, cookies []FileCookie) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""

}

func (cu CookieUtil) FormatCookieString(cookies map[string]string) string {
	str := ""
	for k, v := range cookies {
		str += k + "=" + v + "; "
	}
	return str
}

func (cu CookieUtil) WithPath(filename string) string {
	initWithPath()
	if withPathBaseDir == "" {
		return filename
	}
	return filepath.Join(withPathBaseDir, filename)
}

func (cu CookieUtil) ReadCookiesFileRaw(path_ string) ([]FileCookie, error) {
	v, err := os.ReadFile(cu.WithPath(path_))
	if err != nil {
		return nil, nil
	}
	var cookies []FileCookie
	err = json.Unmarshal(v, &cookies)
	if err != nil {
		return nil, errors.New("failed to json.Unmarshal content of cookie file")
	}
	return cookies, nil
}
func (cu CookieUtil) ReadCookiesFile(path_ string) string {
	res := map[string]string{}
	cookies, err := cu.ReadCookiesFileRaw(path_)
	if err != nil {
		return ""
	}
	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}
	return cu.FormatCookieString(res)
}

func (cu CookieUtil) ParseCookiesFromString(cookiesStr string) map[string]string {
	cookies := map[string]string{}
	for _, cookie := range strings.Split(cookiesStr, ";") {
		cookie = strings.TrimSpace(cookie)
		parts := strings.Split(cookie, "=")
		if len(parts) < 2 {
			continue
		}
		cookies[parts[0]] = strings.Join(parts[1:], "=")
	}
	return cookies
}

type IconUtil struct{}

func (iu IconUtil) IconFromBytes(IconName string, IconBytes []byte) fyne.Resource {
	return fyne.NewStaticResource(IconName, IconBytes)
}

func (iu IconUtil) IconFromRepo(name string) fyne.Resource {
	IconsRepoBaseUrl := "https://raw.githubusercontent.com/ZachC137/MICNS/e2b386038b1465856a3f46735a384b80d5511655/"
	// http get request to the repo to get the icon
	resp, err := http.Get(IconsRepoBaseUrl + name + ".png")
	if err != nil {
		log.Fatalf("Failed to fetch icon: %v", err)
	}
	IconBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read icon response body: %v", err)
	}
	return fyne.NewStaticResource(name, IconBytes)
}

func (iu IconUtil) IconByteLoader(IconName string, IconsFolder string) []byte {
	if IconsFolder == "" {
		IconsFolder = "assets/"
	}
	IconBytes, err := ioutil.ReadFile(IconsFolder + IconName + ".png")
	if err != nil {
		log.Fatalf("Failed to load icon: %v", err)
	}
	return IconBytes
}

func (iu IconUtil) Icon(name string) fyne.Resource {
	return fyne.NewStaticResource(name, iu.IconByteLoader(name, ""))
}

func (iu IconUtil) Icons8(uuid string, name string, category string) fyne.Resource {
	if category == "" {
		category = "color"
	}
	BaseUrl := fmt.Sprintf("https://img.icons8.com/%s/%s/%s", category, uuid, name)
	resp, err := http.Get(BaseUrl)
	if err != nil {
		log.Fatalf("Failed to fetch icon: %v", err)
	}
	IconBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read icon response body: %v", err)
	}
	return fyne.NewStaticResource(name, IconBytes)
}

type ImageUtil struct{}

func NewUuid() string {
	uuid_base, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uuid_base.String()
}

func (iu ImageUtil) LoadImageFromBytes(name string, data []byte) fyne.Resource {
	return fyne.NewStaticResource(name, data)
}
func (iu ImageUtil) LoadImageFromUri(name string, uri string) fyne.Resource {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("Failed to fetch image: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read image response body: %v", err)
	}
	return iu.LoadImageFromBytes(name, data)
}
func (iu ImageUtil) NewCanvasImageUri(w float32, h float32, uri string) *fyne.Container {
	Image_ := canvas.NewImageFromResource(iu.LoadImageFromUri(NewUuid(), uri))
	// Image_.FillMode = canvas.ImageFillContain
	Image_.SetMinSize(fyne.Size{Width: w, Height: h})
	Image_.Resize(fyne.Size{Width: w, Height: h})
	view := container.NewHBox(Image_)
	return view
}

func (iu ImageUtil) NewCanvasImageFile(w float32, h float32, filePath string) *fyne.Container {
	Image_ := canvas.NewImageFromFile(filePath)
	Image_.SetMinSize(fyne.Size{Width: w, Height: h})
	Image_.Resize(fyne.Size{Width: w, Height: h})
	view := container.NewHBox(Image_)
	return view
}

type ChatApp struct {
	App     fyne.App
	Win     fyne.Window
	Input   *widget.Entry
	ChatLog *widget.Entry
}

func (f Funcs) NotificationModal(w fyne.Window, a *ChatApp, title string, message string) {
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
	OKBtn.SetIcon(IconUtil{}.Icons8("256", "checkmark--v1.png", ""))
	OKBtn.OnTapped = func() {
		popup.Hide()
	}
	popup.Resize(fyne.NewSize(300, 150))
	popup.Show()
}
