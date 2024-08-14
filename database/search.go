package database

import (
	"agricultural_meta/network"
	"agricultural_meta/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (nd *NodeDb) SearchNodeById(id int) (*network.Node, error) {
	node := network.Node{}
	res := nd.DB.Where("id=?", id).First(&node)
	if res.Error == nil {
		utils.LogMsg([]string{"SearchNode"}, []string{fmt.Sprintf("search node [id:%d addr:%s] successfully", node.Id, node.Addr)})
	} else {
		logrus.Errorln(res.Error.Error())
	}
	return &node, res.Error
}
