package system

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"oneinstack/web/output"
	"time"
)

var metrics = &output.SystemMetrics{
	NetworkIOMap: make(map[string][]output.IOStat),
	DiskIOMap:    make(map[string][]output.IOStat),
}

func GetSystemMonitor() *output.SystemMetrics {
	return metrics
}

func SystemMonitor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 清理超过30分钟的数据
		cleanupOldData()

		networkStats, _ := net.IOCounters(true)
		diskStats, _ := disk.IOCounters()

		// 收集所有网络接口的数据
		for _, network := range networkStats {
			stats, exists := metrics.NetworkIOMap[network.Name]
			if !exists {
				stats = []output.IOStat{}
			}
			stats = append(stats, output.IOStat{
				Time:      time.Now(),
				Time1:     time.Now().Format("15:04:05"),
				BytesSent: network.BytesSent,
				BytesRecv: network.BytesRecv,
			})
			metrics.NetworkIOMap[network.Name] = stats
		}

		// 收集所有磁盘的数据
		for name, disk := range diskStats {
			stats, exists := metrics.DiskIOMap[name]
			if !exists {
				stats = []output.IOStat{}
			}
			stats = append(stats, output.IOStat{
				Time:    time.Now(),
				Time1:   time.Now().Format("15:04:05"),
				Read:    disk.ReadBytes,
				Written: disk.WriteBytes,
			})
			metrics.DiskIOMap[name] = stats
		}
	}
}

func cleanupOldData() {
	threshold := time.Now().Add(-30 * time.Minute)

	// 移除超过30分钟的网络IO数据
	for name, stats := range metrics.NetworkIOMap {
		filtered := make([]output.IOStat, 0)
		for _, stat := range stats {
			if stat.Time.After(threshold) {
				filtered = append(filtered, stat)
			}
		}
		metrics.NetworkIOMap[name] = filtered
	}

	// 移除超过30分钟的磁盘IO数据
	for name, stats := range metrics.DiskIOMap {
		filtered := make([]output.IOStat, 0)
		for _, stat := range stats {
			if stat.Time.After(threshold) {
				filtered = append(filtered, stat)
			}
		}
		metrics.DiskIOMap[name] = filtered
	}
}

// GetSystemInfo 获取系统信息和磁盘使用情况
func GetSystemInfo() (*output.SystemInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	// CPU Info
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	// CPU Usage
	cpuPercent, _ := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}
	// Memory Info
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	// Disk Info
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, err
	}
	// Network Info
	interfaces, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}
	sysinfo := &output.SystemInfo{
		HostInfo:      info,
		CPUInfo:       cpuInfo,
		CPU:           cpuPercent,
		Memory:        vmem,
		DiskUsage:     partitions,
		NetIOCounters: interfaces,
	}
	return sysinfo, nil
}
