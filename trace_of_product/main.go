package main

import (
	"agricultural_meta/database"
	"agricultural_meta/network"
	"flag"
	"fmt"
)

var DB = &database.NodeDb{}

func main() {

	ip := ""
	port := ""
	clientMode := false
	chainMode := false
	help := false
	flag.StringVar(&ip, "ip", "127.0.0.1", "please set the ip for node")
	flag.StringVar(&port, "port", "8081", "please set the port for node")
	flag.BoolVar(&clientMode, "client", false, "set the node with clientMode")
	flag.BoolVar(&chainMode, "chain", false, "set the node with chainMode")
	flag.BoolVar(&help, "h", false, "note for inputs")
	flag.Parse()
	netAddr := ip + ":" + port
	fmt.Println(netAddr)
	if help {
		fmt.Println("农产追溯通使用详解")
		fmt.Println("-----------------------------")
		fmt.Println("1. -ip 指定节点运行的 ip 地址")
		fmt.Println("2. -port 指定节点运行的端口")
		fmt.Println("3. -client", "运行 client 节点")
		fmt.Println("4 -chain", "运行联盟链节点")
		fmt.Println("-----------------------------")
	} else if clientMode {
		client := network.NewClient(netAddr)
		client.ClientUp()
	} else if chainMode {
		chain := network.NewNodeServer(netAddr)
		chain.StartServer()
	} else {
		fmt.Println("请务必指定节点的运行模式")
		fmt.Println("输入 '-h' 查看帮助指令")
	}
}

// encode => io.Writer  把编码后的数据写入 io.Writer
// decode => io.Reader  把解码后的数据读入 io.Reader
