package network

import "github.com/gin-gonic/gin"

func (s *NodeServer) ClientHandlerUp() {
	s.ClientHandler = gin.Default()

	// 数据上传接口

}
