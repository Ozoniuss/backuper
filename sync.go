package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RcloneSyncOpts struct {
	Checkers       int
	Transfers      int
	BufferSize     string
	Exclude        []string
	DeleteExcluded bool
	DryRun         bool
}

func (c RcloneSyncOpts) AsFlagArray() []string {
	args := []string{"-v", "-P", "--no-update-dir-modtime", "--no-update-modtime", "--create-empty-src-dirs", "--check-first"}

	args = append(args, "--checkers", strconv.Itoa(c.Checkers))
	args = append(args, "--transfers", strconv.Itoa(c.Transfers))
	args = append(args, "--buffer-size", c.BufferSize)

	if len(c.Exclude) != 0 {
		args = append(args, "--exclude", strings.Join(c.Exclude, " "))
		if c.DeleteExcluded {
			args = append(args, "--delete-excluded")
		}
	}
	if c.DryRun {
		args = append(args, "--dry-run")
	}
	return args
}

type RcloneSyncParams struct {
	Pwd          string
	SrcDir       string
	Dest         string
	DryRun       bool
	TransferType TransferType
}

func RcloneSync(params RcloneSyncParams) error {

	st, err := os.Stat(params.SrcDir)
	if err != nil {
		return fmt.Errorf("could not stat: %w", err)
	}
	if !st.IsDir() {
		return ErrWontSyncFiles
	}

	var args []string

	switch params.TransferType {
	case Local:
		// rclone sync -v -P --checkers 16 --transfers 16 --check-first --no-update-dir-modtime --no-update-modtime --create-empty-src-dirs ./dir remote:/dir/
		opts := RcloneSyncOpts{
			BufferSize: "16Mi", // default
			Checkers:   16,
			Transfers:  16,
			DryRun:     params.DryRun,
			Exclude:    []string{"videos/backup.7z"},
		}
		args = opts.AsFlagArray()
	case Cloud:
		// clone sync -v -P --buffer-size 32Mi --checkers 50 --transfers 100 --check-first  --no-update-dir-modtime --no-update-modtime --create-empty-src-dirs --exclude **.git/** --delete-excluded ./dir/ remote:/dir/

		opts := RcloneSyncOpts{
			BufferSize: "32Mi",
			Checkers:   50,
			Transfers:  100, // backblaze deals very well with parallelism

			// remove all .git folders from cloud backup, since it's emmergency only
			Exclude:        []string{"**.git/**", "node_modules/**", "__pycache__/**", "videos/backup.7z"},
			DeleteExcluded: true,
		}
		args = opts.AsFlagArray()

	default:
		return ErrInvalidTransferType
	}

	printRcloneCommand("sync", args, params.SrcDir, params.Dest)

	cmd := buildRcloneCommandCommand("sync", params.Pwd, params.SrcDir, params.Dest, args)
	return cmd.Run()
}
