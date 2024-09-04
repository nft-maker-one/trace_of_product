package network

import (
	"agricultural_meta/utils"
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

func tcpDial(data []byte, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to connect to %s", addr)})
		logrus.Errorln(err.Error())
		return err
	}
	_, err = conn.Write(data)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to write data to %s", addr)})
		logrus.Errorln(err.Error())
		return err
	}
	conn.Close()
	return nil
}

func (p *NodeServer) StartServer() {
	listen, err := net.Listen("tcp", p.Addr)
	if err != nil {
		panic(err)
	}
	if p.IsLeader {
		go p.CreateBlock()
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
		conn.Close()
	}
}

func handleConnection(conn net.Conn) (string, string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	request, _ := reader.ReadString('\n')
	for s, err := reader.ReadString('\n'); err != nil; {
		fmt.Println(s)
	}
	res := strings.Split(request, " ")
	return res[0], res[1]

}
