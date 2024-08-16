package database

import (
	"agricultural_meta/crypto"
	"agricultural_meta/utils"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NodeDb struct {
	DB    *gorm.DB
	Nodes []string
}

type ConsortiumNode struct {
	Id     int    `gorm:"id"`
	Addr   string `gorm:"string"`
	PubKey []byte `gorm:"pubkey"`
}

func initDataBase(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	utils.LogMsg([]string{"initDataBase"}, []string{"success"})
	return db
}

func InitNodeDb(dsn string) *NodeDb {
	db := NodeDb{}
	db.Nodes = make([]string, 0)
	db.DB = initDataBase(dsn)
	return &db
}

// 添加节点
func (nd *NodeDb) AddNode(id int, addr string, key crypto.PublicKey) error {
	node := ConsortiumNode{Id: id, Addr: addr, PubKey: key.ToSlice()}
	res := nd.DB.Create(&node)
	if res.Error == nil {
		utils.LogMsg([]string{"AddNode"}, []string{fmt.Sprintf("add node [id:%d addr:%s] to chain successfully", id, addr)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return res.Error
}

// 删除失去资格或者下线的联盟链节点
func (nd *NodeDb) DeleteNode(id int) error {
	res := nd.DB.Where("id=?", id).Delete(ConsortiumNode{})
	if res.Error == nil {
		utils.LogMsg([]string{"DeleteNode"}, []string{fmt.Sprintf("delete node [id:%d ] successfully", id)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return res.Error
}
