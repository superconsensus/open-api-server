package conf

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

var (
	App    *AppConf
	Server *ServerConf
	Code   *CodeConf
	Log    *LogConf
	Req    *ReqConf
	Cache  *CacheConf
)

type AppConf struct {
	AppMode  string `ini:"app_mode"`
	RootPath string `ini:"root_path"`
}

type ServerConf struct {
	Protocol string `ini:"protocol"`
	HttpPort string `ini:"http_port"`
}

type CodeConf struct {
	CodePath    string   `ini:"code_path"`
	WasmPath    string   `ini:"wasm_path"`
	FileMaxSize int64    `ini:"file_max_size"`
	FileExts    []string `ini:"file_exts"`
}

type LogConf struct {
	FilePath    string `ini:"file_path"`
	FileName    string `ini:"file_name"`
	RouterFile  string `ini:"router_file"`
	RunTimeFile string `ini:"runtime_file"`
}

type ReqConf struct {
	Language int   `json:"language"`
	Strength uint8 `json:"strength"`
}

type CacheConf struct {
	Size int `ini:"size"`
}

//todo 有空了将这个改成yaml来配置 ini太麻烦了
//初始化配置
func init() {

	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		cfg = ini.Empty() //读不了文件就用默认值，避免空指针
	}

	App = &AppConf{
		AppMode:  "debug", //默认值
		RootPath: "runtime/",
	}
	cfg.Section("app").MapTo(App)

	Server = &ServerConf{
		Protocol: "http",
		HttpPort: "8080",
	}
	cfg.Section("server").MapTo(Server)

	Log = &LogConf{
		FilePath:    "logs/",
		FileName:    "app.log",
		RouterFile:  "router.log",
		RunTimeFile: "runtime.log",
	}
	cfg.Section("log").MapTo(Log)

	Code = &CodeConf{
		CodePath:    "contract_code/",
		WasmPath:    "contract_wasm/",
		FileMaxSize: 2,
		FileExts:    []string{"cc", "go"},
	}
	cfg.Section("code").MapTo(Code)

	Req = &ReqConf{
		Language: 2,
		Strength: 1,
	}
	cfg.Section("req").MapTo(Req)

	Cache = &CacheConf{
		Size: 10,
	}
	cfg.Section("cache").MapTo(Cache)

	//生产环境
	if App.AppMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
