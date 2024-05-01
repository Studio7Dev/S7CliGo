package mailraid

import (
	"fmt"

	"guiv1/misc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	icns = misc.IconUtil{}
)

func RaidWindow(a fyne.App, w fyne.Window) {
	// new form
	form_ := widget.NewForm(
		&widget.FormItem{
			Widget: widget.NewLabel("Mail Raid"),
		},
		&widget.FormItem{
			Widget: canvas.NewImageFromResource(icns.Icons8("256", "bomb-emoji.png", "emoji")),
		},
		&widget.FormItem{
			Text:   "Target Email",
			Widget: widget.NewEntry(),
		},
	)
	// Add a submit button
	title_ := form_.Items[0].Widget.(*widget.Label)
	title_.TextStyle.Bold = true
	title_.Alignment = fyne.TextAlignCenter
	title_.Resize(fyne.NewSize(300, 100))
	BombImage := form_.Items[1].Widget.(*canvas.Image)
	BombImage.Resize(fyne.NewSize(256, 256))
	BombImage.FillMode = canvas.ImageFillOriginal
	submitButton := widget.NewButton("Submit", func() {
		// Get the values from the form
		targetEmail := form_.Items[2].Widget.(*widget.Entry).Text
		fmt.Println(targetEmail)
	})
	// Add the form and submit button to the window

	content := container.NewVBox(
		form_,
		submitButton,
	)
	modal := widget.NewModalPopUp(content, w.Canvas())
	CloseBtn := widget.NewButton("Close", func() {
		modal.Hide()
	})
	CloseBtn.OnTapped = func() {
		modal.Hide()
	}
	content.Add(CloseBtn)
	content.Refresh()
	modal.Resize(fyne.NewSize(700, 600))
	modal.Show()

}
