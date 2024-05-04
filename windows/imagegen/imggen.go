package imagegen

import (
	"context"
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"guiv1/models/sydney"
	"path/filepath"
	"runtime"

	"guiv1/models/util"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	wx "fyne.io/x/fyne/widget"
)

var (
	icns = misc.IconUtil{}
	iu   = misc.ImageUtil{}
	f_   = misc.Funcs{}
)

func GenImageBing(image_gen_prompt string) map[string]interface{} {
	if image_gen_prompt == "" {
		return map[string]interface{}{
			"success": false,
			"error":   "No image generation prompt provided",
		}
	}
	cookies, err := util.ReadCookiesFile()
	if err != nil {
		log.Fatalf("Error reading cookies file: %v", err)
	}

	sydneyAPI := sydney.NewSydney(sydney.Options{

		Cookies: cookies,

		ConversationStyle: "creative",
		Locale:            "en-US",
	})

	messageCh, err := sydneyAPI.AskStream(sydney.AskStreamOptions{
		StopCtx: context.TODO(),
		Prompt:  "Create image for the description: " + image_gen_prompt,
	})
	if err != nil {
		log.Fatalf("Error creating Sydney instance: %v", err)
	}

	var generativeImage sydney.GenerativeImage

	for message := range messageCh {
		if message.Type == sydney.MessageTypeGenerativeImage {
			err := json.Unmarshal([]byte(message.Text), &generativeImage)
			if err == nil {
				break
			}
		}
	}
	if generativeImage.URL == "" {
		log.Println("No image URL returned from the API")
		return map[string]interface{}{
			"success": false,
			"error":   "No image URL returned from the API",
		}
	}

	image, err := sydneyAPI.GenerateImage(generativeImage)
	if err != nil {
		log.Fatalf("Error generating image: %v", err)
	}
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	fmt.Println(currentDir)
	fmt.Println(image.Duration)
	if len(image.ImageURLs) > 0 {
		timestamp := time.Now().Format("2006-01-02-15-04-05")
		id_ := 0
		for _, url := range image.ImageURLs {
			id_ += 1
			urlParts := strings.Split(url, "?")
			url = urlParts[0]

			filename := fmt.Sprintf("generated_image_%s_%s.png", timestamp, strconv.Itoa(id_))
			fmt.Println("Image URL:", url)
			fmt.Println("Filename:", filename)
			return map[string]interface{}{
				"success":    true,
				"image_urls": image.ImageURLs,
			}

		}
	}
	return map[string]interface{}{
		"success": false,
		"error":   "No image URL returned from the API",
	}
}

func ImageGenerationWindow(a fyne.App, w fyne.Window) {
	infinite_progress := widget.NewProgressBarInfinite()
	gif, err := wx.NewAnimatedGif(storage.NewFileURI("./assets/loading.gif"))
	gif.SetMinSize(fyne.Size{Width: 30, Height: 30})
	gif.Resize(fyne.Size{Width: 30, Height: 30})

	gif_ := container.NewHBox(gif)
	gif_.Resize(fyne.Size{Width: 30, Height: 30})
	gif_.Hidden = true
	PromptEntry := widget.NewEntry()
	PromptEntry.PlaceHolder = "Enter prompt here"
	PromptSubmit := widget.NewButton("Generate", nil)
	PromptContainer := container.NewGridWithColumns(
		2,
		PromptEntry,
		PromptSubmit,
	)
	Grid := container.New(layout.NewGridLayout(2))
	title_label := widget.NewRichTextFromMarkdown("# Image Generation ( Bing Creative )" + misc.InvisFill + misc.InvisFill)
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
	SaveMenu := container.NewHBox()
	SaveMenu.Layout = layout.NewGridLayout(4)
	Images_ := container.NewBorder(
		TopBorder,
		container.NewVBox(
			PromptContainer,
			infinite_progress,
		),
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
		SaveMenu,
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
		SaveMenu.RemoveAll()
		Images_.Refresh()
		Gen := GenImageBing(text)
		Grid.RemoveAll()

		infinite_progress.Hidden = true
		if Gen["success"].(bool) {
			image_urls := Gen["image_urls"].([]string)
			x := 0
			SaveMenu.RemoveAll()
			for _, url := range image_urls {
				x += 1
				Grid.Add(iu.NewCanvasImageUri(350, 350, url))
				SaveMenu.Add(widget.NewButton("Image "+strconv.Itoa(x), func() {
					timestamp := time.Now().Format("2006-01-02-15-04-05")
					os := runtime.GOOS
					if os == "windows" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					if os == "darwin" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					if os == "linux" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					f_.NotificationModal(w, &misc.ChatApp{}, "Image Saved", fmt.Sprintf("Image %d saved to 'data/generated_images' directory.", x))
				}))

				SaveMenu.Refresh()
				Images_.Refresh()
			}
			Grid.Refresh()
			Images_.Refresh()
		} else {
			err := Gen["error"].(string)
			log.Println("Error generating image:", err)
			f_.NotificationModal(w, &misc.ChatApp{}, "Error Generating Image", err)
		}
	}
	PromptEntry.OnSubmitted = func(text string) {
		infinite_progress.Hidden = false

		if err != nil {
			log.Fatalf("Error creating animated GIF: %v", err)
		}
		Grid.RemoveAll()
		SaveMenu.RemoveAll()
		Images_.Refresh()
		PromptEntry.Text = ""
		PromptEntry.Refresh()
		Gen := GenImageBing(text)
		Grid.RemoveAll()

		infinite_progress.Hidden = true
		if Gen["success"].(bool) {
			SaveMenu.RemoveAll()
			image_urls := Gen["image_urls"].([]string)
			x := 0
			for _, url := range image_urls {
				x += 1
				Grid.Add(iu.NewCanvasImageUri(350, 350, url))
				SaveMenu.Add(widget.NewButton("Image "+strconv.Itoa(x), func() {
					timestamp := time.Now().Format("2006-01-02-15-04-05")
					os := runtime.GOOS
					if os == "windows" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					if os == "darwin" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					if os == "linux" {
						err = f_.DownloadImage(url, filepath.Join("./", "data", "generated_images", "generated_image_"+timestamp+"_"+strconv.Itoa(x)+".png"))
					}
					f_.NotificationModal(w, &misc.ChatApp{}, "Image Saved", "Image saved to 'data/generated_images' directory.")
				}))

				SaveMenu.Refresh()
				Images_.Refresh()
			}
			Grid.Refresh()
			Images_.Refresh()
		} else {
			err := Gen["error"].(string)
			log.Println("Error generating image:", err)
			f_.NotificationModal(w, &misc.ChatApp{}, "Error", err)
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
