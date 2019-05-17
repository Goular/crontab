package main

import (
	"flag"
	"fmt"
	"github.com/Goular/crontab/worker"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// 支持运行指令: worker -config ./master.json
	flag.StringVar(&confFile, "config", "./worker.json", "指定master.json文件")
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
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 启动日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	//启动执行器
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	// 启动调度器
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	// 初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	// 正常退出
	select {}

	return
ERR:
	fmt.Println(err)
}
