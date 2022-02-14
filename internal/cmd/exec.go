package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// ExecuteCmd 执行 command
func ExecuteCmd(c string, timeout time.Duration) (string, string, error) {
	// 超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var ch = make(chan error)
	var cmd = exec.Command("/bin/sh", "-c", c)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", "", fmt.Errorf("start execute cmd failure, nest error: %v", err)
	}

	go func() {
		ch <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		cmd.Process.Signal(syscall.SIGINT)
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-ch
		close(ch)
		return stdout.String(), stderr.String(), fmt.Errorf("execute cmd timeout")

	case err := <-ch:
		close(ch)
		return stdout.String(), stderr.String(), err
	}
}
