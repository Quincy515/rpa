package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// TransferData 转换数据的结构体
type TransferData struct {
	LogId          int64       `json:"log_id"`
	WordsResultNum int64       `json:"words_result_num"`
	WordsResult    []wordsList `json:"words_result"`
}

type wordsList struct {
	Words string `json:"words"`
}

// Sentences 图片转文字之后进一步结构化处理
func Sentences(text string) string {
	fmt.Println("整理前：", text)
	pubData := TransferData{}
	err := json.Unmarshal([]byte(text), &pubData)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	//var buffer string
	var buffer bytes.Buffer
	for _, p := range pubData.WordsResult {
		//buffer += p.Words
		buffer.WriteString(p.Words)
	}

	sentences := buffer.String()
	// 去除特殊字符空格
	sentences = strings.Join(strings.Fields(sentences), "")
	fmt.Println("整理后: ", sentences)
	return sentences
}
