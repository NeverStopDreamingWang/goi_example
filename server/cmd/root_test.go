package cmd

import (
	"fmt"
	"os"
	"testing"
)

// 调试
func TestCmd(t *testing.T) {
	// 在测试中执行命令
	// 服务管理
	// os.Args = []string{"main_test", "version"} // 查看面板版本信息
	// os.Args = []string{"main_test", "info"} // 查看面板信息
	// os.Args = []string{"main_test", "port"} // 修改面板端口
	// os.Args = []string{"main_test", "domain"} // 修改面板绑定域名
	// os.Args = []string{"main_test", "entrance"} // 修改面板安全入口
	// os.Args = []string{"main_test", "listen-ip", "all"} // 修改面板监听所有网络协议（IPv4,IPv6）
	// os.Args = []string{"main_test", "listen-ip", "ipv4"} // 修改面板监听 ipv4
	// os.Args = []string{"main_test", "listen-ip", "ipv6"} // 修改面板监听 ipv6
	// os.Args = []string{"main_test", "ssl", "enable"} // 开启面板 SSL
	// os.Args = []string{"main_test", "ssl", "disable"} // 管理面板 SSL
	// 用户管理
	// os.Args = []string{"main_test", "user", "info", "2", "10"} // 查看所有用户信息，可选参数 page pagesize
	// os.Args = []string{"main_test", "user", "username"} // 修改用户名
	// os.Args = []string{"main_test", "user", "password"} // 修改用户密码
	RootCmd.AddCommand(HelpCmd) // 帮助信息
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
