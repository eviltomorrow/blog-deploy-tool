package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eviltomorrow/blog-deploy-tool/internal/cmd"
	"github.com/eviltomorrow/blog-deploy-tool/pkg/pb"
)

var (
	timeout = 20 * time.Second
)

type Cmd struct {
	pb.UnimplementedCmdServer
}

// Execute(context.Context, *Input) (*Output, error)

func (c *Cmd) Execute(ctx context.Context, req *pb.Input) (*pb.Output, error) {
	t, err := time.ParseDuration(req.Timeout)
	if err != nil {
		return nil, err
	}
	if t < 5*time.Second {
		t = timeout
	}

	stdout, stderr, err := cmd.ExecuteCmd(req.Text, t)
	if err != nil {
		return nil, fmt.Errorf("stdout: %v, stderr: %v, nest error: %v", stdout, stderr, err)
	}
	return &pb.Output{Data: stdout}, nil
}
