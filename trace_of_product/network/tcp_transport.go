package network

import (
	"agricultural_meta/utils"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

func tcpDial(data []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to connect to %s", addr)})
		logrus.Errorln(err.Error())
	}
	_, err = conn.Write(data)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to write data to %s", addr)})
		logrus.Errorln(err.Error())
	}
	conn.Close()
}

func (p *NodeServer) NodeUp() {
	listen, err := net.Listen("tcp", p.Addr)
	if err != nil {
		panic(err)
	}
	utils.LogMsg([]string{"NodeUp"}, []string{"Server started successfully"})
	for {
		conn, err := listen.Accept()
		if err != nil {
			utils.LogError([]string{"NodeUp"}, []string{err.Error()})
		}
		data, err := io.ReadAll(conn)
		if err != nil {
			utils.LogError([]string{"NodeUp"}, []string{err.Error()})
		}
		p.handleRequest(data)
	}
}
