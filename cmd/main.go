package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"oneinstack/app"
	"oneinstack/internal/services/user"
	"oneinstack/server"
	"oneinstack/web"
	"os"
)

func main() {
	//初始化服务
	server.Start()
	// 为 "hello" 命令添加参数
	adduserCmd.Flags().StringP("name", "n", "", "username")
	adduserCmd.Flags().StringP("pwd", "p", "", "password")

	// 将命令添加到根命令
	rootCmd.AddCommand(adduserCmd)
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

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "start http server",
	Example: " go run main.go server",
	Run: func(cmd *cobra.Command, args []string) {
		r := web.SetupRouter()
		if err := r.Run("0.0.0.0:" + app.ONE_CONFIG.System.Port); err != nil {
			log.Fatal("Server run error:", err)
		}
	},
}

var adduserCmd = &cobra.Command{
	Use:     "addu",
	Short:   "add user",
	Example: " go run main.go addu -n abc -p 123 ",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatalf("username not found")
		}
		pwd, _ := cmd.Flags().GetString("pwd")
		if name == "" {
			log.Fatalf("password not found")
		}
		err := user.CreateUser(name, pwd, false)
		if err != nil {
			log.Fatalf("Add user error: %v", err)
		}
	},
}
