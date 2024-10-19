package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("wat shouild i do")
	}
	fmt.Println("Hello world!")
}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func child() {
	slog.Info("Running child process, with arguments: ", slog.String("args", strings.Join(os.Args[2:], ", ")), slog.String("pid", fmt.Sprintf("%d", os.Getpid())))

	slog.Info("Setting hostname in container to 'container'")
	must(syscall.Sethostname([]byte("container")))
	slog.Info("Setting up mount namespace")
	must(syscall.Chroot("/home/user/projects/montainerd/rootfs"))
	must(syscall.Chdir("/"))
	// slog.Info("Mounting rootfs to rootfs")
	// must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	// slog.Info("Creating rootfs/oldrootfs directory")
	// must(os.MkdirAll("rootfs/oldrootfs", 0700))
	// slog.Info("Pivot root to rootfs/oldrootfs")
	// must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	// slog.Info("Changing directory to /")
	syscall.Mount("proc", "proc", "proc", 0, "")
	slog.Info("Running command in child namespace")
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command in child namespace:", err)
		os.Exit(1)
	}
	syscall.Unmount("/proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
