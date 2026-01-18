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
var GlobalChain *network.NodeServer

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
		GlobalChain = chain

		// 启动HTTP服务器用于日志监控
		httpPort := 8080
		go startHTTPServer(httpPort)
		fmt.Printf("HTTP server started on port %d\n", httpPort)

		// 启动区块链TCP服务器 (这个是阻塞的，所以放在最后)
		chain.StartServer()
	} else {
		fmt.Println("请务必指定节点的运行模式")
		fmt.Println("输入 '-h' 查看帮助指令")
	}
}

func startHTTPServer(port int) {
	http.HandleFunc("/logs", handleLogs)
	http.HandleFunc("/api/height", handleGetHeight)
	http.HandleFunc("/api/blocks", handleGetBlocksByRange)
	http.HandleFunc("/api/node/status", handleGetNodeStatus)
	http.HandleFunc("/api/block/transactions", handleGetBlockTransactions)
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

// 获取最高区块高度
func handleGetHeight(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	if GlobalChain == nil {
		http.Error(w, "Chain not initialized", http.StatusInternalServerError)
		return
	}

	height := GlobalChain.Chain.Height()
	response := map[string]interface{}{
		"height":  height,
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 根据区块范围筛选数据
func handleGetBlocksByRange(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	if GlobalChain == nil {
		http.Error(w, "Chain not initialized", http.StatusInternalServerError)
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "Missing start or end parameter", http.StatusBadRequest)
		return
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		http.Error(w, "Invalid start parameter", http.StatusBadRequest)
		return
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		http.Error(w, "Invalid end parameter", http.StatusBadRequest)
		return
	}

	if start > end {
		http.Error(w, "Start height cannot be greater than end height", http.StatusBadRequest)
		return
	}

	maxHeight := GlobalChain.Chain.Height()
	if end > maxHeight {
		end = maxHeight
	}

	blocks := []map[string]interface{}{}
	for i := start; i <= end; i++ {
		if i >= 0 && i < len(GlobalChain.Chain.Chains) {
			block := GlobalChain.Chain.Chains[i]
			blockData := map[string]interface{}{
				"height":    block.Height,
				"leader":    block.Leader,
				"scores":    block.Scores,
				"timestamp": block.Timestamp,
				"dataHash":  fmt.Sprintf("%x", block.DataHash),
				"prevHash":  fmt.Sprintf("%x", block.PrevBlockHash),
				"nonce":     block.Nonce,
				"version":   block.Version,
			}
			blocks = append(blocks, blockData)
		}
	}

	response := map[string]interface{}{
		"blocks":  blocks,
		"count":   len(blocks),
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 获取节点状态
func handleGetNodeStatus(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	if GlobalChain == nil {
		http.Error(w, "Chain not initialized", http.StatusInternalServerError)
		return
	}

	status := map[string]interface{}{
		"nodeId":      GlobalChain.Id,
		"address":     GlobalChain.Addr,
		"isLeader":    GlobalChain.IsLeader,
		"chainHeight": GlobalChain.Chain.Height(),
		"sequenceId":  GlobalChain.SequenceId,
		"poolSize":    GlobalChain.Pool.Size(),
		"poolCap":     GlobalChain.Pool.Cap,
		"success":     true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// 获取区块中的具体交易
func handleGetBlockTransactions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	if GlobalChain == nil {
		http.Error(w, "Chain not initialized", http.StatusInternalServerError)
		return
	}

	heightStr := r.URL.Query().Get("height")
	if heightStr == "" {
		http.Error(w, "Missing height parameter", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		http.Error(w, "Invalid height parameter", http.StatusBadRequest)
		return
	}

	if height < 0 || height >= len(GlobalChain.Chain.Chains) {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	block := GlobalChain.Chain.Chains[height]
	transactions := []map[string]interface{}{}

	// 这里需要从实际的区块数据中获取交易信息
	// 假设区块包含Eggplants作为交易
	for i, eggplant := range block.Eggplants {
		tx := map[string]interface{}{
			"index":         i,
			"eggplantId":    eggplant.EggplantId,
			"nodeId":        eggplant.NodeId,
			"hash":          fmt.Sprintf("%x", eggplant.Hash),
			"firstSeen":     eggplant.FirstSeen,
			"productHeight": eggplant.ProductHeight,
			"productHash":   fmt.Sprintf("%x", eggplant.ProductHash),
		}
		transactions = append(transactions, tx)
	}

	response := map[string]interface{}{
		"height":       height,
		"transactions": transactions,
		"count":        len(transactions),
		"success":      true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 设置CORS头
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
