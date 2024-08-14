package database

import (
	"agricultural_meta/network"
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

func (nd *NodeDb) AddNode(id int, addr string) error {
	node := network.Node{Id: id, Addr: addr}
	res := nd.DB.Create(&node)
	if res.Error == nil {
		utils.LogMsg([]string{"AddNode"}, []string{fmt.Sprintf("add node [id:%d addr:%s] to chain successfully", id, addr)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return res.Error
}

func (nd *NodeDb) DeleteNode(id int) error {
	res := nd.DB.Where("id=?", id).Delete(network.Node{})
	if res.Error == nil {
		utils.LogMsg([]string{"DeleteNode"}, []string{fmt.Sprintf("delete node [id:%d ] successfully", id)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return res.Error
}
