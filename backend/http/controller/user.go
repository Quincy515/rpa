package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
	"yiran/logic"
	"yiran/model"
	myJWT "yiran/pkg/jwt"
)

type UserController struct{}

var DefaultUserController = &UserController{}

func (u UserController) RegisterRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{

		// 邮箱注册
		v1.POST("/email-register", u.emailRegister)
		// 邮箱登录
		v1.POST("/email-login", u.emailLogin)
		// 登出清除session缓存
		v1.GET("/logout", u.logout)
		// 用户通过微信小程序注册登录
		v1.POST("/mini-login", u.miniLogin)

		user := v1.Group("/u")
		{
			// 需要登录后才能访问的放在下面
			user.Use(JWTAuthMiddleware())
			// 微信小程序用户授权更新用户信息
			user.POST("/mini-user", u.miniUpdate)
		}
	}
}

// @Summary 邮箱注册
// @Description 通过邮箱创建新用户
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body model.AccountByEmail true "需要提交的注册信息"
// @Success 200 {object} ResponseData { "code": 1000, "message": "success", "data": null }
// @Router /api/v1/email-register [post]
func (u UserController) emailRegister(c *gin.Context) {
	// 1. 提取用户提交的注册信息
	var addAccount model.AccountByEmail
	if err := c.ShouldBindJSON(&addAccount); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 进行数据校验
	if err := addAccount.Validation(); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 3. 进行邮箱注册的逻辑处理
	user, err := logic.DefaultUserLogic.EmailRegister(addAccount)
	if err != nil {
		ResponseErrorWithMsg(c, CodeEmailRegister, err.Error())
		return
	}
	// 4. 注册成功构造token,自动登录
	tokenString, err := myJWT.GenToken(user.UserSn)
	if err != nil {
		zap.L().Error("myJWT.GenToken failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeGenToken, err.Error())
		return
	}
	// 但不记住登录状态
	session := sessions.Default(c)
	session.Set("token", tokenString)
	session.Save()
	// 5. 构造返回数据
	ResponseSuccess(c, map[string]interface{}{"userSn": user.UserSn, "access_token": tokenString})
}

// @Summary 邮箱登录
// @Description 通过邮箱密码进行登录
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body model.AccountByEmail true "需要提交的登录信息"
// @Success 200 {object} ResponseData { "code": 1000, "message": "success", "data": null }
// @Router /api/v1/email-login [post]
func (u UserController) emailLogin(c *gin.Context) {
	// 1. 提取用户提交的登录信息
	var account model.AccountByEmail
	if err := c.ShouldBindJSON(&account); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 进行数据校验
	if len(account.Email) == 0 || len(account.Password) == 0 {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 3. 通过邮箱登录的逻辑处理
	user, err := logic.DefaultUserLogic.EmailLogin(&account)
	if err != nil {
		ResponseErrorWithMsg(c, CodeEmailLogin, err.Error())
		return
	}
	// 4. 登录成功构造token
	rememberMe, err := getIDFromQuery(c, "remember_me")
	remember := cast.ToBool(rememberMe)
	var expireDuration time.Duration
	if remember {
		// 增加token的过期时间 session过期时间设置为30天，token过期时间默认为7天，这里增加到30天
		expireDuration = viper.GetDuration("jwt.expire")
	}
	tokenString, err := myJWT.GenToken(user.UserSn, int(expireDuration.Seconds()))
	if err != nil {
		zap.L().Error("myJWT.GenToken failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeGenToken, err.Error())
		return
	}
	session := sessions.Default(c)
	session.Set("token", tokenString)
	session.Save()

	// 5. 构造返回数据
	ResponseSuccess(c, map[string]interface{}{"userSn": user.UserSn, "access_token": tokenString})
}

// @Summary 微信小程序注册登录
// @Description 微信小程序的注册登录逻辑
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body model.MiniAccount true "微信小程序注册登录请求转结构体"
// @Success 200 {object} ResponseData { "code": 1000, "message": "success", "data": null }
// @Router /api/v1/mini-login [post]
func (u UserController) miniLogin(c *gin.Context) {
	// 1. 请求转结构体
	var addAccount model.MiniAccount
	if err := c.ShouldBindJSON(&addAccount); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 判断参数中是否有 code
	if len(addAccount.Code) < 1 {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 3. 处理微信小程序注册登录逻辑
	account, err := logic.DefaultUserLogic.MiniLogin(&addAccount)
	if err != nil {
		zap.L().Error("logic.DefaultUserLogic.MiniLogin failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeWXLogin, err.Error())
		return
	}
	// 4. 生成 JWT Token
	tokenString, err := myJWT.GenToken(account.UserSn)
	if err != nil {
		zap.L().Error("myJWT.GenToken failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeGenToken, err.Error())
		return
	}
	// 但不记住登录状态
	session := sessions.Default(c)
	session.Set("token", tokenString)
	session.Save()
	// 5. 构造返回数据
	ResponseSuccess(c, map[string]interface{}{"userSn": account.UserSn, "access_token": tokenString})
}

// @Summary 登出
// @Description 清除session缓存
// @Tags accounts
// @Success 200
// @Router /api/v1/mini-logout [get]
func (u UserController) logout(c *gin.Context) {
	// 清除用户登录状态的数据
	session := sessions.Default(c)
	session.Delete("token")
	session.Save()
}

// @Summary 微信小程序用户授权
// @Description 更新用户信息
// @Tags accounts
// @Param account body model.MiniAccount true "微信小程序注册登录请求转结构体"
// @Success 200 {object} ResponseData { "code": 1000, "message": "success", "data": model.Account }
// @Router /api/v1/u/mini-user [post]
func (u UserController) miniUpdate(c *gin.Context) {
	// 1. 请求转结构体
	var account model.MiniAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 2. 获取当前登录的用户信息
	userSn, err := getCurrentUserSn(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	// 3. 处理更新用户信息的逻辑
	userInfo, err := logic.DefaultUserLogic.WxUpdateUser(account, userSn)
	if err != nil {
		ResponseErrorWithMsg(c, CodeUpdateUser, err.Error())
		return
	}
	// 4. 构造返回数据
	ResponseSuccess(c, map[string]interface{}{"userInfo": userInfo})
}
