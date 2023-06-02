package main

import (
	"bufio"
	"encoding/json"
	"github.com/go-toast/toast"
	"io"
	"os"
	"strings"
)

type JsonLog struct {
	TraceId     string `json:"traceId"`
	SpanId      string `json:"spanId"`
	RequestType string `json:"request_type"`
	Timestamp   string `json:"@timestamp"`
	Msg         string `json:"msg"`
	Version     string `json:"@version"`
	Level       string `json:"level"`
	StackTrace  string `json:"stack_trace"`
	Thread      string `json:"thread"`
	Class       string `json:"class"`
	Method      string `json:"method"`
	Line        string `json:"line"`
	App         string `json:"app"`
	LogPos      string `json:"log_pos"`
	Data        string `json:"data"`
}

func main() {
	args := os.Args

	if len(args) == 1 {
		return
	}

	if len(args) > 1 {
		for i, arg := range args {
			if i >= 1 {
				convert(arg)
			}
		}
	}
}

func toastPush(msg string) {
	notification := toast.Notification{
		AppID:   "j2l",
		Message: msg,
	}
	notification.Push()
}

func convert(filepath string) {

	//打开文件
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		toastPush(err.Error())
		return
	}

	//创建新的输出文件
	fileName := file.Name()
	index := strings.LastIndex(file.Name(), ".")
	newFileName := fileName[:index] + ".j2l" + fileName[index:]
	newFile, err := os.Create(newFileName)
	defer newFile.Close()
	if err != nil {
		toastPush(err.Error())
		return
	}

	//按行读取
	var isFirstLine = true
	reader := bufio.NewReader(file)
	for {
		//按行读取字节
		bytes, rerr := reader.ReadBytes('\n')

		if len(bytes) == 0 {
			break
		}

		//去BOM
		if isFirstLine {
			isFirstLine = false
			if len(bytes) >= 3 && bytes[0] == 0xef && bytes[1] == 0xbb && bytes[2] == 0xbf {
				bytes = bytes[3:]
			}
		}

		//转结构体
		var jl JsonLog
		err := json.Unmarshal(bytes, &jl)
		if err != nil {
			if err.Error() == "unexpected end of JSON input" {
				continue
			} else {
				toastPush(err.Error())
				return
			}
		}

		logStr := jl.App + "|"
		logStr += jl.Timestamp + "|"
		logStr += jl.Thread + "|"
		logStr += jl.Level + "|"
		logStr += jl.Class + "|"
		logStr += jl.Method + "|"
		logStr += jl.Line + "|"
		logStr += jl.Msg + "\n"
		newFile.WriteString(logStr)
		if jl.StackTrace != "" {
			newFile.WriteString(jl.StackTrace)
		}

		if rerr == io.EOF {
			break
		}
	}

	toastPush("转换完成")
}
