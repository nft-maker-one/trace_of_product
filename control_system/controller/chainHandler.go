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
	NodeIp          string `json:"node_ip"`
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
	dialAddress := ur.NodeIp
	err = tcpDial(reqByte, dialAddress)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "连接节点失败，尝试向其他节点发起请求",
		})
		return
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

func NewChainModel(path string) *ChainModel {
	cfg, err := utils.NewConfig(path)
	if err != nil {
		utils.LogMsg([]string{"NewChainModel"}, []string{"read config file failed err =" + err.Error()})
	}
	cm := ChainModel{}
	cm.myDb = database.NewMyDb(cfg.Mysql.Dsn)
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
