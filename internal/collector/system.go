package collector

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	PID        int32
	Name       string
	CPUPercent float64
	MemMB      float32
	Status     string
}

type SystemSnapshot struct {
	CPUPercent  float64
	MemTotal    uint64
	MemUsed     uint64
	MemPercent  float64
	DiskTotal   uint64
	DiskUsed    uint64
	DiskPercent float64
	Hostname    string
	OS          string
	Uptime      uint64
	Processes   []Process
}

func GetSnapshot() (*SystemSnapshot, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	hostStat, err := host.Info()
	if err != nil {
		return nil, err
	}

	procs, _ := process.Processes()
	var processes []Process
	for _, p := range procs {
		name, _ := p.Name()
		cpu, _ := p.CPUPercent()
		mem, _ := p.MemoryInfo()
		status, _ := p.Status()
		var memMB float32
		if mem != nil {
			memMB = float32(mem.RSS) / 1024 / 1024
		}
		processes = append(processes, Process{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpu,
			MemMB:      memMB,
			Status:     status[0],
		})
	}

	return &SystemSnapshot{
		CPUPercent:  cpuPercent[0],
		MemTotal:    memStat.Total,
		MemUsed:     memStat.Used,
		MemPercent:  memStat.UsedPercent,
		DiskTotal:   diskStat.Total,
		DiskUsed:    diskStat.Used,
		DiskPercent: diskStat.UsedPercent,
		Hostname:    hostStat.Hostname,
		OS:          hostStat.OS,
		Uptime:      hostStat.Uptime,
		Processes:   processes,
	}, nil
}
