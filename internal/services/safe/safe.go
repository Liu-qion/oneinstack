package safe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// FirewallInfo 存储防火墙信息
type FirewallInfo struct {
	Status      string         `json:"status"`
	PingBlocked bool           `json:"ping_blocked"`
	PortRules   []PortRuleInfo `json:"port_rules"`
}

// PortRuleInfo 存储单个端口规则信息
type PortRuleInfo struct {
	Rule string `json:"rule"`
}

// GetFirewallInfo 获取防火墙信息并格式化为 JSON
func GetFirewallInfo() (string, error) {
	var cmd *exec.Cmd
	var output bytes.Buffer

	// 根据操作系统选择命令
	if runtime.GOOS == "linux" {
		// Linux 使用 iptables
		cmd = exec.Command("bash", "-c", "sudo iptables -L -n -v; sudo iptables -L INPUT -n -v | grep 'DROP' || echo 'Ping is allowed'")
	} else if runtime.GOOS == "darwin" {
		// macOS 使用 pfctl
		cmd = exec.Command("bash", "-c", "sudo pfctl -sr; sudo pfctl -s info; sudo ipfw list")
	} else {
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令并获取输出
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error executing command: %v, output: %s", err, output.String())
	}

	// 获取端口规则
	portRules, err := GetPortRules()
	if err != nil {
		return "", err
	}

	// 创建防火墙信息结构体
	firewallInfo := FirewallInfo{
		Status:      output.String(),
		PingBlocked: isPingBlocked(output.String()),
		PortRules:   portRules,
	}

	// 格式化为 JSON
	jsonData, err := json.Marshal(firewallInfo)
	if err != nil {
		return "", fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return string(jsonData), nil
}

// GetPortRules 获取端口规则并格式化为 JSON
func GetPortRules() ([]PortRuleInfo, error) {
	var cmd *exec.Cmd
	var output bytes.Buffer

	if runtime.GOOS == "linux" {
		cmd = exec.Command("bash", "-c", "sudo iptables -L -n -v | grep 'dpt:'")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("bash", "-c", "sudo pfctl -sr | grep 'pass in' || sudo ipfw list")
	} else {
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令并获取输出
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error executing command: %v, output: %s", err, output.String())
	}

	// 解析端口规则
	var portRules []PortRuleInfo
	lines := bytes.Split(output.Bytes(), []byte{'\n'})
	for _, line := range lines {
		if len(line) > 0 {
			portRules = append(portRules, PortRuleInfo{Rule: string(line)})
		}
	}

	return portRules, nil
}

// isPingBlocked 检查是否禁用 ping
func isPingBlocked(output string) bool {
	return bytes.Contains([]byte(output), []byte("DROP"))
}

// GetFirewallStatus 获取防火墙状态，返回是否启用
func GetFirewallStatus() (bool, error) {
	var cmd *exec.Cmd
	var output bytes.Buffer

	// 根据操作系统选择命令
	if runtime.GOOS == "linux" {
		// Linux 使用 systemctl 检查防火墙状态
		cmd = exec.Command("bash", "-c", "sudo systemctl is-active firewalld || echo 'inactive'")
	} else if runtime.GOOS == "darwin" {
		// macOS 使用 pfctl 检查防火墙状态
		cmd = exec.Command("bash", "-c", "sudo pfctl -s info | grep 'Status'")
	} else {
		return false, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 执行命令并获取输出
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("error executing command: %v, output: %s", err, output.String())
	}

	// 判断防火墙是否启用
	if runtime.GOOS == "linux" {
		return bytes.Contains(output.Bytes(), []byte("active")), nil
	} else if runtime.GOOS == "darwin" {
		return bytes.Contains(output.Bytes(), []byte("Enabled")), nil
	}

	return false, nil
}
