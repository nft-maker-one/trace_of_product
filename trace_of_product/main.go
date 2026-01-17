package main

import (
	"agricultural_meta/database"
	"agricultural_meta/network"
	"agricultural_meta/utils"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
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

	// 初始化日志
	if err := utils.InitLogger("./logs"); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}

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

		// 启动HTTP服务器用于日志监控
		httpPort := 8080
		go startHTTPServer(httpPort)
		fmt.Printf("HTTP server started on port %d\n", httpPort)
	} else {
		fmt.Println("请务必指定节点的运行模式")
		fmt.Println("输入 '-h' 查看帮助指令")
	}
}

func startHTTPServer(port int) {
	http.HandleFunc("/logs", handleLogs)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	limit := 100 // 默认返回100条日志
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs := utils.GetLogs(limit)
	json.NewEncoder(w).Encode(logs)
}
