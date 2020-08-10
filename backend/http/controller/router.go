package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"        // session包 定义了一套session操作的接口 类似于 database/sql
	"github.com/gin-contrib/sessions/cookie" // session具体存储的介质
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strings"
	"time"
	_ "yiran/docs"
	"yiran/logger"
)

// SetupRouters 配置项目路由信息
func SetupRouters() *gin.Engine {
	gin.SetMode(viper.GetString("app.mode"))
	r := gin.New() //r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(false))
	r.Use(Cors())

	// swagger api docs 仅在开发模式下查看
	if mode := gin.Mode(); mode == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 静态文件服务
	r.StaticFS("/static/upload", http.Dir(viper.GetString("upload.save_path")))

	// 设置session midddleware
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(720 * time.Hour), // 24h*30天
		Path:   "/",
	})
	r.Use(sessions.Sessions("mysession", store))

	DefaultUserController.RegisterRouter(r)  // 用户服务路由组
	DefaultFileController.RegisterRouter(r)  // 文件服务路由组
	DefaultAudioController.RegisterRouter(r) // rpa 服务路由组
	return r
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
