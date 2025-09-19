package cmd

import (
	"fmt"
	"os/user"
	"runtime"

	"goi_example/backend/server"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(HelpCmd) // 帮助信息
}

// 根命令
var RootCmd = &cobra.Command{
	Use:   "goi_example",
	Short: `Goi 示例项目`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}

var Help = `Goi 示例项目
`

// help
var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "help 帮助",
	RunE: func(cmd *cobra.Command, args []string) error {
		help_txt := fmt.Sprintf(Help)
		fmt.Print(help_txt)
		return nil
	},
}

func isRoot() bool {
	// 如果是 Windows 系统，不进行判断
	if runtime.GOOS == "windows" {
		return true
	}
	// 在 Unix 系统下检查是否为 root 用户
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}
