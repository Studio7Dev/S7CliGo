package rathandler

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"fyne.io/fyne/v2/widget"
)

var (
	ClientsConnected int
)
var ConsoleLog = widget.NewMultiLineEntry()
var StringClients []string
var Clients []net.Conn
var Tasks []string
var FileUploadInProgress = false

type ClientReply struct {
	Message string `json:"message"`
}

type TCPServer struct {
	listener net.Listener
	running  bool
	Addr     string
}

func NewTCPServer() *TCPServer {
	return &TCPServer{}
}

func (s *TCPServer) WriteLog(msg string) {
	TimeStamp := fmt.Sprintf("[%s] ", time.Now().Format("15:04:05"))
	BaseString := fmt.Sprintf("%s%s", TimeStamp, msg)
	ConsoleLog.SetText(ConsoleLog.Text + BaseString + "\n")
}

func (s *TCPServer) PingClient(i int, conn net.Conn) error {
	_, err := conn.Write([]byte("ping"))
	if err != nil {
		Clients = append(Clients[:i], Clients[i+1:]...)
		StringClients = append(StringClients[:i], StringClients[i+1:]...)
		conn.Close()
		s.WriteLog(fmt.Sprintf("Error pinging client: %v", err))
		return err
	}
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	msg := string(buffer[:n])
	if err != nil {

		return err
	}
	if msg == "pong" {
		return nil
	} else {
		if FileUploadInProgress {
			s.WriteLog("File upload in progress, cannot process new tasks")

		} else {

			var ClientReply ClientReply
			err = json.Unmarshal([]byte(msg), &ClientReply)
			if err != nil {
				s.WriteLog(fmt.Sprintf("Unexpected response from client %s: %s", conn.RemoteAddr().String(), msg))
				if strings.Contains(err.Error(), "invalid character") {
					s.WriteLog(fmt.Sprintf("Error unmarshaling client response: %v", err))
				} else {
					s.WriteLog("Disconnecting client: " + conn.RemoteAddr().String())
					conn.Close()
					Clients = append(Clients[:i], Clients[i+1:]...)
					StringClients = append(StringClients[:i], StringClients[i+1:]...)
				}

			}
			fmt.Println(ClientReply)
			//s.WriteLog(fmt.Sprintf("Received response from client %s: %s", conn.RemoteAddr().String(), ClientReply.Message))

		}
	}
	return err
}

func (s *TCPServer) MonitorClients() {
	for i, client := range Clients {

		CurrentPingClientAddr := StringClients[i]
		PingClient_ := s.PingClient(i, client)
		if PingClient_ != nil {
			s.WriteLog(fmt.Sprintf("Error pinging client %s: %v", CurrentPingClientAddr, PingClient_))
			break
		}
	}
}

func (s *TCPServer) Start() error {
	s.WriteLog(fmt.Sprintf("Starting TCP server on %s", s.Addr))
	go func() {
		for {
			s.MonitorClients()
			time.Sleep(5 * time.Second)
			if !s.running {
				break
			}
		}
	}()
	if s.running {
		s.WriteLog("Server is already running")
		return fmt.Errorf("server is already running")
	}
	ln, err := net.Listen("tcp", s.Addr)
	s.WriteLog(fmt.Sprintf("Server Listening on %s", s.Addr))
	if err != nil {
		return err
	}
	s.listener = ln
	s.running = true
	go s.handleConnections()
	return nil
}

func (s *TCPServer) Stop() error {
	for client := range Clients {
		Clients[client].Close()
	}
	Clients = []net.Conn{}
	StringClients = []string{}
	if !s.running {
		s.WriteLog("Server is not running, there are no clients to disconnect.")
		return fmt.Errorf("server is not running")
	}
	s.running = false
	s.WriteLog("Server stopped")
	return s.listener.Close()
}

func (s *TCPServer) handleConnections() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.running {
				break
			}
			fmt.Println(err)
			continue
		}
		if !FileUploadInProgress {
			Clients = append(Clients, conn)
			StringClients = append(StringClients, conn.RemoteAddr().String())
		}

	}
}

var ClientsTaskCompletedNo []string

func (s *TCPServer) TaskWatcher() {
	for {
		if len(Clients) > 0 {
			if len(Clients) == len(ClientsTaskCompletedNo) {
				s.WriteLog("All clients have completed their tasks.")
				ClientsTaskCompletedNo = []string{}
				Tasks = []string{}
			}
		}
	}
}
