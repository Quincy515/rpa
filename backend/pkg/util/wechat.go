package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"yiran/model"

	"github.com/spf13/viper"
)

func GetWeChatId(code string) (data *model.Code2Session, err error) {
	// 1. 发起 Get 请求
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		viper.GetString("wechat.appid"), viper.GetString("wechat.secret"), code)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 2. 读取请求返回的结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 3. 返回结果转换
	data = &model.Code2Session{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	return
}
