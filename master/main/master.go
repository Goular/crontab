package main

import (
	"flag"
	"fmt"
	"github.com/Goular/crontab/master"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// 支持运行指令: master -config ./master.json
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json文件")
	flag.Parse()
}

// 初始化协程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	err error
)

func main() {
	// 初始化命令行参数
	initArgs()

	// 初始化协程
	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	// 启动etcd任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 启动HTTP API服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 正常退出
	select {}

	return
ERR:
	fmt.Println(err)
}
