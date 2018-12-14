package http

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/wonderivan/logger"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func HttpListerner() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("系统运维服务接口"))
	})
	http.HandleFunc("/ops", Execmd)
	logger.Info("启动运维服务 .端口1210..")
	logger.Info(http.ListenAndServe(":1210", nil))
}

var mutex sync.Mutex

//关闭应用
func Execmd(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()         //先把这个方法锁住。避免两个人同时操作记错日志
	defer mutex.Unlock() //执行完成之后解锁
	username := r.FormValue("username")
	filename := r.FormValue("filename")
	if len(username) > 0 {
		var cmd *exec.Cmd
		var command string
		switch runtime.GOOS {
		case "windows":
			command = fmt.Sprintf(`./%s.bat `, filename)
			cmd = exec.Command("cmd.exe", "/c", command)
			logger.Info(fmt.Sprintf("操作人：%s,操作任务bat：%s", username, command))
			break
		case "linux":
			command = fmt.Sprintf(`./%s.sh `, filename)
			cmd = exec.Command("/bin/bash", "-c", command)
			logger.Info(fmt.Sprintf("操作人：%s,操作任务sh：%s", username, command))
			break
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("执行命令:%s 失败，原因是:%s", command, err.Error())))
			logger.Error(fmt.Sprintf("执行命令:%s 失败，原因是:%s", command, err.Error()))
			return
		}
		cmd.Start()
		reader := bufio.NewReader(stdout)
		//实时循环读取输出流中的一行内容
		w.Write([]byte(fmt.Sprintf("执行命令:%s 开始，输出:\n", command)))
		var buffer bytes.Buffer
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				buffer.WriteString("执行命令错误：" + err2.Error())
				break
			}
			b, _ := simplifiedchinese.GB18030.NewDecoder().Bytes([]byte(line))
			line = string(b)
			if runtime.GOOS == "windows" && strings.Contains(line, "exit") {
				buffer.WriteString(line)
				break
			}
			if !strings.EqualFold(line, "\r\n") {
				buffer.WriteString(line)
			}
		}
		logger.Info(buffer.String())
		w.Write([]byte(buffer.String()))
		cmd.Wait()
		w.Write([]byte(fmt.Sprintf("\n执行命令:%s 完成", command)))
		logger.Info(fmt.Sprintf("执行命令:%s 完成", command))
	} else {
		w.Write([]byte("不允许闲杂人等乱操作"))
		logger.Info(fmt.Sprintf("闲杂人操作一次,ip地址：%s", r.RemoteAddr))
	}
	return
}
