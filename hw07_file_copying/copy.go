package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrCreateDestinationFile = errors.New("wrong destination path or other error")
	ErrSeek                  = errors.New("seek operation error. check offset value")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcStat, err := srcFile.Stat()
	if err != nil {
		return err
	}

	var sourceFileSize int64
	if sourceFileSize = srcStat.Size(); sourceFileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > sourceFileSize {
		return ErrOffsetExceedsFileSize
	}

	if _, err = srcFile.Seek(offset, io.SeekStart); err != nil {
		return ErrSeek
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return ErrCreateDestinationFile
	}
	defer dstFile.Close()

	if limit <= 0 || limit > sourceFileSize-offset {
		limit = sourceFileSize - offset
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(srcFile)
	if _, err = io.CopyN(dstFile, barReader, limit); err != nil {
		return err
	}
	defer bar.Finish()

	return dstFile.Close()
}
