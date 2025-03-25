package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildRcloneCommandCommand(rcloneCmd string, pwd string, src, dest string, args []string) *exec.Cmd {
	fullcmd := append([]string{rcloneCmd}, args...)
	fullcmd = append(fullcmd, src, dest)
	cmd := exec.Command("rclone", fullcmd...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("RCLONE_CONFIG_PASS=%s", pwd),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func printRcloneCommand(rcloneCmd string, args []string, src, dest string) {
	fullcmd := append([]string{"rclone", rcloneCmd}, args...)
	fullcmd = append(fullcmd, src, dest)
	fmt.Println("paste command > ", strings.Join(fullcmd, " "))
}
