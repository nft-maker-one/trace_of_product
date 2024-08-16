package database

import (
	"agricultural_meta/utils"
	"fmt"

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
