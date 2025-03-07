package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"wsync/conf"
	"wsync/global"
)

func SyncFile(_filePath string, _remoteUser string, _remoteIP string, _remotePath string) error {
	var cmd *exec.Cmd

	switch conf.GlobalConfig.Role {
	case "sender":
		cmd = exec.Command("sudo", "-u", global.GlobalWsyncer.WorkUser, "rsync", "-avz", _filePath, _remoteUser+"@"+_remoteIP+":"+_remotePath)
		// cmd = exec.Command("rsync", "-avz", _filePath, _remoteUser+"@"+_remoteIP+":"+_remotePath)	
	case "accepter":

	case "puller":
		cmd = exec.Command("sudo", "-u", global.GlobalWsyncer.WorkUser, "rsync", "-avz", _remoteUser+"@"+_remoteIP+":"+_remotePath, _filePath)
		// cmd = exec.Command("rsync", "-avz", _remoteUser+"@"+_remoteIP+":"+_remotePath, _filePath)
	}

	cmd.Env = append(os.Environ(), fmt.Sprintf("TMPDIR=%s/tmp", global.GlobalWsyncer.WorkDir))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdoutpipe error: %v", err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderrpipe error: %v", err.Error())
	}

	global.GlobalWsyncer.Wg.Add(1)
	go func() {
		defer global.GlobalWsyncer.Wg.Done()
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				log.Printf("cannot read from stderr: %s", err.Error())
				return
			}
			select {
			case <-global.GlobalWsyncer.Cancelctx.Done():
				return
			default:
				fmt.Printf("\033[31m[stderr]\033[0m%s", line)
			}
		}
	}()

	global.GlobalWsyncer.Wg.Add(1)
	go func() {
		defer global.GlobalWsyncer.Wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadBytes('\n')
			if len(line) == 0 || errors.Is(err, io.EOF) {
				return
			}
			select {
			case <-global.GlobalWsyncer.Cancelctx.Done():
				return
			default:
				fmt.Printf("[stdout]%s", line)
			}
		}
	}()

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error syncing file %s: %s", _filePath, err.Error())
	}
	return nil
}
