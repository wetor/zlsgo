//go:build !windows
// +build !windows

package zshell

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func RunNewProcess(file string, args []string) (pid int, err error) {
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	if tmp, _ := ioutil.TempDir("", ""); tmp != "" {
		tmp = filepath.Dir(tmp)
		if strings.HasPrefix(file, tmp) {
			return 0, errors.New("temporary program does not support startup")
		}
	}
	return syscall.ForkExec(file, args, execSpec)
}

func RunBash(ctx context.Context, command string) (code int, outStr, errStr string, err error) {
	return ExecCommand(ctx, []string{
		"bash",
		"-c",
		command,
	}, nil, nil, nil)
}
