package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yiran/dao/mysql"
	"yiran/global"
	"yiran/http/controller"
	"yiran/logger"
	"yiran/pkg/snowflake"
)

var (
	// 编译相关信息
	// 初始化为 unknown，如果编译时没有传入这些值，则为 unknown
	gitCommitLog = "unknown"
	buildTime    = "unknown"
	gitRelease   = "unknown"
)
func init() {
	// 加载配置信息
	global.Init()
	global.App.FillBuildInfo(gitCommitLog, buildTime, gitRelease)
	fmt.Println("load config success")
}

// @title Gin swagger
// @version 1.0
// @description Gin swagger
func main() {
	// 初始化日志库
	logger.Init()
	zap.L().Info("init logger success")
	defer zap.L().Sync() // 将缓冲区内的日志条目落盘

	// 初始化MySQL链接
	if err := dao.Init(); err != nil {
		zap.L().Error("init mysql failed",
			zap.Error(err),
			zap.String("mysql", "xxx"),
			zap.Int("port", 3306))
		return
	}
	defer dao.Close()
	zap.L().Info("init mysql success")

	// 初始化分布式ID生成器
	if err := snowflake.Init(uint16(viper.GetInt("app.machine_id"))); err != nil {
		zap.L().Error("init snowflake failed", zap.Error(err))
	}

	// 加载路由信息
	router := controller.SetupRouters()

	/**
	服务器重启时对于正在访问网站的用户来说，直接就报服务端异常。

	优雅关机就是指

	1. 停止接收新请求
	2. 等待正在访问网站的用户收到响应后再关机。

	`net/http` 通过`srv.Shutdown(ctx)`原生支持优雅关机
	*/
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: router,
	}
	// 开启一个 goroutine 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	// 相当于告诉程序我给你5秒钟的时间你把没完成的请求处理一下，之后我们就要关机啦
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting...")
}