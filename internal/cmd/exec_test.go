package cmd

import (
	"log"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	stdout, stderr, err := ExecuteCmd("pwd", 10*time.Second)
	log.Printf("stdout: %s\r\n", stdout)
	log.Printf("stderr: %s\r\n", stderr)
	log.Printf("error: %v\r\n", err)
}
