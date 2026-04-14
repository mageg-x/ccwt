package service

import (
	"log"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo 系统信息
type SystemInfo struct {
	CPUPercent float64 `json:"cpu_percent"`
	MemTotal   uint64  `json:"mem_total"`
	MemUsed    uint64  `json:"mem_used"`
	MemPercent float64 `json:"mem_percent"`
	GoRoutines int     `json:"goroutines"`
}

// GetSystemInfo 获取系统资源信息
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		GoRoutines: runtime.NumGoroutine(),
	}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("获取CPU信息失败: %v", err)
	} else if len(cpuPercent) > 0 {
		info.CPUPercent = cpuPercent[0]
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("获取内存信息失败: %v", err)
	} else {
		info.MemTotal = memInfo.Total
		info.MemUsed = memInfo.Used
		info.MemPercent = memInfo.UsedPercent
	}

	return info
}
