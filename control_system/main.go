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
	engine.POST("/login", AuthHandler.Login)
	engine.POST("/meta_data", ChainHandler.HandleChainResponse)
	engine.POST("/upload", AuthHandler.VerifyMiddleWare, ChainHandler.UpdateData)
	engine.GET("/nodes", AuthHandler.VerifyMiddleWare, ChainHandler.GetNodes)
	engine.GET("/message", AuthHandler.VerifyMiddleWare, ChainHandler.SendMessage)
	engine.GET("/menu", AuthHandler.VerifyMiddleWare, AuthHandler.Menu)
	engine.Run(":8081")
}

func init() {
	ChainHandler = *controller.NewChainModel("./config.yaml")
	AuthHandler = *controller.NewAuthModel("./config.yaml")
	gob.Register(models.Hash{})
}
