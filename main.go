package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	packageName := ""
	if len(os.Args) > 1 {
		packageName = os.Args[1]
	} else {
		fmt.Println("command error")
		os.Exit(0)
	}
	fmt.Println(getPID(packageName))
	var cmd *exec.Cmd
	for {
		pid := getPID(packageName)
		if cmd == nil && pid != "" {
			cmd = logcat(pid)
		}
		if pid == "" && cmd != nil {
			err := cmd.Process.Kill()
			fmt.Println("LOGCAT: process kill ", pid)
			if err == nil {
				cmd = nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func getPID(packageName string) string {
	cmd := exec.Command("adb", "shell", "pidof", "-s", packageName)
	buf, _ := cmd.Output()
	pid := strings.Trim(string(buf), "\n ")
	return pid
}

func logcat(pid string) *exec.Cmd {
	params := []string{"logcat", "--pid=" + pid}
	if len(os.Args) > 2 {
		params = append(params, os.Args[2:]...)
	}
	cmd := exec.Command("adb", params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	go cmd.Start()
	return cmd
}
