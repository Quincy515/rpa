package logic

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
	"yiran/dao/mysql"
	"yiran/model"
	"yiran/pkg/util"
)

type fileLogic struct{}

var DefaultFile = fileLogic{}

// 上传文件的处理的逻辑
func (f fileLogic) Create(userSn uint64, filename string, buf []byte) (fileMeta *model.FileMeta, err error) {
	fileHash := util.Sha1(buf)
	// 1. 构建文件元信息
	fileMeta = &model.FileMeta{
		UserSn:   userSn,          // 文件属主
		FileName: filename,        // 文件名
		FileSha1: fileHash,        //　计算文件sha1
		FileSize: int64(len(buf)), // 文件大小
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	// 2. 从文件表中查询是否有相同的hash文件记录
	fileMetaResp, err := dao.DefaultFile.GetFileMeta(fileHash)
	if fileMetaResp != nil && err == nil {
		// 3. 如果查询到记录则实现秒传
		// 秒传实现：将上传过的文件信息写入用户文件表
		fileMeta.Location = fileMetaResp.Location
		err = dao.DefaultUser.SaveFile(fileMeta)
		if err != nil {
			zap.L().Error("秒传失败", zap.Error(err))
			return
		}
		zap.L().Debug("秒传成功")
		return
	}
	// 4. 查不到记录则将文件写入临时存储位置
	uploadDir := viper.GetString("upload.save_path") + time.Now().Format("2006/01/02/")
	err = os.MkdirAll(uploadDir, 0744) // os.ModePerm
	if err != nil {
		zap.L().Error("根据当前日期来创建文件夹失败", zap.Error(err))
		return
	}
	ext := path.Ext(filename)
	fileMeta.Location = uploadDir + fileMeta.FileSha1 + ext // 存储地址
	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		zap.L().Error("create file failed", zap.Error(err))
		return
	}
	defer newFile.Close()

	nByte, err := newFile.Write(buf)
	if int64(nByte) != fileMeta.FileSize || err != nil {
		zap.L().Error("save data into file failed", zap.Any("data", nByte), zap.Error(err))
		return
	}
	// 5. TODO: 同步或异步将文件转移到Ceph/OSS
	// 6. 更新文件表记录
	err = dao.DefaultFile.Create(fileMeta)
	if err != nil {
		zap.L().Error("mysql file failed", zap.Error(err))
		return
	}
	return
}
