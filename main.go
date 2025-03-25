package main

import (
	"fmt"
	"os"
)

func main() {
	rclonePwd, ok := os.LookupEnv("RCLONE_CONFIG_PASS")
	if !ok {
		panic("pass not set")
	}

	hddSyncParams := RcloneSyncParams{
		Pwd:          rclonePwd,
		SrcDir:       "../D/mystuff",
		Dest:         "/media/ozoniuss/Expansion/mystuff/",
		DryRun:       true,
		TransferType: Local,
	}
	fmt.Println(RcloneSync(hddSyncParams))
}
