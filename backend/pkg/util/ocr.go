package util

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ocr(baseUrl, path string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("获取文件失败:", err)
	}

	sourceString := base64.StdEncoding.EncodeToString(fileBytes)

	token, err := FetchToken(viper.GetString("ocr.apiKey"), viper.GetString("ocr.secretKey"))
	if err != nil {
		return nil, err
	}

	urlStr := fmt.Sprintf("%s?access_token=%s", baseUrl, token)
	//todo options参数抽出来
	params := url.Values{
		"image": {sourceString},
	}
	res, err := http.PostForm(urlStr, params)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	return body, err
}

func GeneralBasic(path string) ([]byte, error) {
	//通用文字识别
	return ocr(viper.GetString("ocr.generalBasic"), path)
}
