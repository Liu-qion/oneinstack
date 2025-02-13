package safe

import (
	"fmt"
	"oneinstack/app"
	"oneinstack/internal/models"
	"oneinstack/internal/services"
	"oneinstack/router/input"
	"oneinstack/router/output"
	"os/exec"
	"strings"
)

func GetUfwStatus() (*output.IptablesStatus, error) {
	// 获取 UFW 状态
	cmd := exec.Command("ufw", "status", "verbose")
	var out strings.Builder
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return &output.IptablesStatus{Enabled: false, PingBlocked: false}, nil
	}

	// 检查 UFW 是否启用
	enabled := strings.Contains(out.String(), "Status: active")

	// 检查是否有阻止 ICMP 请求的规则
	pingBlocked := strings.Contains(out.String(), "icmp") && strings.Contains(out.String(), "DENY")

	return &output.IptablesStatus{Enabled: enabled, PingBlocked: pingBlocked}, nil
}

func GetUfwRules(param *input.IptablesRuleParam) (*services.PaginatedResult[models.IptablesRule], error) {
	tx := app.DB()
	if param.Q != "" {
		tx = tx.Where("remark LIKE ?", "%"+param.Q+"%").Or("source LIKE ?", "%"+param.Q+"%").Or("dest LIKE ?", "%"+param.Q+"%")
	}
	if param.Target != "" {
		tx = tx.Where("target = ?", param.Target)
	}
	return services.Paginate[models.IptablesRule](tx, &models.IptablesRule{}, &input.Page{
		Page:     param.Page.Page,
		PageSize: param.Page.PageSize,
	})
}

func AddUfwRule(param *models.IptablesRule) error {
	if param.State == 0 {
		return fmt.Errorf("状态不能为禁用")
	}
	err := addUfwRule(param)
	if err != nil {
		return err
	}
	tx := app.DB().Create(param)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 删除UFW规则函数
func DeleteUfwRule(id int64) error {
	rule := &models.IptablesRule{}
	tx := app.DB().Where("id = ?", id).First(rule)
	if tx.Error != nil {
		return tx.Error
	}
	validProtocols := []string{"tcp", "udp", "icmp"}
	if !contains(validProtocols, rule.Protocol) {
		return fmt.Errorf("invalid protocol: %s. Valid options are: tcp, udp, icmp", rule.Protocol)
	}

	ipList := strings.Split(rule.IPs, ",")
	if len(ipList) == 0 || ipList[0] == "" {
		ipList = []string{"any"}
	}

	portList := strings.Split(rule.Ports, ",")
	if len(portList) == 0 || portList[0] == "" {
		portList = []string{"0"}
	}

	for _, ip := range ipList {
		for _, port := range portList {
			var cmdArgs []string
			cmdArgs = append(cmdArgs, "delete", "allow", rule.Direction) // 关键变化：添加delete参数

			if rule.Direction == "in" {
				cmdArgs = append(cmdArgs, "from", ip, "to", "any")
			} else if rule.Direction == "out" {
				cmdArgs = append(cmdArgs, "to", ip)
			} else {
				return fmt.Errorf("invalid direction: %s", rule.Direction)
			}

			if port != "0" {
				cmdArgs = append(cmdArgs, "port", port)
			}

			cmdArgs = append(cmdArgs, "proto", rule.Protocol)

			cmd := exec.Command("ufw", cmdArgs...)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to delete ufw rule: %v\nCommand: %s\nOutput: %s",
					err,
					strings.Join(cmdArgs, " "),
					string(output))
			}
			fmt.Printf("UFW rule deleted: ufw %s\n", strings.Join(cmdArgs, " "))
		}
	}
	return nil
}

// 更新UFW规则函数（需传入新旧规则）
func UpdateUfwRule(new *models.IptablesRule) error {
	oldRule := &models.IptablesRule{}
	tx := app.DB().Where("id = ?", new.ID).First(oldRule)
	if tx.Error != nil {
		return tx.Error
	}
	if err := DeleteUfwRule(oldRule.ID); err != nil {
		return fmt.Errorf("failed to remove old rule: %v", err)
	}
	if err := addUfwRule(new); err != nil {
		// 尝试恢复旧规则
		_ = addUfwRule(oldRule)
		return fmt.Errorf("failed to apply new rule: %v (rolled back)", err)
	}
	return nil
}

func addUfwRule(rule *models.IptablesRule) error {
	validProtocols := []string{"tcp", "udp", "icmp"}
	if !contains(validProtocols, rule.Protocol) {
		return fmt.Errorf("invalid protocol: %s. Valid options are: tcp, udp, icmp", rule.Protocol)
	}

	ipList := strings.Split(rule.IPs, ",")
	if len(ipList) == 0 || ipList[0] == "" {
		ipList = []string{"any"}
	}

	portList := strings.Split(rule.Ports, ",")
	if len(portList) == 0 || portList[0] == "" {
		portList = []string{"0"}
	}

	for _, ip := range ipList {
		for _, port := range portList {
			var cmdArgs []string
			cmdArgs = append(cmdArgs, "allow", rule.Direction)

			// 根据方向设置地址参数
			if rule.Direction == "in" {
				cmdArgs = append(cmdArgs, "from", ip, "to", "any")
			} else if rule.Direction == "out" {
				cmdArgs = append(cmdArgs, "to", ip)
			} else {
				return fmt.Errorf("invalid direction: %s", rule.Direction)
			}

			// 处理端口参数（0表示所有端口）
			if port != "0" {
				cmdArgs = append(cmdArgs, "port", port)
			}

			// 添加协议参数
			cmdArgs = append(cmdArgs, "proto", rule.Protocol)

			// 执行UFW命令
			cmd := exec.Command("ufw", cmdArgs...)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to add ufw rule: %v\nCommand: %s\nOutput: %s",
					err,
					strings.Join(cmdArgs, " "),
					string(output))
			}

			fmt.Printf("UFW rule added: ufw %s\n", strings.Join(cmdArgs, " "))
		}
	}

	return nil
}

// ToggleUfw 切换 ufw 的启用和禁用状态
func ToggleUfw() error {
	// 获取 ufw 当前的状态
	cmdStatus := exec.Command("ufw", "status")
	output, err := cmdStatus.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to check ufw status: %v, output: %s", err, string(output))
	}

	// 判断 UFW 当前状态，查找 "Status" 字段
	if strings.Contains(string(output), "Status: active") {
		// 如果当前是启用的，则禁用它
		cmdDisable := exec.Command("ufw", "disable")
		disableOutput, err := cmdDisable.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to disable ufw: %v, output: %s", err, string(disableOutput))
		}
	} else {
		// 如果当前是禁用的，则启用它
		cmdEnable := exec.Command("ufw", "enable")
		enableOutput, err := cmdEnable.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to enable ufw: %v, output: %s", err, string(enableOutput))
		}
	}

	return nil
}

// 检查是否已经禁用 ping（ICMP 请求）
func isPingBlocked() (bool, error) {
	// 查询现有的 ufw 状态，检查是否有阻止 ICMP 请求的规则
	cmd := exec.Command("ufw", "status", "verbose")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check ufw: %v, output: %s", err, string(output))
	}

	// 检查是否存在阻止 ICMP 请求的规则
	if strings.Contains(string(output), "icmp") && strings.Contains(string(output), "DENY") {
		return true, nil
	}

	return false, nil
}

// 禁止 ping 请求（ICMP）
func BlockPing() error {
	// 检查是否已经禁用 ping
	isBlocked, err := isPingBlocked()
	if err != nil {
		return err
	}
	// 如果已经禁用 ping，则删除现有规则
	if isBlocked {
		return deletePingBlockRule()
	}

	// 创建阻止 ICMP 请求的 ufw 规则
	cmd := exec.Command("ufw", "deny", "icmp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to block ping: %v, output: %s", err, string(output))
	}

	return nil
}

// 删除禁ping的规则
func deletePingBlockRule() error {
	// 列出当前的 ufw 规则并带有详细信息
	cmd := exec.Command("ufw", "status", "verbose")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list ufw rules: %v, output: %s", err, string(output))
	}

	// 查找禁ping规则
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// 查找包含禁ping的规则（icmp）
		if strings.Contains(line, "icmp") && strings.Contains(line, "DENY") {
			// 删除禁ping规则
			cmd := exec.Command("ufw", "delete", "deny", "icmp")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("failed to delete block ping rule: %v, output: %s", err, string(output))
			}
			break
		}
	}

	return nil
}

// contains 检查元素是否在列表中
func contains(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}
