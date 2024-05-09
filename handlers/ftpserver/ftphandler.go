package ftpserver

import (
	"log"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

func NewFtpServer(root, user, pass, host string, port int) {
	factory := &filedriver.FileDriverFactory{
		RootPath: root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     port,
		Hostname: host,
		Auth:     &server.SimpleAuth{Name: user, Password: pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", user, pass)
	server := server.NewServer(opts)
	server.ListenAndServe()
}

func RunFtp() {
	go NewFtpServer(
		"./data/clientuploads",
		"uploads",
		"uploads",
		"0.0.0.0",
		2121,
	)
	go NewFtpServer(
		"./data/serverdownloads",
		"downloads",
		"downloads",
		"0.0.0.0",
		2122,
	)
	go NewFtpServer(
		"./data/screenshots",
		"screenshots",
		"screenshots",
		"0.0.0.0",
		2123,
	)
}
