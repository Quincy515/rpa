package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"yiran/logic"
)

type AudioController struct{}

var DefaultAudioController = &AudioController{}

// 管理 tts 服务路由
func (t AudioController) RegisterRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// 需要登录后才能访问的放在下面
		v1.Use(JWTAuthMiddleware() /*, BasicAuthMiddleware()*/)
		v1.POST("/audio", t.play)
	}
}

// @Summary 图片识别语音播报
// @Description 可以把上传的图片文件进行文字识别、文字处理、语音合成返回MP3文件
// @Tags audio
// @Accept json
// @Produce json
// @Param fileName query string true "文件名"
// @Param Authorization header string true "Bearer ***"
// @Success 200 {object} ResponseData {"code": 1000, "message": "success", "data": { "uri": ""} }
// @Router /api/v1/audio [post]
func (t AudioController) play(c *gin.Context) {
	fileName := c.Query("fileName")
	mp3, err := logic.DefaultOCR.OcrReq(fileName)
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "处理图片识别语音合成失败")
	}
	// 3. 返回语音mp3文件的地址
	ResponseSuccess(c, map[string]interface{}{
		"uri": fmt.Sprintf("%s:%d/%s",
			viper.GetString("app.host"),
			viper.GetInt("app.port"),
			mp3),
	})
}
