package main

import "errors"

type TransferType int

const (
	Cloud TransferType = iota
	Local
)

var (
	ErrWontCopyDir         = errors.New("copy is only used with files")
	ErrInvalidTransferType = errors.New("transfer type can only be cloud or hdd")
	ErrWontSyncFiles       = errors.New("sync is only used with directories")
)
