package collector

import (
	"github.com/shirou/gopsutil/v3/net"
)

type Port struct {
	Port   uint32
	PID    int32
	Type   string
	Status string
}

func GetPorts() ([]Port, error) {
	connections, err := net.Connections("all")
	if err != nil {
		return nil, err
	}

	var ports []Port
	seen := map[uint32]bool{}

	for _, c := range connections {
		if c.Laddr.Port == 0 || seen[c.Laddr.Port] {
			continue
		}
		seen[c.Laddr.Port] = true

		connType := "TCP"
		if c.Type == 2 {
			connType = "UDP"
		}

		ports = append(ports, Port{
			Port:   c.Laddr.Port,
			PID:    c.Pid,
			Type:   connType,
			Status: c.Status,
		})
	}

	return ports, nil
}
