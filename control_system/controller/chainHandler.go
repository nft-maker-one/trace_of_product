package controller

import (
	"bytes"
	"control_system/database"
	"control_system/models"
	"control_system/utils"
	"encoding/gob"
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
	c.myDb.Db.Find(&nodes)
	resp := make(map[string]string)
	for _, v := range nodes {
		resp[strconv.Itoa(v.Id)] = v.Addr
	}
	ctx.JSON(200, resp)
}

func (c *ChainModel) SendMessage(ctx *gin.Context) {
	id := ctx.DefaultQuery("id", "0")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id is invalid",
		})
	}
	req := models.ClientRequest{}
	req.Header = "Search"
	req.Content = []byte(id)
	req.RespAddr = "http://localhost:8081/meta_data"
	reqByte, err := json.Marshal(req)
	if err != nil {
		utils.LogMsg([]string{"SendMessage"}, []string{"clientRequest encode failed err =" + err.Error()})
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Server Internal error",
		})
	}
	tcpDial(reqByte, "127.0.0.1:8083")
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
	reqMeta := models.ReqMetaData{}
	err := ctx.ShouldBindJSON(&reqMeta)
	fmt.Println(reqMeta)
	if err != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"failed to decode user metadata err=" + err.Error()})
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to decode user metadata",
		})
		return
	}
	req := models.ClientRequest{}
	meta := reqMeta.ToMetaData()
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
			"error": "Internal server error",
		})
		return
	}
	if err != nil {
		utils.LogMsg([]string{"UpdateData"}, []string{"failed to encode rpc"})
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Internal server error",
		})
		return
	}
	tcpDial(reqByte, "127.0.0.1:8083")
	ctx.JSON(200, gin.H{
		"msg": "upload completed",
	})
}

func (c *ChainModel) HandleChainResponse(ctx *gin.Context) {
	metaData := models.MetaData{}
	if err := ctx.ShouldBindJSON(&metaData); err != nil {
		utils.LogMsg([]string{"HandleChainResponse"}, []string{"can not decode metadata err=" + err.Error()})
		return
	}
	fmt.Println(metaData)
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

func tcpDial(data []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to connect to %s", addr)})
		logrus.Errorln(err.Error())
		utils.LogError([]string{"tcpDial"}, []string{fmt.Sprintf("failed to write data to %s", addr)})
		logrus.Errorln(err.Error())
	}
	rpc := RPC{}
	rpc.ContentType = cRequest
	rpc.Payload = data
	req, err := json.Marshal(rpc)
	if err != nil {
		utils.LogMsg([]string{"tcpDial"}, []string{"failed to encode rpc"})
	}

	conn.Write(req)
	conn.Close()
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
