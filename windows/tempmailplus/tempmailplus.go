package tempmailplus

import (
	"fmt"
	"time"

	"guiv1/misc"
	"guiv1/parsers/verif"

	"guiv1/handlers/tempmailplus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"golang.design/x/clipboard"
)

var (
	icns = misc.IconUtil{}
	iu   = misc.ImageUtil{}
	f_   = misc.Funcs{}
)

func EmailElement(message tempmailplus.Email) *widget.Card {
	subject := widget.NewLabel(message.Subject)
	subject.Wrapping = fyne.TextWrapWord
	from := widget.NewLabel(fmt.Sprintf("From: %s", message.From))
	date := widget.NewLabel(message.Date)
	converter := md.NewConverter("", true, nil)

	html := message.Html

	markdown, err := converter.ConvertString(html)
	if err != nil {
		panic(err)
	}
	fmt.Println("=====================")
	fmt.Println(len(markdown))
	fmt.Println("=====================")

	content := widget.NewRichTextFromMarkdown(markdown)
	content.Wrapping = fyne.TextWrapWord
	WrapperContainer := widget.NewCard(
		"Email",
		"",
		container.NewVBox(
			subject,
			from,
			date,
			content,
		),
	)
	WrapperContainer.Resize(fyne.NewSize(600, 150))
	return WrapperContainer
}

func RaidWindow(a fyne.App, w fyne.Window) {

	EmailClient := tempmailplus.TmailPlusClient{}

	Account := EmailClient.NewAccount()
	EmailClient.Domain = Account.Domain
	EmailClient.Username = Account.Username
	Email := Account.Email

	EmailAccountDetails := container.NewVBox(
		widget.NewRichTextFromMarkdown(fmt.Sprintf(`
		# Temp Mail Account Details
		**Email:** %s
		`, Email)),
		widget.NewSeparator(),
		widget.NewButton("Copy Email", func() {
			clipboard.Write(clipboard.FmtText, []byte(Email))
			f_.NotificationModal(w, &misc.ChatApp{}, "Email Copied", "The email address has been copied to the clipboard.")
		}),
	)

	EmailMessages := container.NewVBox()
	EmailMessages.Layout = layout.NewVBoxLayout()
	// EmailMessages.Layout = layout.NewStackLayout()
	MessagesScroll := container.NewScroll(EmailMessages)
	MessagesScroll.Resize(fyne.NewSize(600, 400))
	MessagesScroll.SetMinSize(fyne.NewSize(600, 400))

	MainContainer := container.NewVBox(
		EmailAccountDetails,
		widget.NewSeparator(),
		MessagesScroll,
	)

	content := container.NewBorder(
		MainContainer,
		nil,
		nil,
		nil,
		nil,
	)
	modal := widget.NewModalPopUp(content, w.Canvas())
	CloseBtn := widget.NewButton("Close", func() {
		modal.Hide()
	})
	CloseBtn.OnTapped = func() {
		modal.Hide()
	}
	EmailAccountDetails.Add(CloseBtn)
	content.Refresh()
	modal.Resize(fyne.NewSize(700, 600))
	modal.Canvas.SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEscape {
			modal.Hide()
		}
	})
	RefreshBtn := widget.NewButton("Refresh", func() {

		messages := EmailClient.CheckInbox().MailList

		time.Sleep(3 * time.Second)
		if len(messages) > 0 {
			LatestMessage := messages[0]
			DetailedMessage := EmailClient.GetMessageByID(LatestMessage.MailId)
			NewEmailElement := EmailElement(DetailedMessage)
			// check if the email element is already in the EmailMessages container
			EmailMessages.Add(NewEmailElement)

			sep := widget.NewLabel("")
			EmailMessages.Add(sep)
			EmailMessages.Refresh()
			MessagesScroll.Refresh()
			MainContainer.Refresh()

		}

	})

	EmailAccountDetails.Add(RefreshBtn)
	go func() {
		for {
			messages := EmailClient.CheckInbox().MailList

			time.Sleep(2 * time.Second)
			if len(messages) > 0 {
				LatestMessage := messages[0]
				DetailedMessage := EmailClient.GetMessageByID(LatestMessage.MailId)
				NewEmailElement := EmailElement(DetailedMessage)
				// check if the email element is already in the EmailMessages container
				EmailMessages.Add(NewEmailElement)
				// extraction notification
				f_.NotificationModal(w, &misc.ChatApp{}, "Info", "Extracting verification information from the latest email...")
				verifObj, err := verif.ExtractVerificationInfo(DetailedMessage.Html)
				if err != nil {
					verifObj, err = verif.ExtractVerificationInfo(DetailedMessage.Text)
					if err != nil {
						f_.NotificationModal(w, &misc.ChatApp{}, "Error", "Failed to extract verification information from email")
						break
						return
					}
				}
				Info := fmt.Sprintf(`
				### VerificationCode: %s
				### AccountUsername: %s
				### Instructions: %s
				### ExpirationTime: %s
				### Verification Link: %s`, verifObj.VerificationCode, verifObj.AccountUsername, verifObj.Instructions, verifObj.ExpirationTime, verifObj.VerificationURL)
				clipboard.Write(clipboard.FmtText, []byte(Info))
				f_.CNotificationModal(w, &misc.ChatApp{}, "Verification Info Copied to Clipboard", Info)
				sep := widget.NewLabel("")
				EmailMessages.Add(sep)
				EmailMessages.Refresh()
				MessagesScroll.Refresh()
				MainContainer.Refresh()
				break
			}
		}

	}()

	modal.Show()

}
