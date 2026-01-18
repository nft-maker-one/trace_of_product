package controller

import (
	"bytes"
	"control_system/database"
	"control_system/models"
	"control_system/utils"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MessageType byte
type CommandType string

type ChainModel struct {
	myDb *database.MyDb
	reDB *database.ReDb
}

type Node struct {
	Id   int    `json:"id"`
	Addr string `json:"addr"`
}

type UploadRequest struct {
	EggplantId      int    `json:"eggplant_id"`
	ProductHeight   int    `json:"product_height"`
	ProductHash     string `json:"product_hash"`
	TransportHeight int    `json:"transport_height"`
	TransportHash   string `json:"transport_hash"`
	ProcessHeight   int    `json:"process_height"`
	ProcessHash     string `json:"process_hash"`
	StorageHeight   int    `json:"storage_height"`
	StorageHash     string `json:"storage_hash"`
	SellHeight      int    `json:"sell_height"`
	SellHash        string `json:"sell_hash"`
}

const (
	MessageTypeProduce   MessageType = 0x1
	MessageTypeTransport MessageType = 0x2
	MessageTypeProcess   MessageType = 0x3
	MessageTypeStorage   MessageType = 0x4
	MessageTypeSell      MessageType = 0x5
	MessageTypeBlock     MessageType = 0x6
	MessageTypeEggplant  MessageType = 0x7
	MessageTypeTest      MessageType = 0x8
)

const (
	cRequest CommandType = "request"
)

type RPC struct {
	ContentType CommandType `json:"content_type"`
	Payload     []byte      `json:"payload"`
}

func (c *ChainModel) GetNodes(ctx *gin.Context) {
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "获取节点错误：" + res.Error.Error(),
		})
		return
	}
	ctx.JSON(200, nodes)
}

func (c *ChainModel) SendMessage(ctx *gin.Context) {
	id := ctx.DefaultQuery("id", "")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"msg":    "id lost",
		})
		return
	}
	node_address := ctx.DefaultQuery("node", "")
	if node_address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"msg":    "node not selected",
		})
		return
	}
	if res, err := c.reDB.Db.Exists(id).Result(); res == 1 && err == nil {
		data, err := c.reDB.Db.Get(id).Result()
		if err != nil {
			utils.LogMsg([]string{"SendMessage"}, []string{"get data from redis failed"})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    data,
		})
		return
	}

	req := models.ClientRequest{}
	req.Header = "Search"
	req.Content = []byte(id)
	req.RespAddr = "http://localhost:8081/meta_data"
	reqByte, err := json.Marshal(req)
	if err != nil {
		utils.LogMsg([]string{"SendMessage"}, []string{"clientRequest encode failed err =" + err.Error()})
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"statut": "error",
			"msg":    "Server Internal error",
		})
	}
	tcpDial(reqByte, node_address)
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		if res, err := c.reDB.Db.Exists(id).Result(); res == 0 || err != nil {
			continue
		} else {
			data, err := c.reDB.Db.Get(id).Result()
			if err != nil {
				utils.LogMsg([]string{"SendMessage"}, []string{"get data from redis failed"})
			}
			ctx.JSON(http.StatusOK, gin.H{
				"msg": data,
			})
			return
		}
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": "no such eggplant",
	})

}

func (c *ChainModel) UpdateData(ctx *gin.Context) {
	ur := UploadRequest{}
	err := ctx.ShouldBindJSON(&ur)
	if err != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"failed to decode user metadata err=" + err.Error()})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"msg":    "failed to decode user metadata",
		})
		return
	}
	req := models.ClientRequest{}
	meta, err := ur.ToMetaData()
	if err != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"to metadata failed err = " + err.Error()})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"msg":    "update format error " + err.Error(),
		})
		return
	}
	egg := models.Eggplant{}
	egg.MetaData = meta
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(&egg)
	req.Content = buf.Bytes()
	req.Header = "Upload"
	req.RespAddr = "no need to response"

	reqByte, err := json.Marshal(&req)
	if err != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"failed to encode clientRequest"})
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"msg":    "Internal server error",
		})
		return
	}
	// 自动选择可用的区块链节点
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"获取区块链节点失败: " + res.Error.Error()})
		// 如果数据库中没有节点，尝试使用容器中的默认节点
		defaultNodes := []string{"blockchain_node1:8081", "blockchain_node2:8081", "blockchain_node3:8081", "blockchain_node4:8081"}
		success := false
		for _, nodeAddr := range defaultNodes {
			err = tcpDial(reqByte, nodeAddr)
			if err == nil {
				success = true
				break
			}
			utils.LogMsg([]string{"UpdateData"}, []string{"连接节点 " + nodeAddr + " 失败: " + err.Error()})
		}
		if !success {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"msg":    "所有区块链节点连接失败",
			})
			return
		}
	} else {
		// 使用数据库中的节点
		success := false
		for _, node := range nodes {
			err = tcpDial(reqByte, node.Addr)
			if err == nil {
				success = true
				utils.LogMsg([]string{"UpdateData"}, []string{"成功连接节点: " + node.Addr})
				break
			}
			utils.LogMsg([]string{"UpdateData"}, []string{"连接节点 " + node.Addr + " 失败: " + err.Error()})
		}
		if !success {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"msg":    "所有区块链节点连接失败",
			})
			return
		}
	}
	ctx.JSON(200, gin.H{
		"status": "ok",
		"msg":    "上传成功",
	})
}

func (c *ChainModel) HandleChainResponse(ctx *gin.Context) {
	metaData := models.MetaData{}
	if err := ctx.ShouldBindJSON(&metaData); err != nil {
		utils.LogMsg([]string{"HandleChainResponse"}, []string{"can not decode metadata err=" + err.Error()})
		return
	}
	if !metaData.Verify() {
		utils.LogMsg([]string{"HandleChainResponse"}, []string{"metadata verify failed"})
		return
	}
	key := strconv.Itoa(metaData.EggplantId)
	byteData, err := json.Marshal(metaData)
	if err != nil {
		utils.LogMsg([]string{"HandleChainResponse"}, []string{"metaData marshal failed err=" + err.Error()})
	}

	res := c.reDB.Db.Set(key, string(byteData), 60*60*time.Second)
	if res.Err() != nil {
		utils.LogMsg([]string{""}, []string{"set key failed err = " + res.Err().Error()})
	}
}

func tcpDial(data []byte, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to connect to %s", addr)})
		logrus.Errorln(err.Error())
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to write data to %s", addr)})
		logrus.Errorln(err.Error())
		return fmt.Errorf("failed to connct to consortium nodes")
	}
	rpc := RPC{}
	rpc.ContentType = cRequest
	rpc.Payload = data
	req, err := json.Marshal(rpc)
	if err != nil {
		utils.LogMsg([]string{"tcpDial"}, []string{"failed to encode rpc"})
		return fmt.Errorf("rpc format error")
	}

	_, err = conn.Write(req)
	if err != nil {
		utils.LogMsg([]string{"tcpDial"}, []string{"failed to write data err =" + err.Error()})
	}
	conn.Close()
	return nil
}

// 获取区块链节点列表
func (c *ChainModel) GetBlockchainNodes(ctx *gin.Context) {
	// 从数据库获取区块链节点列表
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "获取区块链节点失败：" + res.Error.Error(),
		})
		return
	}

	// 转换为适合前端的格式
	blockchainNodes := make([]gin.H, 0, len(nodes))
	for _, node := range nodes {
		blockchainNodes = append(blockchainNodes, gin.H{
			"id":          node.Id,
			"addr":        node.Addr,
			"pub_key":     node.PubKey,
			"create_time": node.CreateTime,
			"verify_time": node.VerifyTime,
			"http_addr":   "http://" + node.Addr[:len(node.Addr)-1] + "0", // 将8081改为8080用于HTTP调用
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   blockchainNodes,
	})
}

// 获取区块链最高区块高度
func (c *ChainModel) GetBlockchainHeight(ctx *gin.Context) {
	// 获取第一个可用的区块链节点
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil || len(nodes) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "未找到可用的区块链节点",
		})
		return
	}

	// 使用第一个节点获取高度信息
	nodeAddr := nodes[0].Addr
	// 将端口从8081改为8080用于HTTP调用
	httpAddr := "http://" + nodeAddr[:len(nodeAddr)-1] + "0/api/height"

	resp, err := http.Get(httpAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "连接区块链节点失败：" + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var heightResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&heightResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "解析区块链响应失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, heightResponse)
}

// 根据区块范围获取区块数据
func (c *ChainModel) GetBlocksByRange(ctx *gin.Context) {
	start := ctx.DefaultQuery("start", "0")
	end := ctx.DefaultQuery("end", "10")

	// 获取第一个可用的区块链节点
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil || len(nodes) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "未找到可用的区块链节点",
		})
		return
	}

	// 使用第一个节点获取区块数据
	nodeAddr := nodes[0].Addr
	// 将端口从8081改为8080用于HTTP调用
	httpAddr := fmt.Sprintf("http://%s0/api/blocks?start=%s&end=%s", nodeAddr[:len(nodeAddr)-1], start, end)

	resp, err := http.Get(httpAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "连接区块链节点失败：" + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var blocksResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&blocksResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "解析区块链响应失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, blocksResponse)
}

// 获取区块链节点状态
func (c *ChainModel) GetNodeStatus(ctx *gin.Context) {
	// 获取第一个可用的区块链节点
	var nodes []models.ConsortiumNode
	if res := c.myDb.Db.Find(&nodes); res.Error != nil || len(nodes) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "未找到可用的区块链节点",
		})
		return
	}

	// 使用第一个节点获取状态信息
	nodeAddr := nodes[0].Addr
	// 将端口从8081改为8080用于HTTP调用
	httpAddr := "http://" + nodeAddr[:len(nodeAddr)-1] + "0/api/node/status"

	resp, err := http.Get(httpAddr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "连接区块链节点失败：" + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var statusResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"msg":    "解析区块链响应失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, statusResponse)
}

func NewChainModel(path string) *ChainModel {
	cfg, err := utils.NewConfig(path)
	if err != nil {
		utils.LogMsg([]string{"NewChainModel"}, []string{"read config file failed err =" + err.Error()})
	}
	cm := ChainModel{}
	cm.myDb = database.NewMyDb(cfg.Postgresql.Dsn)
	cm.reDB = database.NewReDb(cfg.Redis.Ip, cfg.Redis.Port)
	return &cm

}

func (ur *UploadRequest) ToMetaData() (models.MetaData, error) {

	pHash, err := hex.DecodeString(ur.ProductHash)
	if err != nil {
		return models.MetaData{}, fmt.Errorf("decode ProductHash[%v] failed", ur.ProductHash)
	}
	tHash, err := hex.DecodeString(ur.TransportHash)
	if err != nil {
		return models.MetaData{}, fmt.Errorf("decode TransportHash[%v] failed", ur.TransportHash)
	}
	pcHash, err := hex.DecodeString(ur.ProcessHash)
	if err != nil {
		return models.MetaData{}, fmt.Errorf("decode ProcessHash[%v] failed", ur.ProcessHash)
	}
	sHash, err := hex.DecodeString(ur.StorageHash)
	if err != nil {
		return models.MetaData{}, fmt.Errorf("decode StorageHash[%v] failed", ur.StorageHash)
	}
	seHash, err := hex.DecodeString(ur.SellHash)
	if err != nil {
		return models.MetaData{}, fmt.Errorf("decode SellHash[%v] failed", seHash)
	}
	res := models.MetaData{}
	res.EggplantId = ur.EggplantId
	res.ProductHash, _ = models.BytesToHash(pHash)
	res.ProductHeight = ur.ProductHeight
	res.TransportHash, _ = models.BytesToHash(tHash)
	res.TransportHeight = ur.TransportHeight
	res.ProcessHash, _ = models.BytesToHash(pcHash)
	res.ProcessHeight = ur.ProcessHeight
	res.StorageHash, _ = models.BytesToHash(sHash)
	res.StorageHeight = ur.StorageHeight
	res.SellHash, _ = models.BytesToHash(seHash)
	res.SellHeight = ur.SellHeight
	return res, nil

}
