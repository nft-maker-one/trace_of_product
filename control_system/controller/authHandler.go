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
	utils.LogMsg([]string{"Login"}, []string{rUser.Nickname})
	userInfo := gin.H{
		"id":            rUser.ID,
		"username":      rUser.UserName,
		"nickname":      rUser.Nickname,
		"email":         rUser.Email,
		"avatar_url":    rUser.AvatarUrl,
		"last_login_at": rUser.LastLoginAt,
		"created_at":    rUser.CreatedAt,
		"updated_at":    rUser.UpdatedAt,
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "login successful",
		"data": gin.H{
			"token":      tokenString,
			"user":       userInfo,
			"expires_in": 86400,
		},
	})

}

func (am *AuthModel) Register(ctx *gin.Context) {
	var registerData struct {
		Username   string `json:"username" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=8"`
		InviteCode string `json:"invite_code"`
		AgreeTerms bool   `json:"agree_terms" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&registerData); err != nil {
		utils.LogMsg([]string{"Register"}, []string{"decode register data failed err = " + err.Error()})
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	am.myDb.Db.Where("user_name = ?", registerData.Username).Find(&existingUser)

	// 如果找到了记录（ID不为零），说明用户名已存在
	if existingUser.ID != 0 {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "username already exists",
		})
		return
	}

	// 检查邮箱是否已存在
	var emailExists int64
	emailRes := am.myDb.Db.Model(&models.User{}).Where("email = ?", registerData.Email).Count(&emailExists)
	if emailRes.Error != nil {
		utils.LogMsg([]string{"Register"}, []string{"email check error: " + emailRes.Error.Error()})
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	if emailExists > 0 {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "email already registered",
		})
		return
	}

	// 创建新用户
	newUser := models.User{
		UserName:   registerData.Username,
		Nickname:   registerData.Username, // 初始化nickname为username
		Email:      registerData.Email,
		Password:   registerData.Password,
		InviteCode: registerData.InviteCode,
		AgreeTerms: registerData.AgreeTerms,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := am.myDb.Db.Create(&newUser).Error; err != nil {
		utils.LogMsg([]string{"Register"}, []string{"failed to create user err = " + err.Error()})
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	utils.LogMsg([]string{"Register"}, []string{"user " + registerData.Username + " registered successfully"})
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Registration successful",
		"user_id": newUser.ID,
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

	// 将claims存储到context中供后续处理器使用
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		ctx.Set("claims", claims)
	}

	ctx.Next()

}

func (am *AuthModel) Menu(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "欢迎登陆农产追溯通客户端",
	})
}

func (am *AuthModel) UpdateProfile(ctx *gin.Context) {
	// 从JWT中获取用户名
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户信息获取失败",
		})
		return
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无效的用户信息",
		})
		return
	}

	username, ok := jwtClaims["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无法获取用户名",
		})
		return
	}

	// 解析更新数据，支持nickname、profile_public、email_notifications
	var updateData struct {
		Email              string `json:"email" binding:"omitempty,email"`
		Password           string `json:"password" binding:"omitempty,min=8"`
		InviteCode         string `json:"invite_code"`
		Nickname           string `json:"nickname"`
		ProfilePublic      *bool  `json:"profile_public"`
		EmailNotifications *bool  `json:"email_notifications"`
	}

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		utils.LogMsg([]string{"UpdateProfile"}, []string{"decode update data failed err = " + err.Error()})
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "无效的请求数据",
		})
		return
	}

	// 查找当前用户
	var currentUser models.User
	res := am.myDb.Db.Where("user_name = ?", username).First(&currentUser)
	if res.Error != nil {
		utils.LogMsg([]string{"UpdateProfile"}, []string{"user not found: " + username})
		ctx.JSON(404, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}
	// 如果更新邮箱，检查邮箱是否已被其他用户使用
	if updateData.Email != "" && updateData.Email != currentUser.Email {
		var emailExists int64
		am.myDb.Db.Model(&models.User{}).Where("email = ? AND id != ?", updateData.Email, currentUser.ID).Count(&emailExists)
		if emailExists > 0 {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "邮箱已被其他用户注册",
			})
			return
		}
	}
	// 构建更新字段
	updates := make(map[string]interface{})
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}
	if updateData.Password != "" {
		updates["password"] = updateData.Password
	}
	if updateData.InviteCode != "" {
		updates["invite_code"] = updateData.InviteCode
	}
	if updateData.Nickname != "" {
		updates["nick_name"] = updateData.Nickname
	}
	if updateData.ProfilePublic != nil {
		updates["profile_public"] = *updateData.ProfilePublic
	}
	if updateData.EmailNotifications != nil {
		updates["email_notifications"] = *updateData.EmailNotifications
	}
	// 如果有更新字段，则执行更新
	if len(updates) > 0 {
		updates["updated_at"] = time.Now()

		if err := am.myDb.Db.Model(&currentUser).Updates(updates).Error; err != nil {
			utils.LogMsg([]string{"UpdateProfile"}, []string{"failed to update user err = " + err.Error()})
			ctx.JSON(500, gin.H{
				"success": false,
				"message": "更新用户信息失败",
			})
			return
		}
	}
	// 重新查询更新后的用户信息
	am.myDb.Db.Where("user_name = ?", username).First(&currentUser)

	utils.LogMsg([]string{"UpdateProfile"}, []string{"user " + username + " profile updated successfully"})
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "个人信息更新成功",
		"data": gin.H{
			"id":         currentUser.ID,
			"username":   currentUser.UserName,
			"email":      currentUser.Email,
			"avatar_url": currentUser.AvatarUrl,
			"updated_at": currentUser.UpdatedAt,
		},
	})
}

func (am *AuthModel) UpdateAvatar(ctx *gin.Context) {
	// 从JWT中获取用户名
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户信息获取失败",
		})
		return
	}

	jwtClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无效的用户信息",
		})
		return
	}

	username, ok := jwtClaims["username"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "无法获取用户名",
		})
		return
	}

	// 解析头像数据
	var avatarData struct {
		AvatarUrl string `json:"avatar_url" binding:"required,url"`
	}

	if err := ctx.ShouldBindJSON(&avatarData); err != nil {
		utils.LogMsg([]string{"UpdateAvatar"}, []string{"decode avatar data failed err = " + err.Error()})
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "无效的头像URL",
		})
		return
	}

	// 查找当前用户
	var currentUser models.User
	res := am.myDb.Db.Where("user_name = ?", username).First(&currentUser)
	if res.Error != nil {
		utils.LogMsg([]string{"UpdateAvatar"}, []string{"user not found: " + username})
		ctx.JSON(404, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 更新头像URL
	if err := am.myDb.Db.Model(&currentUser).Updates(map[string]interface{}{
		"avatar_url": avatarData.AvatarUrl,
		"updated_at": time.Now(),
	}).Error; err != nil {
		utils.LogMsg([]string{"UpdateAvatar"}, []string{"failed to update avatar err = " + err.Error()})
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "更新头像失败",
		})
		return
	}

	utils.LogMsg([]string{"UpdateAvatar"}, []string{"user " + username + " avatar updated successfully"})
	ctx.JSON(200, gin.H{
		"success": true,
		"message": "头像更新成功",
		"data": gin.H{
			"id":         currentUser.ID,
			"username":   currentUser.UserName,
			"avatar_url": avatarData.AvatarUrl,
			"updated_at": currentUser.UpdatedAt,
		},
	})
}

func NewAuthModel(path string) *AuthModel {
	cfg, err := utils.NewConfig(path)
	if err != nil {
		utils.LogMsg([]string{"NewAuthModel"}, []string{"read config file failed err =" + err.Error()})
	}
	am := &AuthModel{}
	am.SecretKey = cfg.Key
	am.myDb = database.NewMyDb(cfg.Postgresql.Dsn)
	return am
}
