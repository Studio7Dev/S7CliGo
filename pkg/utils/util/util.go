package util

import (
	"CLI/pkg/misc"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/imroc/req/v3"

	"github.com/samber/lo"

	"github.com/ncruces/zenity"
	getproxy "github.com/rapid7/go-get-proxied/proxy"
)

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
)

func RandIntInclusive(min int, max int) int {
	return min + rand.Intn(max-min+1)
}
func Ternary[T any](expression bool, trueResult T, falseResult T) T {
	if expression {
		return trueResult
	} else {
		return falseResult
	}
}
func MakeHTTPClient(proxy string, timeout time.Duration) (*http.Client, *req.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	reqClient := req.C().ImpersonateChrome().SetCommonHeader("accept-language", "en-US,en;q=0.9")
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, nil, err
		}
		transport.Proxy = http.ProxyURL(proxyURL)
		reqClient.SetProxyURL(proxy)
	} else {
		log.SetOutput(io.Discard)
		proxies := []getproxy.Proxy{
			getproxy.NewProvider("").GetHTTPProxy("https://www.bing.com"),
			getproxy.NewProvider("").GetHTTPSProxy("https://www.bing.com"),
			getproxy.NewProvider("").GetSOCKSProxy("https://www.bing.com"),
		}
		log.SetOutput(os.Stdout)
		var sysProxy getproxy.Proxy
		for _, p := range proxies {
			p := p
			if p != nil {
				sysProxy = p
				break
			}
		}
		if sysProxy != nil {
			transport.Proxy = http.ProxyURL(sysProxy.URL())
			reqClient.SetProxyURL(sysProxy.URL().String())
		}
	}
	client := &http.Client{}
	client.Transport = transport
	if timeout != time.Duration(0) {
		client.Timeout = timeout
		reqClient.SetTimeout(timeout)
	}
	return client, reqClient, nil
}
func FormatCookieString(cookies map[string]string) string {
	str := ""
	for k, v := range cookies {
		str += k + "=" + v + "; "
	}
	return str
}
func ParseCookiesFromString(cookiesStr string) map[string]string {
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
func CopyMap[T comparable, E any](source map[T]E) map[T]E {
	res := map[T]E{}
	for k, v := range source {
		res[k] = v
	}
	return res
}
func CreateTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
func CreateCancelContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
func MustGenerateRandomHex(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.New(rand.NewSource(time.Now().Unix())).Read(randomBytes)
	if err != nil {
		GracefulPanic(err)
	}
	randomString := hex.EncodeToString(randomBytes)
	return randomString
}

type FileCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func ReadCookiesFileRaw() ([]FileCookie, error) {
	v, err := os.ReadFile(WithPath(settings.BingCookie))
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
func ReadCookiesFile() (map[string]string, error) {
	res := map[string]string{}
	cookies, err := ReadCookiesFileRaw()
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}
	return res, nil
}
func UpdateCookiesFile(cookies map[string]string) error {
	var arr []FileCookie
	for k, v := range cookies {
		arr = append(arr, FileCookie{
			Name:  k,
			Value: v,
		})
	}
	v, err := json.MarshalIndent(&arr, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(WithPath(settings.BingCookie), v, 0644)
	if err != nil {
		return err
	}
	return nil
}
func Map[T any, E any](arr []T, function func(value T) E) []E {
	var result []E
	for _, item := range arr {
		result = append(result, function(item))
	}
	return result
}
func FindFirst[T any](arr []T, function func(value T) bool) (T, bool) {
	var empty T
	for _, item := range arr {
		if function(item) {
			return item, true
		}
	}
	return empty, false
}
func ConvertImageToJpg(img []byte) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, src, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func GenerateSecMSGec() string {
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	randomBytes := make([]byte, 32)

	for i := range randomBytes {
		randomBytes[i] = byte(rng.Intn(256))
	}

	return hex.EncodeToString(randomBytes)
}
func GracefulPanic(err error) {
	_, file, line, _ := runtime.Caller(1)
	zenity.Error(fmt.Sprintf("Error: %v\nDetails: file(%s), line(%d).\n"+
		"Instruction: This is probably an unknown bug. Please take a screenshot and report this issue.",
		err, file, line))
	lo.Must0(OpenURL("https://github.com/juzeon/SydneyQt/issues"))
	os.Exit(-1)
}
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
func ReadDebugOptionSets() (debugOptionsSets []string) {
	debugOptionsSetsFile, err := os.ReadFile(WithPath("debug_options_sets.json"))
	if err != nil {
		return
	}
	if strings.TrimSpace(string(debugOptionsSetsFile)) == "" {
		return
	}
	err = json.Unmarshal(debugOptionsSetsFile, &debugOptionsSets)
	if err != nil {
		GracefulPanic(err)
	}
	if len(debugOptionsSets) != 0 {
	}
	return
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

func WithPath(filename string) string {
	initWithPath()
	if withPathBaseDir == "" {
		return filename
	}
	return filepath.Join(withPathBaseDir, filename)
}
