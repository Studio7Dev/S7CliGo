package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

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
