package output

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"time"
)

// SystemInfo 用于承载系统和磁盘相关信息的数据结构
type SystemInfo struct {
	HostInfo      *host.InfoStat         `json:"host_info"`
	CPU           []float64              `json:"cpu_usage"`
	CPUInfo       []cpu.InfoStat         `json:"cpu_info"`
	Memory        *mem.VirtualMemoryStat `json:"memory_usage"`
	DiskUsage     []disk.PartitionStat   `json:"disk_usage"`
	NetIOCounters []net.IOCountersStat   `json:"network_io"`
}

type IOStat struct {
	Time      time.Time `json:"-"`
	Time1     string    `json:"time"`
	BytesSent uint64    `json:"bytes_sent"`
	BytesRecv uint64    `json:"bytes_recv"`
	Read      uint64    `json:"read_bytes"`
	Written   uint64    `json:"write_bytes"`
}

// SystemMetrics 存储所有的统计数据
type SystemMetrics struct {
	NetworkIOMap map[string][]IOStat `json:"network_io"`
	DiskIOMap    map[string][]IOStat `json:"disk_io"`
}
