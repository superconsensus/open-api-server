package main

import (
	_ "github.com/superconsensus/open-api-server/conf"
	"github.com/superconsensus/open-api-server/routers"
)

func main() {
	// 从配置文件读取配置
	//conf.Init()

	// 装载路由
	r := routers.NewRouter()
	r.Run(":50400")
}
