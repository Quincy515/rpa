package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"path"
	"yiran/logic"
)

type FileController struct{}

var DefaultFileController = &FileController{}

// 管理文件服务路由
func (f FileController) RegisterRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// 需要登录后才能访问的放在下面
		v1.Use(JWTAuthMiddleware() /*, BasicAuthMiddleware()*/)
		v1.POST("/upload", f.upload)
	}
}

var allowFile = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".mp4":  true,
}

// @Summary 上传文件
// @Description 可以上传图片、视频文件
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param Authorization header string true "Bearer ***"
// @Success 200 {object} ResponseData {"code": 1000, "message": "success", "data": { "uri": ""} }
// @Router /api/v1/upload [post]
func (f FileController) upload(c *gin.Context) {
	// 1. 因为上传文件一定要登录，获取当前登录的用户信息
	userSn, err := getCurrentUserSn(c)
	if err != nil {
		zap.L().Error("用户上传操作时未登录")
		ResponseError(c, CodeNotLogin)
		return
	}

	// 2. 从 form 表单中获取文件内容句柄
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	defer file.Close()

	// 3. 确保上传的是允许的格式
	ext := path.Ext(head.Filename)
	if _, ok := allowFile[ext]; !ok {
		ResponseErrorWithMsg(c, CodeServerBusy, "未被允许的文件格式类型")
		return
	}

	// 4. 把文件内容转化为 []byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "文件解析失败")
		return
	}

	// 5. 文件格式对应的上传大小限制
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif":
		//zap.L().Debug("文件大小：", zap.Any("大小", len(buf.Bytes())), zap.Any("限制大小:", viper.GetInt("upload.img_max_size")))
		if len(buf.Bytes()) > viper.GetInt("upload.img_max_size") {
			ResponseError(c, CodeFileSizeImg)
			return
		}
	case ".mp4":
		if len(buf.Bytes()) > viper.GetInt("upload.video_max_size") {
			ResponseError(c, CodeFileSizeVideo)
			return
		}
	}

	// 6. 保存文件信息到数据库
	fileMeta, err := logic.DefaultFile.Create(userSn, head.Filename, buf.Bytes())
	if err != nil {
		ResponseErrorWithMsg(c, CodeServerBusy, "保存文件失败")
		return
	}
	ResponseSuccess(c, map[string]interface{}{
		"uri": fmt.Sprintf("%s:%d/%s",
			viper.GetString("app.host"),
			viper.GetInt("app.port"),
			fileMeta.Location),
	})
}
