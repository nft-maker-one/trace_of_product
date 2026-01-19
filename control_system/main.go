package main

import (
	"control_system/controller"
	"control_system/models"
	"encoding/gob"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ChainHandler controller.ChainModel
var AuthHandler controller.AuthModel

func main() {
	engine := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	engine.Use(cors.New(config))

	// 设置jwt以及登录校验服务
	// 请求参数
	// {
	// 	"user_name": "string",
	// 	"password": "string"
	// }
	// 成功
	// {
	// 	"token": "string"
	// }
	// 失败
	// {
	// 	"error": "string"
	// }
	engine.POST("/login", AuthHandler.Login)

	// 注册新用户
	// 请求参数
	// {
	// 	"username": "string",
	// 	"email": "string",
	// 	"password": "string",
	// 	"invite_code": "string",
	// 	"agree_terms": true
	// }
	// 成功
	// {
	// 	"success": true,
	// 	"message": "Registration successful",
	// 	"user_id": 0
	// }
	// 失败
	// {
	// 	"success": false,
	// 	"message": "string"
	// }
	engine.POST("/register", AuthHandler.Register)

	// 更新用户资料
	// 请求参数
	// {
	// 	"nickname": "string",
	// 	"email": "string",
	// 	"profile_public": true,
	// 	"email_notifications": true
	// }
	engine.POST("/user/update/profile", AuthHandler.VerifyMiddleWare, AuthHandler.UpdateProfile)

	// 更新用户头像
	// 请求参数
	// {
	// 	"avatar_url": "string"
	// }
	engine.POST("/user/update/avatar", AuthHandler.VerifyMiddleWare, AuthHandler.UpdateAvatar)

	// 处理节点发送的农产品信息
	engine.POST("/meta_data", ChainHandler.HandleChainResponse)

	engine.POST("/upload", AuthHandler.VerifyMiddleWare, ChainHandler.UpdateData)

	// 获取节点列表
	// 请求参数
	// 无，但需要 JWT Token（Authorization Header）
	// 成功
	// [
	// {
	// 	"id": 0,                    // 节点ID
	// 	"addr": "string",            // 节点地址（IP:Port）
	// 	"pub_key": "[]byte",         // 节点公钥（字节数组）
	// 	"create_time": 0,            // 创建时间戳
	// 	"verify_time": 0             // 验证时间
	//   },
	//   ...
	// ]
	// 失败
	// {
	// 	"status": "error",
	// 	"msg": "获取节点错误：..."
	//   }
	engine.GET("/nodes", AuthHandler.VerifyMiddleWare, ChainHandler.GetNodes)

	// 测试节点连接状态
	// 请求参数
	// 需要 JWT Token（Authorization Header）
	// ?addr=127.0.0.1:8080
	// 成功
	// {
	// 	"status": "ok",
	// 	"msg": "节点连接成功",
	// 	"response_time": 123,
	// 	"online": true
	// }
	// 失败
	// {
	// 	"status": "error",
	// 	"msg": "节点连接失败: ...",
	// 	"response_time": -1,
	// 	"online": false
	// }
	engine.GET("/node/ping", AuthHandler.VerifyMiddleWare, ChainHandler.PingNode)

	// 根据农产品ID查询区块链中的农产品追溯信息
	// 请求参数
	// 需要 JWT Token（Authorization Header）
	// &id=0&node=127.0.0.1:8080
	// 成功
	// {
	// 	"status": "ok",
	// 	"msg": "..."
	//   }
	// 失败
	// {
	// 	"status": "error",
	// 	"msg": "查询错误：..."
	//   }
	engine.GET("/message", AuthHandler.VerifyMiddleWare, ChainHandler.SendMessage)
	// 菜单，欢迎信息
	// 主要用于前端验证用户是否已登录
	// 请求参数
	// 无，但需要 JWT Token（Authorization Header）
	// 如果 Token 无效，中间件会拦截并返回 401 错误
	// {
	// 	"status": "ok",
	// 	"msg": "欢迎登陆农产追溯通客户端"
	//   }
	engine.GET("/menu", AuthHandler.VerifyMiddleWare, AuthHandler.Menu)

	// 新增用户相关接口
	engine.GET("/profile", AuthHandler.VerifyMiddleWare, AuthHandler.GetProfile)
	engine.POST("/profile", AuthHandler.VerifyMiddleWare, AuthHandler.UpdateProfile)
	engine.POST("/upload/avatar", AuthHandler.VerifyMiddleWare, AuthHandler.UploadAvatar)

	// 区块链节点相关接口
	engine.GET("/blockchain/nodes", AuthHandler.VerifyMiddleWare, ChainHandler.GetBlockchainNodes)
	engine.GET("/blockchain/height", AuthHandler.VerifyMiddleWare, ChainHandler.GetBlockchainHeight)
	engine.GET("/blockchain/blocks", AuthHandler.VerifyMiddleWare, ChainHandler.GetBlocksByRange)
	engine.GET("/blockchain/status", AuthHandler.VerifyMiddleWare, ChainHandler.GetNodeStatus)

	// 添加静态文件服务用于访问上传的头像文件
	engine.Static("/public", "./public")

	engine.Run(":8080")
}

func init() {
	ChainHandler = *controller.NewChainModel("./config.yaml")
	AuthHandler = *controller.NewAuthModel("./config.yaml")
	gob.Register(models.Hash{})
}
