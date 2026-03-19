package collector

import (
	"github.com/shirou/gopsutil/v3/net"
)

type NetworkStats struct {
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

func GetNetwork() (*NetworkStats, error) {
	stats, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	if len(stats) == 0 {
		return &NetworkStats{}, nil
	}

	return &NetworkStats{
		BytesSent:   stats[0].BytesSent,
		BytesRecv:   stats[0].BytesRecv,
		PacketsSent: stats[0].PacketsSent,
		PacketsRecv: stats[0].PacketsRecv,
	}, nil
}
