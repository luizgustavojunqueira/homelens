package client

import (
	"homelens/shared"

	"github.com/shirou/gopsutil/v4/host"
)

func readHostInfo() shared.HostInfo {
	info, err := host.Info()
	if err != nil {
		return shared.HostInfo{}
	}
	return shared.HostInfo{
		Hostname:      info.Hostname,
		OS:            info.OS,
		Platform:      info.Platform,
		KernelVersion: info.KernelVersion,
		Uptime:        info.Uptime,
	}
}
