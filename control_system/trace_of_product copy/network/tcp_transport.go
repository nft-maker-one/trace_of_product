package network

import (
	"agricultural_meta/utils"
	"bufio"
	"fmt"
	"net"
	"strings"

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

func (p *NodeServer) StartServer() {
	listen, err := net.Listen("tcp", p.Addr)
	if err != nil {
		panic(err)
	}
	if p.IsLeader {
		go p.CreateBlock()
	}
	utils.LogMsg([]string{"NodeUp"}, []string{"Server started successfully"})

	// for {
	// 	conn, err := listen.Accept()
	// 	if err != nil {
	// 		utils.LogError([]string{"NodeUp"}, []string{err.Error()})
	// 	}
	// 	data, err := io.ReadAll(conn)
	// 	if err != nil {
	// 		utils.LogError([]string{"NodeUp"}, []string{err.Error()})
	// 	}
	// 	p.handleRequest(data)
	// 	conn.Close()
	// }
	for {
		conn, err := listen.Accept()
		if err != nil {
			utils.LogError([]string{"NodeUp"}, []string{err.Error()})
		}
		method, path := handleConnection(conn)
		if method == "GET" {
			if path == "hello" {
				continue
			}
		} else if method == "POST" {

		}
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

func parseRequest(request string) (method, path string) {

	parts := strings.Split(request, " ")
	fmt.Println(parts)
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
