package about

import (
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	Hostname     string    `json:"hostname"`
	OS           string    `json:"os"`
	Architecture string    `json:"architecture"`
	NumCPU       int       `json:"num_cpu"`
	CPUInfo      []CPUInfo `json:"cpu_info"`
	Memory       MemInfo   `json:"memory"`
	Disk         DiskInfo  `json:"disk"`
	GoVersion    string    `json:"go_version"`
	StartTime    time.Time `json:"start_time"`
}

type CPUInfo struct {
	Model    string  `json:"model"`
	Cores    int     `json:"cores"`
	Mhz      float64 `json:"mhz"`
	UserTime float64 `json:"user_time"`
	SysTime  float64 `json:"sys_time"`
	IdleTime float64 `json:"idle_time"`
}

type MemInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	UsageRate float64 `json:"usage_rate"`
}

type DiskInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	UsageRate float64 `json:"usage_rate"`
}

func GetSystemInfo() (*SystemInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpuTimes, err := cpu.Times(false)
	if err != nil {
		return nil, err
	}

	cpuDetails := make([]CPUInfo, len(cpuInfo))
	for i, info := range cpuInfo {
		cpuDetails[i] = CPUInfo{
			Model: info.ModelName,
			Cores: runtime.NumCPU(),
			Mhz:   info.Mhz,
		}
		if i < len(cpuTimes) {
			cpuDetails[i].UserTime = cpuTimes[i].User
			cpuDetails[i].SysTime = cpuTimes[i].System
			cpuDetails[i].IdleTime = cpuTimes[i].Idle
		}
	}

	return &SystemInfo{
		Hostname:     hostname,
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		CPUInfo:      cpuDetails,
		Memory: MemInfo{
			Total:     memInfo.Total,
			Used:      memInfo.Used,
			Free:      memInfo.Free,
			UsageRate: memInfo.UsedPercent,
		},
		Disk: DiskInfo{
			Total:     diskInfo.Total,
			Used:      diskInfo.Used,
			Free:      diskInfo.Free,
			UsageRate: diskInfo.UsedPercent,
		},
		GoVersion: runtime.Version(),
		StartTime: time.Now(),
	}, nil
}
