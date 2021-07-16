package common

import (
	"os"
	"path"
	"strconv"
	"syscall"
)

func CreateFileWithPerm(filepath string, permcode string) (*os.File, error) {
	perm, err := strconv.ParseInt(permcode, 8, 64)
	if err != nil {
		return nil, err
	}
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	filedir := path.Dir(filepath)
	os.MkdirAll(filedir, os.FileMode(0770))
	fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(filepath, os.FileMode(perm))
	}
	return fd, err
}
