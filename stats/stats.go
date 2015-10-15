package stats

import (
	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const (
	// Message types
	STAT    = 0
	STATALL = 1

	// Stats
	ALL = 0
	// CPU  = 1
	// MEM  = 2
	// DISK = 3
)

type NodeStats struct {
	Name              string  `json:"name"`
	CPUUsed           float64 `json:"cpuUsed"`
	DiskTotal         uint64  `json:"diskTotal"`
	DiskUsed          uint64  `json:"diskUsed"`
	DiskUsedPercent   float64 `json:"diskUsedPercent"`
	MemoryTotal       uint64  `json:"memoryTotal"`
	MemoryUsed        uint64  `json:"memoryUsed"`
	MemoryUsedPercent float64 `json:"memoryUsedPercent"`
	Error             string  `json:"error"`
}

type MessageRequest struct {
	Type int `json:"type"`

	StatType int `json:"statType"`
}

type Stats map[string]NodeStats

type MessageResponse struct {
	Host string `json:"host"`

	Nodes Stats `json:"nodes"`

	Error string `json:"error"`
}

func GetStats(name string) (NodeStats, error) {
	nodeStats := NodeStats{Name: name}

	nodeStats.CPUUsed, _ = getCpuUsage()

	memStat, _ := mem.VirtualMemory()
	nodeStats.MemoryTotal = memStat.Total
	nodeStats.MemoryUsed = memStat.Total - memStat.Available
	nodeStats.MemoryUsedPercent = memStat.UsedPercent

	nodeStats.DiskTotal, nodeStats.DiskUsed, nodeStats.DiskUsedPercent, _ = getDiskUsage()

	return nodeStats, nil
}

func getCpuUsage() (float64, error) {
	if res, err := cpu.CPUPercent(0, false); err == nil {
		return res[0], nil
	}

	return 0.0, nil
}

func getDiskUsage() (total uint64, used uint64, usedPercent float64, err error) {
	if partitions, pErr := disk.DiskPartitions(true); pErr == nil {
		for _, partition := range partitions {
			if diskStat, statErr := disk.DiskUsage(partition.Mountpoint); statErr == nil {
				total += diskStat.Total
				used += diskStat.Used
			} else {
				log.Error("Failed to get diskStat:", statErr)
				err = statErr
			}
		}

		if total > 0 {
			usedPercent = float64(used) / float64(total) * 100
		}
	} else {
		log.Error("Failed to get partitions:", pErr)
		err = pErr
	}

	return
}
