package main

import (
	"fmt"
	"os"
	"strconv"
)

type RcloneCopyOpts struct {
	BufferSize string
	DryRun     bool
}

func (c RcloneCopyOpts) AsFlagArray() []string {
	args := []string{"-v", "-P", "--no-update-modtime"}

	args = append(args, "--buffer-size", c.BufferSize)
	if c.DryRun {
		args = append(args, "--dry-run")
	}
	return args
}

type HddLargeFileCopyParam struct {
	Copytype TransferType
	Pwd      string
	SrcFile  string
	DestPath string
	DryRun   bool
	Filename string
}

type RcloneB2CopyOpts struct {
	copyOpts          RcloneCopyOpts
	ChunkSize         string
	UploadConcurrency int
}

func (c RcloneB2CopyOpts) AsFlagArray() []string {
	args := c.copyOpts.AsFlagArray()
	args = append(args, "--b2-chunk-size", c.ChunkSize)
	args = append(args, "--b2-upload-concurrency", strconv.Itoa(c.UploadConcurrency))
	args = append(args, "--b2-disable-checksum") // very slow for large files
	return args
}

func LargeFileCopy(params HddLargeFileCopyParam) error {

	st, err := os.Stat(params.SrcFile)
	if err != nil {
		return fmt.Errorf("could not stat: %w", err)
	}
	if st.IsDir() {
		return ErrWontCopyDir
	}

	var args []string

	switch params.Copytype {
	case Local:
		// rclone copy -v -P --buffer-size 0 --dry-run? ./src/file remote:/src/
		opts := RcloneCopyOpts{
			BufferSize: "0", // give it all you got
			DryRun:     params.DryRun,
		}
		args = opts.AsFlagArray()
	case Cloud:
		// rclone copy -v -P --b2-disable-checksum  --buffer-size 0 --b2-chunk-size 96Mi --b2-upload-concurrency 64 --dry-run? ./src/file remote:/src/
		opts := RcloneB2CopyOpts{
			copyOpts: RcloneCopyOpts{
				BufferSize: "0", // give it all you got
				DryRun:     params.DryRun,
			},
			ChunkSize:         "96Mi",
			UploadConcurrency: 64,
		}
		args = opts.AsFlagArray()
	default:
		return ErrInvalidTransferType
	}

	cmd := buildRcloneCommandCommand("copy", params.Pwd, params.SrcFile, params.DestPath, args)
	return cmd.Run()
}
