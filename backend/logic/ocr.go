package logic

import (
	"go.uber.org/zap"
	"strings"
	dao "yiran/dao/mysql"
	"yiran/pkg/util"
)

type ocrLogic struct{}

var DefaultOCR = ocrLogic{}

func (o ocrLogic) OcrReq(filename string) (string, error) {
	// 1. 根据文件名截取 fileHash
	fileHash := strings.Split(filename, ".")[0]
	// 2. 从文件表中查询是否有相同的hash文件记录
	fileMetaResp, err := dao.DefaultFile.GetFileMeta(fileHash)
	if err != nil {
		zap.L().Error("未找到该文件", zap.Error(err))
		return "", err
	}
	// 3. 根据文件的位置读取文件进行 OCR 识别
	result, err := util.GeneralBasic(fileMetaResp.Location)
	if err != nil {
		zap.L().Error("图片识别文字失败", zap.Error(err))
		return "", err
	}
	if string(result) == "" {
		zap.L().Error("未识别到文字", zap.Error(err))
		return "", err
	}
	zap.L().Debug("图片识别文字成功", zap.Any("", string(result)))
	// 4. 文字处理
	text := util.Sentences(string(result))
	// 5. 文字转语音合成
	resp, err, mp3Name := util.Text2audio(text)
	if resp != 200 || err != nil {
		zap.L().Error("处理文字转语音失败", zap.Error(err))
		return "", err
	}
	return mp3Name, nil
}
