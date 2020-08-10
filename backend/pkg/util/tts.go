package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	formats = map[int]string{3: "mp3", 4: "pcm-16k", 5: "pcm-8k", 6: "wav"}
	format  = formats[viper.GetInt("tts.AUE")]
	token   string
	//ctx     *oto.Context // 如果已有音频正在播放，则关闭后开始播放新音频
)

// Text2audio 文字转音频
func Text2audio(text string) (code int, err error, filename string) {
	token, err = FetchToken(viper.GetString("tts.apiKey"), viper.GetString("tts.secretKey"))
	if err != nil {
		return 500, err, ""
	} else if token == "" {
		return 502, nil, ""
	}

	fmt.Println("text:", text)
	//urlT2a := "http://tsn.baidu.com/text2audio?tex=" + text + "&lan=zh&cuid=" + viper.GetInt("tts.CUID") + "&ctp=1&tok=" + token
	urlT2a := fmt.Sprintf("%s?tex=%s&lan=zh&cuid=%d&ctp=1&tok=%s",
		viper.GetString("tts.TTS_URL"), text,
		viper.GetInt("tts.CUID"), token)

	fmt.Println("urlT2a: ", urlT2a)

	resp, err := http.Get(urlT2a)
	if err == nil {
		contentType := resp.Header.Get("Content-Type")
		if contentType == "application/json" {
			// 返回json失败处理
			type ErrResp struct {
				ErrMsg string `json:"err_msg"`
				ErrNo  int    `json:"err_no"`
			}
			var respContent ErrResp

			body, err := ioutil.ReadAll(resp.Body)
			if err == nil && resp.StatusCode == 200 {
				_ = json.Unmarshal(body, &respContent)
			}
			_ = resp.Body.Close()
		} else if contentType == "audio/mp3" {
			// 返回audio mp3
			var audio []byte
			audio, err := ioutil.ReadAll(resp.Body)

			if err == nil {
				uploadDir := viper.GetString("upload.save_path") + time.Now().Format("2006/01/02/")
				err = os.MkdirAll(uploadDir, 0744) // os.ModePerm
				if err != nil {
					zap.L().Error("根据当前日期来创建文件夹失败", zap.Error(err))
					return 500, err, ""
				}
				// 1. 计算sha1
				fileHash := Sha1(audio)
				filename = uploadDir + fileHash + ".mp3"

				// write to file
				err = ioutil.WriteFile(filename, audio, 0666)
				if err != nil {
					fmt.Println("存储音频文件失败!", err)
					return 500, err, ""
				}
				fmt.Println("存储音频文件成功！")
			}
			_ = resp.Body.Close()
		} else {
			err = errors.New("Unknown Content-Type ! =>" + contentType + " | url: " + urlT2a)
		}
	}
	return 200, err, filename
}
