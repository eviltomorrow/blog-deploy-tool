package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/eviltomorrow/blog-deploy-tool/pkg/znet"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/ztime"
)

var (
	Pid         = os.Getpid()
	Pwd         string
	LaunchTime  = time.Now()
	HostName    string
	OS          = runtime.GOOS
	Arch        = runtime.GOARCH
	RunningTime = func() string {
		return ztime.FormatDuration(time.Since(LaunchTime))
	}
	IP string
)

func init() {
	path, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("get execute path failure, nest error: %v", err))
	}
	path = strings.ReplaceAll(path, "bin/blog-deploy-server", "")
	path = strings.ReplaceAll(path, "bin/blog-deploy", "")

	Pwd, err = filepath.Abs(path)
	if err != nil {
		panic(fmt.Errorf("get current folder failure, nest error: %v", err))
	}
	HostName, err = os.Hostname()
	if err != nil {
		panic(fmt.Errorf("get host name failure, nest error: %v", err))
	}
	IP, err = znet.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("get local ip failure, nest error: %v", err))
	}
}
