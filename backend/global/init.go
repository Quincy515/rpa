package global

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"sync"
	"time"
)

var once = new(sync.Once)

var (
	config = flag.String("config", "config", "配置文件名称，默认 config")
)

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		if !flag.Parsed() {
			flag.Parsed()
		}

		// 随机种子
		rand.Seed(time.Now().UnixNano())

		// 配置文件名称
		viper.SetConfigName(*config)
		// 配置文件查找路径
		viper.AddConfigPath("/etc/yiran/")
		viper.AddConfigPath("$HOME/.yiran")
		viper.AddConfigPath(App.RootDir + "/config")
		// 读取配置文件
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}

		// 监控配置文件变化
		viper.WatchConfig()

		// 填充 global.App 需要的数据
		App.fillOtherField()
	})
}
