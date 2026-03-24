package analyzer

import (
	"fmt"
	"math"
)

type HistoryPoint struct {
	CPUPercent  float64
	MemPercent  float64
	DiskPercent float64
}

type TrendDirection string

const (
	TrendUp     TrendDirection = "up"
	TrendDown   TrendDirection = "down"
	TrendStable TrendDirection = "stable"
)

type Finding struct {
	Level  string // "critical", "warning", "info"
	Title  string
	Detail string
}

type Report struct {
	CPUTrend  TrendDirection
	MemTrend  TrendDirection
	DiskTrend TrendDirection
	Findings  []Finding
}

// trend calculates direction of a series of values
func trend(values []float64) TrendDirection {
	if len(values) < 4 {
		return TrendStable
	}
	// compare average of first half vs second half
	mid := len(values) / 2
	var first, second float64
	for _, v := range values[:mid] {
		first += v
	}
	for _, v := range values[mid:] {
		second += v
	}
	first /= float64(mid)
	second /= float64(len(values) - mid)

	diff := second - first
	if math.Abs(diff) < 2.0 {
		return TrendStable
	}
	if diff > 0 {
		return TrendUp
	}
	return TrendDown
}

// average returns mean of a slice
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// sustained checks if values have been above threshold for all recent N points
func sustained(values []float64, threshold float64, n int) bool {
	if len(values) < n {
		return false
	}
	recent := values[len(values)-n:]
	for _, v := range recent {
		if v < threshold {
			return false
		}
	}
	return true
}

// Analyze takes a history of snapshots and returns a report
func Analyze(history []HistoryPoint) Report {
	if len(history) == 0 {
		return Report{}
	}

	cpuVals := make([]float64, len(history))
	memVals := make([]float64, len(history))
	diskVals := make([]float64, len(history))
	for i, h := range history {
		cpuVals[i] = h.CPUPercent
		memVals[i] = h.MemPercent
		diskVals[i] = h.DiskPercent
	}

	report := Report{
		CPUTrend:  trend(cpuVals),
		MemTrend:  trend(memVals),
		DiskTrend: trend(diskVals),
	}

	avgCPU := average(cpuVals)
	avgMem := average(memVals)
	latest := history[len(history)-1]

	// CPU findings
	if sustained(cpuVals, 90, 10) {
		report.Findings = append(report.Findings, Finding{
			Level:  "critical",
			Title:  "Sustained High CPU",
			Detail: fmt.Sprintf("CPU has been above 90%% for the last 20 seconds. Average: %.1f%%.", avgCPU),
		})
	} else if sustained(cpuVals, 70, 15) {
		report.Findings = append(report.Findings, Finding{
			Level:  "warning",
			Title:  "Elevated CPU Load",
			Detail: fmt.Sprintf("CPU has been above 70%% for the last 30 seconds. Average: %.1f%%.", avgCPU),
		})
	} else if report.CPUTrend == TrendUp && avgCPU > 50 {
		report.Findings = append(report.Findings, Finding{
			Level:  "warning",
			Title:  "CPU Trending Up",
			Detail: fmt.Sprintf("CPU usage is climbing. Currently at %.1f%%, up from %.1f%% recently.", latest.CPUPercent, average(cpuVals[:len(cpuVals)/2])),
		})
	}

	// Memory findings
	if sustained(memVals, 90, 5) {
		report.Findings = append(report.Findings, Finding{
			Level:  "critical",
			Title:  "Critical Memory Pressure",
			Detail: fmt.Sprintf("Memory has been above 90%% for the last 10 seconds. Average: %.1f%%.", avgMem),
		})
	} else if report.MemTrend == TrendUp && latest.MemPercent > 70 {
		report.Findings = append(report.Findings, Finding{
			Level:  "warning",
			Title:  "Memory Trending Up",
			Detail: fmt.Sprintf("Memory usage is climbing steadily. Currently at %.1f%%, was %.1f%% recently.", latest.MemPercent, average(memVals[:len(memVals)/2])),
		})
	}

	// Disk findings
	if latest.DiskPercent >= 90 {
		report.Findings = append(report.Findings, Finding{
			Level:  "critical",
			Title:  "Disk Space Critical",
			Detail: fmt.Sprintf("Disk usage is at %.1f%%. Free space is dangerously low.", latest.DiskPercent),
		})
	} else if report.DiskTrend == TrendUp && latest.DiskPercent > 75 {
		report.Findings = append(report.Findings, Finding{
			Level:  "warning",
			Title:  "Disk Usage Growing",
			Detail: fmt.Sprintf("Disk usage is trending upward and currently at %.1f%%.", latest.DiskPercent),
		})
	}

	return report
}
