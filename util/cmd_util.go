package util

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"github.com/go-cmd/run"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func ExecuteShell(filePath string) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "-c", filePath)
	case "linux":
		cmd = exec.Command("/bin/bash", "-c", filePath)
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	} else {
		return string(output)
	}
}

func ExecuteCommands(c ...*exec.Cmd) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var commands []*exec.Cmd
	for _, v := range c {
		commands = append(commands, v)
	}
	l := len(commands)
	// 前一个命令的输出，是后一个命令的输入
	for i := 1; i < l; i++ {
		commands[i].Stdin, _ = commands[i-1].StdoutPipe()
	}
	commands[l-1].Stdout = os.Stdout
	// 启动
	for i := 1; i < l; i++ {
		err := commands[i].Start()
		if err != nil {
			panic(err)
		}
	}
	// 第一个命令启动
	commands[0].Run()
	// 等候其他命令执行结果
	for i := 1; i < l; i++ {
		err := commands[i].Wait()
		if err != nil {
			panic(err)
		}
	}
	// 获取最后一个命令输出结果
	output, err := commands[l-1].Output()
	if err != nil {
		return ""
	} else {
		return string(output)
	}
}

// @see https://github.com/go-cmd/run/blob/master/sync_test.go
func ExecuteCommand(c cmd.Cmd) string {
	commands := []cmd.Cmd{
		c,
	}
	r := run.NewRunSync(true)
	var gotStatus []cmd.Status
	var gotErr error
	doneChan := make(chan struct{})
	go func() {
		gotErr = r.Run(commands)
		gotStatus, _ = r.Status()
		close(doneChan)
	}()

	var timeOut = 3 * time.Second
	time.Sleep(timeOut)

	// Check that Run returns ErrRunning on 2nd+ call
	//err := r.Run(commands)
	//if err != run.ErrRunning {
	//	return ""
	//}

	// Stop the first cmd which will return 0
	if err := r.Stop(); err != nil {
		return ""
	}

	// Run should return instantly after Stop
	select {
	case <-doneChan:
	case <-time.After(timeOut):
		return ""
	}

	if len(gotStatus) != len(commands) {
		return ""
	} else {
		return gotStatus[0].Stdout[0]
	}

}
