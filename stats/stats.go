package stats

import (
	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

const (
	// Message types
	DISCOVERY = 0
	STAT      = 1
	STATALL   = 2

	// Stats
	ALL = 0
	// CPU  = 1
	// MEM  = 2
	// DISK = 3
)

type NodeStats struct {
	Address           string  `json:"address"`
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

type MessageResponse struct {
	Host string `json:"host"`

	Nodes map[string]NodeStats `json:"nodes"`

	Error string `json:"error"`
}

func GetStats(currentAddress string) (NodeStats, error) {
	nodeStats := NodeStats{Address: currentAddress}

	nodeStats.CPUUsed, _ = getCpuUsage()

	memStat, _ := mem.VirtualMemory()
	nodeStats.MemoryTotal = memStat.Total
	nodeStats.MemoryUsed = memStat.Total - memStat.Available
	nodeStats.MemoryUsedPercent = memStat.UsedPercent

	nodeStats.DiskTotal, nodeStats.DiskUsed, nodeStats.DiskUsedPercent, _ = getDiskUsage()

	return nodeStats, nil
}

func GetAllStats(allNodes []*Node) ([]NodeStats, error) {
	allStats := []NodeStats{}
	chStat := make(chan NodeStats)

	for _, node := range allNodes {

		go func(node *Node) {
			log.Debug("Stat for node: ", node.Address)

			if err := node.Stat(); err == nil {
				select {
				case nStat := <-node.Result:
					log.Debug("Received stats from: ", node.Address)
					chStat <- nStat
				case nErr := <-node.Error:
					log.Error("Received error from: ", node.Address)
					chStat <- NodeStats{
						Address: node.Address,
						Error:   nErr,
					}
				}
			} else {
				chStat <- NodeStats{
					Address: node.Address,
					Error:   err.Error(),
				}
			}
		}(node)
	}

	// All nodes stats should be available within 10 secs!
	timeout := time.After(10 * time.Second)
	statCount := 0
	timeoutErr := false

	for {

		if statCount == len(allNodes) || timeoutErr {
			break
		}

		select {
		case nodeStat := <-chStat:
			log.Debug("Received stat channel: ", nodeStat)
			allStats = append(allStats, nodeStat)
			statCount++
			log.Debug("Stats updated: ", allStats, statCount)
		case <-timeout:
			log.Warn("Get all stats timedout")
			timeoutErr = true
		}
	}

	return allStats, nil
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
