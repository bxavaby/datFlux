package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"datflux/internal/entropy"
)

type SystemMonitor struct {
	CPUUsage        float64
	MemoryUsage     float64
	MemoryTotal     uint64
	MemoryUsed      uint64
	NetworkRxBytes  uint64
	NetworkTxBytes  uint64
	NetworkRxSpeed  float64
	NetworkTxSpeed  float64
	ActiveInterface string

	lastRxBytes    uint64
	lastTxBytes    uint64
	lastUpdateTime time.Time
}

func NewSystemMonitor() *SystemMonitor {
	monitor := &SystemMonitor{
		lastUpdateTime: time.Now(),
	}
	monitor.Update()
	return monitor
}

func (sm *SystemMonitor) Update() {
	sm.updateCPU()
	sm.updateMemory()
	sm.updateNetwork()
}

func (sm *SystemMonitor) GetEntropySource() entropy.EntropySource {
	return entropy.EntropySource{
		CPU:       sm.CPUUsage,
		Memory:    sm.MemoryUsage,
		NetworkRx: sm.NetworkRxSpeed,
		NetworkTx: sm.NetworkTxSpeed,
		Timestamp: time.Now().UnixNano(),
	}
}

func (sm *SystemMonitor) updateCPU() {
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		sm.CPUUsage = cpuPercent[0]
	}
}

func (sm *SystemMonitor) updateMemory() {
	memStats, err := mem.VirtualMemory()
	if err == nil {
		sm.MemoryUsage = memStats.UsedPercent
		sm.MemoryTotal = memStats.Total
		sm.MemoryUsed = memStats.Used
	}
}

func (sm *SystemMonitor) updateNetwork() {
	netStats, err := net.IOCounters(true)
	if err != nil {
		return
	}

	var activeIface net.IOCountersStat
	found := false

	for _, iface := range netStats {
		if iface.Name != "lo" {
			activeIface = iface
			sm.ActiveInterface = iface.Name
			found = true
			break
		}
	}

	if !found && len(netStats) > 0 {
		activeIface = netStats[0]
		sm.ActiveInterface = activeIface.Name
	} else if !found {
		return
	}

	now := time.Now()
	elapsed := now.Sub(sm.lastUpdateTime).Seconds()

	if elapsed > 0 && sm.lastRxBytes > 0 && sm.lastTxBytes > 0 {
		sm.NetworkRxSpeed = float64(activeIface.BytesRecv-sm.lastRxBytes) / elapsed
		sm.NetworkTxSpeed = float64(activeIface.BytesSent-sm.lastTxBytes) / elapsed
	}

	sm.NetworkRxBytes = activeIface.BytesRecv
	sm.NetworkTxBytes = activeIface.BytesSent
	sm.lastRxBytes = activeIface.BytesRecv
	sm.lastTxBytes = activeIface.BytesSent
	sm.lastUpdateTime = now
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)
}

func FormatSpeed(bytesPerSec float64) string {
	return FormatBytes(uint64(bytesPerSec)) + "/s"
}
