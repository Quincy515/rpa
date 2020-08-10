package global

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"yiran/pkg/util"
)

func init() {
	App.Version = "V1.0"
	App.LaunchTime = time.Now()
	App.Year = time.Now().Year()

	// 默认在项目根目录运行程序
	App.RootDir = "."

	// 用来处理单元测试或当前目录不在根目录，获取项目根目录
	if !viper.InConfig("http.port") {
		App.RootDir = inferRootDir()
	}
	App.TemplateDir = App.RootDir + "/template/"

	fileInfo, err := os.Stat(os.Args[0])
	if err != nil {
		panic(err)
	}

	App.Date = fileInfo.ModTime()

	App.Build.GoVersion = runtime.Version()
	App.Build.GinVersion = gin.Version
}

// inferRootDir 递归推导项目根目录
func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 先定义一个函数变量，然后在赋值函数时，里面使用了这个函数变量
	var infer func(d string) string
	infer = func(d string) string {
		if d == "/" {
			panic("请确保在项目根目录或子目录下运行程序，当前在: " + cwd)
		}
		if util.Exist(d + "/config") {
			return d
		}
		// 返回路径最后一个元素的目录,路径为空则返回.
		return infer(filepath.Dir(d))
	}
	return infer(cwd)
}

var App = &app{}

type app struct {
	Name    string
	Version string
	Date    time.Time

	// 项目根目录
	RootDir string
	// 模板根目录
	TemplateDir string

	// 启动时间
	LaunchTime time.Time
	Uptime     time.Duration

	Year int

	Domain string
	SEO    map[string]string

	Build struct {
		GitCommitLog string
		BuildTime    string
		GitRelease   string
		GoVersion    string
		GinVersion   string
	}

	Env string

	locker sync.Mutex
}

func (a *app) SetUptime() {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.Uptime = time.Now().Sub(a.LaunchTime)
}

func (a *app) FillBuildInfo(gitCommitLog, buildTime, gitRelease string) {
	a.Build.GitCommitLog = gitCommitLog
	a.Build.BuildTime = buildTime

	pos := strings.Index(gitRelease, "/")
	if pos >= -1 {
		a.Build.GitRelease = gitRelease[pos+1:]
	}

	fmt.Println(a)
}

func (a *app) fillOtherField() {
	a.Name = viper.GetString("app.name")
	a.Domain = viper.GetString("app.host")
	a.SEO = viper.GetStringMapString("seo")
	a.Env = viper.GetString("app.mode")
}

func (a *app) String() string {
	return "Build Info:" +
		"\nGit Commit Log: " + a.Build.GitCommitLog +
		"\nGit Release Info: " + a.Build.GitRelease +
		"\nBuild Time: " + a.Build.BuildTime +
		"\nGo Version: " + a.Build.GoVersion
}
