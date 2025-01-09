package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"oneinstack/app"
	"oneinstack/internal/services/software"
	"oneinstack/internal/services/user"
	"oneinstack/server"
	"oneinstack/web"
	"oneinstack/web/input"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func main() {
	//初始化服务
	server.Start()
	resetPwdCmd.Flags().StringP("name", "n", "", "username")
	resetPwdCmd.Flags().StringP("pwd", "p", "", "password")

	resetUserCmd.Flags().StringP("oldn", "", "", "old username")
	resetUserCmd.Flags().StringP("newn", "", "", "new username")
	// 将命令添加到根命令
	rootCmd.AddCommand(install)
	rootCmd.AddCommand(resetPwdCmd)
	rootCmd.AddCommand(resetUserCmd)
	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "one",
	Short: "oneinstack",
}

const pidFile = "server.pid" // 存储 PID 的文件路径

// serverStopCmd 定义启动和停止服务的命令
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start, restart, or stop HTTP server",
	Example: "go run main.go server [start|restart|stop]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			switch args[0] {
			case "start":
				startServer()
			case "restart":
				restartServer()
			case "stop":
				stopServer()
			default:
				fmt.Println("Invalid argument. Use 'start', 'restart', or 'stop'.")
			}
		} else {
			fmt.Println("Invalid argument. Use 'start', 'restart', or 'stop'.")
		}
	},
}

// startServer 启动服务并记录 PID
func startServer() {
	r := web.SetupRouter()
	fmt.Println("HTTP Server starting...")

	// 创建 PID 文件
	pid := os.Getpid()
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Fatalf("Failed to write PID file: %v", err)
	}

	// 启动服务
	if app.ONE_CONFIG.System.Port == "" {
		app.ONE_CONFIG.System.Port = "8089"
	}
	if err := r.Run("0.0.0.0:" + app.ONE_CONFIG.System.Port); err != nil {
		log.Fatal("Server run error:", err)
	}

	// 删除 PID 文件（在服务正常退出时）
	os.Remove(pidFile)
}

// restartServer 重启服务
func restartServer() {
	fmt.Println("Restarting HTTP Server...")

	// 调用 stopServer() 停止当前服务
	stopServer()

	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	// 启动新服务
	cmd := exec.Command(execPath, "server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to restart server: %v", err)
	}

	fmt.Printf("Server restarted successfully, new PID: %d\n", cmd.Process.Pid)
}

// stopServer 停止服务
func stopServer() {
	fmt.Println("Stopping HTTP Server...")

	// 读取 PID 文件
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No running server found.")
			return
		}
		log.Fatalf("Failed to read PID file: %v", err)
	}

	// 转换 PID 并发送终止信号
	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		log.Fatalf("Invalid PID in file: %v", err)
	}

	// 向目标进程发送 SIGTERM 信号
	err = syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}

	// 删除 PID 文件
	os.Remove(pidFile)
	fmt.Println("Server stopped successfully.")
}

var resetPwdCmd = &cobra.Command{
	Use:     "resetpwd",
	Short:   "reset user password",
	Example: " resetpwd -n admin -p 123123 ",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatalf("username not found")
		}
		pwd, _ := cmd.Flags().GetString("pwd")
		if name == "" {
			log.Fatalf("password not found")
		}
		err := user.ChangePassword(name, pwd)
		if err != nil {
			log.Fatalf("Add user error: %v", err)
		}
	},
}

var resetUserCmd = &cobra.Command{
	Use:     "resetUsername",
	Short:   "reset user username",
	Example: " resetUsername --oldn AHMPotFoxig --newn admin ",
	Run: func(cmd *cobra.Command, args []string) {
		on, _ := cmd.Flags().GetString("oldn")
		if on == "" {
			log.Fatalf("old username not found")
		}
		nn, _ := cmd.Flags().GetString("newn")
		if nn == "" {
			log.Fatalf("new username not found")
		}
		err := user.ResetUsername(on, nn)
		if err != nil {
			log.Fatalf("Add user error: %v", err)
		}
	},
}

var install = &cobra.Command{
	Use:     "install",
	Short:   "安装 php nginx  phpmyadmin",
	Example: "  install -s php",
	Run: func(cmd *cobra.Command, args []string) {
		ls := []*input.InstallParams{
			&input.InstallParams{
				Key:     "php",
				Version: "7.4",
			}, &input.InstallParams{
				Key:     "webserver",
				Version: "1.24.0",
			}, &input.InstallParams{
				Key:     "phpmyadmin",
				Version: "5.2.1",
			},
		}
		for _, v := range ls {
			fmt.Println("开始安装" + v.Key)
			op, err := software.NewInstallOP(v)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fn, err := op.Install(true)
			fmt.Println("开始安装：日志位于:", fn)
			if err != nil {
				fmt.Println("安装失败" + err.Error())
			}
		}

	},
}
