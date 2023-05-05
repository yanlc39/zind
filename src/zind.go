package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("Process PID => %v [%d]\n", os.Args, os.Getpid())
	switch os.Args[1] {
	case "run":
		Run()
	case "INIT_CONTAINER":
		initContainer()
	default:
		panic(os.Args[1] + " has not been defined!")
	}
}

func Run() {
	cmd := exec.Command(os.Args[0], "INIT_CONTAINER", os.Args[2])
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	err := cmd.Wait()
	if err != nil {
		return
	}
}

func initContainer() {
	containerId := GenerateContainerId(64)
	imageFolderPath := "/var/lib/zind/images/base"
	rootFolderPath := "/var/lib/zind/containers/" + containerId

	// copy on write

	if _, err := os.Stat(rootFolderPath); os.IsNotExist(err) {
		if err := copyFileOrDirectory(imageFolderPath, rootFolderPath); err != nil {
			panic(err)
		}
	}

	if err := syscall.Sethostname([]byte(containerId[0:15])); err != nil {
		panic(err)
	}

	if err := syscall.Chroot(rootFolderPath); err != nil {
		panic(err)
	}

	if err := syscall.Chdir("/"); err != nil {
		panic(err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		panic(err)
	}

	path, err := exec.LookPath(os.Args[2])

	if err != nil {
		panic(err)
	}

	if err := syscall.Exec(path, os.Args[2:], os.Environ()); err != nil {
		panic(err)
	}
}

func copyFileOrDirectory(src string, dst string) error {
	fmt.Printf("Copy %s => %s\n", src, dst)
	cmd := exec.Command("cp", "-r", src, dst)
	return cmd.Run()
}

func GenerateContainerId(n uint) string {
	rand.Seed(time.Now().UnixNano())
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	length := len(letters)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return string(b)
}
