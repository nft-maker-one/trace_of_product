package database

import (
	"agricultural_meta/utils"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

// 通过 id 查询联盟链节点的公钥和 ip 地址
func (nd *NodeDb) SearchNodeById(id int) (*ConsortiumNode, error) {
	node := ConsortiumNode{}
	res := nd.DB.Where("id=?", id).First(&node)
	if res.Error == nil {
		utils.LogMsg([]string{"SearchNode"}, []string{fmt.Sprintf("search node [id:%d addr:%s] successfully", node.Id, node.Addr)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return &node, res.Error
}

func (nd *NodeDb) GetNum() int {
	node := ConsortiumNode{}
	var count int64
	nd.DB.Model(&node).Count(&count)
	utils.LogMsg([]string{"GetNum"}, []string{"当前联盟链共有节点：" + strconv.Itoa(int(count)) + " 个"})
	return int(count)
}
