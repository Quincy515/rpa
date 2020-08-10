package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	myJWT "yiran/pkg/jwt"
)

const (
	ContextUserIDKey = "UserSn"
)

// 基于 JWT 实现的的认证中间件
// 对于需要登录才能访问的 API 来说
// 该中间件需要从请求头中获取 JWT Token
// 如果没有 Token --> /login
// 如果 Token 过期 --> /login
// 从 JWT 中解析我们需要的 UserSn 字段 --> 根据 UserSn 我们就能从数据库中查询到当前请求的用户是谁

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 有三种方式 1. 放在请求头 2. 放在请求体 3. 放在 URI
		// 这里假设 Token 放在 Header 的 Authorization: Bearer token_string 中
		// 这里的具体实现方式要依据实际业务情况修改
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少 Auth Token")
			c.Abort()
			return
		}
		// 按空格分隔
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ResponseErrorWithMsg(c, CodeInvalidToken, "请求头中 Auth 格式不正确")
			c.Abort()
			return
		}
		// parts[1] 是获取到的 tokenString, 使用解析 JWT 函数来解析它
		mc, err := myJWT.ParseToken(parts[1])
		if err != nil {
			ResponseError(c, CodeInvalidToken)
			zap.L().Warn("invalid JWT token", zap.Error(err))
			c.Abort()
			return
		}
		// 将当前请求的 UserSn 信息保存到请求的上下文 c 中
		c.Set(ContextUserIDKey, mc.UserSn)
		c.Next() // 后续的处理函数可以用 c.Get("UserSn") 来获取当前请求的用户信息
		// 返回响应的时候可以做 Token/ Cookie 续期
	}
}

// 基于Cookie和Session认证的中间件
// 对于需要登陆才能访问的API来说
// 该中间件需要从请求中获取Cookie值
// 如果没有Cookie --> /login
// 拿到Cookie值取session数据中找对应的数据，找不到(1.session过期了2.无效的cookie值) --> /login
// session值 也可以通过 c.Set() 直接赋值到 上下文c 上
func BasicAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c) // c 代表了请求相关的所有内容，获取当前请求对应的session数据
		token := session.Get("token")
		if token == nil { // 请求对应的 session 中找不到则说明不是登录用户
			//c.Redirect(http.StatusFound, "/login")
			ResponseError(c, CodeInvalidToken)
			c.Abort() // 终止当前请求的处理函数调用链
			return    // 终止当前处理函数
		}
		// 使用解析 JWT 函数来解析它
		mc, err := myJWT.ParseToken(token.(string))
		if err != nil {
			ResponseError(c, CodeInvalidToken)
			zap.L().Warn("invalid JWT token", zap.Error(err))
			c.Abort()
			return
		}
		// 将当前请求的 UserSn 信息保存到请求的上下文 c 中
		c.Set(ContextUserIDKey, mc.UserSn)
		c.Next() // 后续的处理函数可以用 c.Get("UserSn") 来获取当前请求的用户信息
		// 返回响应的时候可以做 Token/ Cookie 续期
	}
}
