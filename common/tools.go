package common

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
)

func CreateFileWithPerm(filePath string, permCode string) (*os.File, error) {

	if abs := filepath.IsAbs(filePath); !abs {
		return nil, errors.New("file path should be absolute path")
	}

	perm, err := strconv.ParseInt(permCode, 8, 64)
	if err != nil {
		return nil, err
	}
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	filedir := path.Dir(filePath)
	os.MkdirAll(filedir, os.FileMode(perm))
	fd, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if os.IsExist(err) {
		os.Chmod(filePath, os.FileMode(perm))
	}
	return fd, err
}

func GetConnectAddress(id string) string {
	return id + "-graphd-svc"
}
