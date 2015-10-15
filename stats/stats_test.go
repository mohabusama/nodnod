package stats

import (
	"testing"
)

func TestGetStats(t *testing.T) {
	currentName := "stan"

	// Call to CPUPercent to avoid 0 CPUUsage on first call!
	GetStats(currentName)

	nodeStats, err := GetStats(currentName)

	if err != nil {
		t.Fatal("Failed to GetStats!", err)
	}

	if nodeStats.Name != currentName {
		t.Errorf("Invalid NodeStats address. Got: %s, Expected: %s",
			nodeStats.Name, currentName)
	}

	if nodeStats.Error != "" {
		t.Error("NodeStats had errors", nodeStats.Error)
	}

	if nodeStats.CPUUsed > 100 {
		t.Error("Invalid CPUUsed: ", nodeStats.CPUUsed)
	}

	if nodeStats.DiskUsedPercent <= 0 || nodeStats.DiskUsedPercent > 100 {
		t.Error("Invalid DiskUsedPercent: ", nodeStats.DiskUsedPercent)
	}

	if nodeStats.DiskTotal < nodeStats.DiskUsed {
		t.Errorf(
			"Invalid DiskTotal. Expected DiskTotal (%d) to be Greater than DiskUsed (%d)",
			nodeStats.DiskTotal, nodeStats.DiskUsed)
	}

	if nodeStats.MemoryTotal < nodeStats.MemoryUsed {
		t.Errorf(
			"Invalid MemoryTotal. Expected MemoryTotal (%d) to be Greater than MemoryUsed (%d)",
			nodeStats.MemoryTotal, nodeStats.MemoryUsed)
	}

	if nodeStats.MemoryUsedPercent <= 0 || nodeStats.MemoryUsedPercent > 100 {
		t.Error("Invalid MemoryUsedPercent: ", nodeStats.MemoryUsedPercent)
	}
}
