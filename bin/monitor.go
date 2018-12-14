package main

import (
	"custom/test/monitor/plugins/http"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"path/filepath"
	"sync"
)

var dir string

func initLog() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	// 通过配置文件配置
	logger.SetLogger(fmt.Sprintf("%s/log.json", dir))
}

var wg sync.WaitGroup

func main() {
	//首先初始化日志配置
	initLog()
	fmt.Println("初始化日志成功")
	//启动http服务
	wg.Add(1)
	go func() {
		http.HttpListerner()
	}()
	wg.Add(1)
	go func() {
		fmt.Println("启动程序更新服务")
	}()
	wg.Wait()
	//
	//启动发送日志数据
	logger.Debug("123")
}
