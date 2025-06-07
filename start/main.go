package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "child" { //通过传入的参数来运行
		fmt.Println("child")
		if err := syscall.Sethostname([]byte("mini-container")); err != nil { //设置hostname
			log.Fatal("设置hostname出错", err)
		}

		if err := syscall.Chroot("/home/codebo/ubuntu-rootfs"); err != nil { //切换root
			log.Fatal("设置root", err)
		}

		if err := os.Chdir("/"); err != nil {
			log.Fatal("设置根目录", err)
		}

		if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
			log.Fatal("设置程序", err)
		}

		proc := os.Args[2:] //获取指定的运行程序名

		if len(proc) == 0 {
			log.Fatal("no process specified")
			return
		}

		cmd := exec.Command(proc[0], proc[1:]...)
		log.Printf("第一个参数 %s", proc[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if err := syscall.Unmount("/proc", 0); err != nil {
			log.Fatal(err)
			log.Println("容器程序启动异常")
		}
		log.Println("程序结束")

	} else {
		fmt.Printf("[Parent] PID : %d\n", os.Getpid())
		cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)

		//绑定标准输入 输出
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		// 设置要使用的新 namespace
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS | // 隔离 hostname
				syscall.CLONE_NEWPID | // 隔离 PID
				syscall.CLONE_NEWNS, // 隔离挂载点
			Unshareflags: syscall.CLONE_NEWNS, // 确保当前进程不继承父 mount namespace
		}

		// 运行子进程
		if err := cmd.Run(); err != nil {
			log.Fatalf("parent run child error: %v", err)
		}

		log.Println("父亲进程结束")

	}
}
