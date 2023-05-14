package lxc

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

const (
	MaxLengthContainerId = 32
	ImageFolderPath      = "/var/lib/zind/images/base"
	RootFolderPathPrefix = "/var/lib/zind/containers/"
	LetterDict           = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func CreateChildProcess(args []string) error {
	newContainerId := generateContainerId(MaxLengthContainerId)
	newRootFolderPath := RootFolderPathPrefix + newContainerId
	if _, err := os.Stat(newRootFolderPath); os.IsNotExist(err) {
		if err := copyFileOrDirectory(ImageFolderPath, newRootFolderPath); err != nil {
			return err
		}
	}
	if err := syscall.Sethostname([]byte(newContainerId)); err != nil {
		return err
	}
	if err := syscall.Chroot(newRootFolderPath); err != nil {
		return err
	}
	if err := syscall.Chdir("/"); err != nil {
		return err
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return err
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	if err := syscall.Exec(path, args[0:], os.Environ()); err != nil {
		return err
	}
	return nil
}

func copyFileOrDirectory(src string, dst string) error {
	msg, err := os.Stat(src)
	if err != nil {
		return err
	}
	if msg.IsDir() {
		if err := os.MkdirAll(dst, 0777); err != nil {
			return err
		}
		if fileList, err := os.ReadDir(src); err == nil {
			for _, file := range fileList {
				if err = copyFileOrDirectory(filepath.Join(src, file.Name()), filepath.Join(dst, file.Name())); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		fileContent, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, fileContent, 0777); err != nil {
			return err
		}
	}
	return nil
}

func generateContainerId(n uint) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	containerId := make([]byte, n)
	length := len(LetterDict)
	for i := range containerId {
		containerId[i] = LetterDict[rand.Intn(length)]
	}
	return string(containerId)
}
