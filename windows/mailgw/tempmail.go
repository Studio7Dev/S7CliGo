package mailgw

import (
	"fmt"
	"time"

	"guiv1/misc"
	"guiv1/parsers/verif"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/felixstrobel/mailtm"
	"golang.design/x/clipboard"
)

var (
	icns = misc.IconUtil{}
	iu   = misc.ImageUtil{}
	f_   = misc.Funcs{}
)

func GetMessage(client *mailtm.MailClient, Account *mailtm.Account, Id string) mailtm.DetailedMessage {
	DetailedMessage, err := client.GetMessageByID(Account, Id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Detailed Message:", DetailedMessage)
	return *DetailedMessage
}

func NewAccount(client *mailtm.MailClient) (*mailtm.Account, error) {
	account, err := client.NewAccount()
	if err != nil {
		return nil, err
	}
	fmt.Println("New email account created:", account.Address)
	fmt.Println("Account Password:", account.Password)
	fmt.Println("Account ID:", account.ID)
	return account, nil
}

func EmailElement(message mailtm.DetailedMessage) *widget.Card {
	subject := widget.NewLabel(message.Subject)
	subject.Wrapping = fyne.TextWrapWord
	from := widget.NewLabel(fmt.Sprintf("From: %s", message.From))
	date := widget.NewLabel(message.CreatedAt.Format(time.RFC822))
	converter := md.NewConverter("", true, nil)

	html := message.Html[0]

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

	EmailClient, err := mailtm.New()
	if err != nil {
		panic(err)
	}
	Account, err := EmailClient.NewAccount()
	if err != nil {
		panic(err)
	}
	Email := Account.Address
	ID := Account.ID
	Password := Account.Password

	EmailAccountDetails := container.NewVBox(
		widget.NewRichTextFromMarkdown(fmt.Sprintf(`
		# Temp Mail Account Details
		**Email:** %s
		**ID:** %s
		**Password:** %s
		`, Email, ID, Password)),
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
	RetrievedAccount, err := EmailClient.RetrieveAccount(Email, Password)
	RefreshBtn := widget.NewButton("Refresh", func() {

		messages, err := EmailClient.GetMessages(RetrievedAccount, 1)
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
		if len(messages) > 0 {
			LatestMessage := messages[0]
			DetailedMessage := GetMessage(EmailClient, RetrievedAccount, LatestMessage.ID)
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
			messages, err := EmailClient.GetMessages(RetrievedAccount, 1)
			if err != nil {
				panic(err)
			}
			time.Sleep(2 * time.Second)
			if len(messages) > 0 {
				LatestMessage := messages[0]
				DetailedMessage := GetMessage(EmailClient, RetrievedAccount, LatestMessage.ID)
				NewEmailElement := EmailElement(DetailedMessage)
				// check if the email element is already in the EmailMessages container
				EmailMessages.Add(NewEmailElement)
				// extraction notification
				f_.NotificationModal(w, &misc.ChatApp{}, "Info", "Extracting verification information from the latest email...")
				verifObj, err := verif.ExtractVerificationInfo(DetailedMessage.Html[0])
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
