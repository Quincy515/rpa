package controller

import (
	"fmt"
	"github.com/spf13/viper"
)

type MyCode int64

const (
	CodeSuccess         MyCode = 1000
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008

	CodeFileSizeImg   MyCode = 1009
	CodeFileSizeVideo MyCode = 1010

	CodeWXLogin       MyCode = 1011
	CodeGenToken      MyCode = 1012
	CodeEmailRegister MyCode = 1013
	CodeEmailLogin    MyCode = 1014
	CodeUpdateUser    MyCode = 1015
)

var msgFlags = map[MyCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",

	CodeFileSizeImg: fmt.Sprintf("上传图片大小超过限制%d",
		viper.GetInt("upload.img_max_size")),
	CodeFileSizeVideo: fmt.Sprintf("上传视频大小超过限制%d",
		viper.GetInt("upload.video_max_size")),

	CodeWXLogin:       "微信小程序注册登录失败",
	CodeGenToken:      "生成token失败",
	CodeEmailRegister: "邮箱注册失败",
	CodeEmailLogin:    "邮箱登录失败",
	CodeUpdateUser:    "微信小程序用户授权失败",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
