package safe

import (
	"fmt"
	"oneinstack/router/output"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func GetFirewallStatus() (*output.FirewallStatus, error) {
	var enabled bool
	var err error

	switch runtime.GOOS {
	case "darwin":
		enabled, err = isFirewallEnabledMac()
	case "linux":
		enabled, err = isFirewallEnabledLinux()
	default:
		return nil, fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return nil, err
	}

	return &output.FirewallStatus{Enabled: enabled}, nil
}

func GetFirewallPorts() (*output.FirewallPorts, error) {
	var ports []output.FirewallPort
	var err error

	switch runtime.GOOS {
	case "darwin":
		ports, err = getListeningPortsMac()
	case "linux":
		ports, err = getListeningPortsLinux()
	default:
		return nil, fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return nil, err
	}

	return &output.FirewallPorts{Ports: ports}, nil
}

func isFirewallEnabledMac() (bool, error) {
	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), "enabled"), nil
}

func getListeningPortsMac() ([]output.FirewallPort, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePortsMac(string(output)), nil
}

func isFirewallEnabledLinux() (bool, error) {
	cmd := exec.Command("sudo", "ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), "active"), nil
}

func getListeningPortsLinux() ([]output.FirewallPort, error) {
	cmd := exec.Command("ss", "-tulnp")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePortsLinux(string(output)), nil
}

func parsePortsLinux(s string) []output.FirewallPort {
	lines := strings.Split(s, "\n")
	var ports []output.FirewallPort
	re := regexp.MustCompile(`(tcp|udp)\s+[^\s]+\s+[^\s]+\s+([^\s]+):([^\s]+)\s+([^\s]+):(\*|[\d.]+|[\w]+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 6 {
			ports = append(ports, output.FirewallPort{
				Protocol:  matches[1],
				Port:      matches[3],
				State:     "LISTEN",
				Policy:    "ALLOW",
				Direction: "INBOUND",
				Source:    matches[4],
			})
		}
	}
	return ports
}

func parsePortsMac(s string) []output.FirewallPort {
	lines := strings.Split(s, "\n")
	var ports []output.FirewallPort
	re := regexp.MustCompile(`([^\s]+)\s+[^\s]+\s+[^\s]+\s+[^\s]+\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) == 5 {
			protocol := matches[2]
			port := strings.Split(matches[3], ":")[1]
			source := matches[4]
			ports = append(ports, output.FirewallPort{
				Protocol:  protocol,
				Port:      port,
				State:     "LISTEN",
				Policy:    "ALLOW",
				Direction: "INBOUND",
				Source:    source,
			})
		}
	}
	return ports
}
