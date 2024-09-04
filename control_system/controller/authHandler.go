package controller

import (
	"control_system/database"
	"control_system/models"
	"control_system/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthModel struct {
	myDb      *database.MyDb
	SecretKey string
}

func (am *AuthModel) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.LogMsg([]string{"Login"}, []string{"decode user failed err = " + err.Error()})
		ctx.JSON(400, gin.H{
			"error": "Invalid request Body",
		})
		return
	}
	var rUser models.User
	res := am.myDb.Db.Where("user_name=?", user.UserName).First(&rUser)
	if res.Error != nil {
		utils.LogMsg([]string{"login"}, []string{"login failed err = " + res.Error.Error()})
		ctx.JSON(401, gin.H{
			"error": "user do not exist",
		})
		return
	}
	if user.Password != rUser.Password {
		ctx.JSON(401, gin.H{
			"error": "password not correct",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(am.SecretKey))
	fmt.Println("login string = ", tokenString)
	// verToken,err := jwt.Parse(tokenString,)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Internal Server error",
		})
	}
	ctx.SetCookie("jwt-token", tokenString, 3600, "/", "example.com", false, true)
	ctx.JSON(200, gin.H{
		"token": tokenString,
	})

}

func (am *AuthModel) VerifyMiddleWare(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"msg":    "请完成登陆",
		})
		ctx.Abort()
		return
	}
	fmt.Println("verify string = " + tokenString)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(am.SecretKey), nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			utils.LogMsg([]string{"VerifyMiddleWare"}, []string{"verify jwt failed err = " + err.Error()})
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"msg":    "请更新登陆token",
		})
		ctx.Abort()
		return
	}

	ctx.Next()

}

func (am *AuthModel) Menu(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "欢迎登陆农产追溯通客户端",
	})
}

func NewAuthModel(path string) *AuthModel {
	cfg, err := utils.NewConfig(path)
	if err != nil {
		utils.LogMsg([]string{"NewAuthModel"}, []string{"read config file failed err =" + err.Error()})
	}
	am := &AuthModel{}
	am.SecretKey = cfg.Key
	am.myDb = database.NewMyDb(cfg.Mysql.Dsn)
	return am
}
