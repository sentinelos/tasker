package diagnostic

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func hostMetadata() map[string]string {
	meta := map[string]string{
		"go_version": strings.Replace(runtime.Version(), "go", "", -1),
	}

	if cpuInfo, err := cpu.Info(); err == nil {
		meta["cpu_total"] = strconv.Itoa(len(cpuInfo))
		for _, logical := range cpuInfo {
			if _, found := meta["cpu_"+logical.PhysicalID]; !found {
				meta["cpu_"+logical.PhysicalID] = logical.ModelName
			}
		}
	}

	if hostInfo, err := host.Info(); err == nil {
		meta["os"] = runtime.GOOS
		meta["arch"] = runtime.GOARCH
		meta["hostname"] = hostInfo.Hostname
		meta["platform"] = hostInfo.Platform
		meta["platform_version"] = hostInfo.PlatformVersion
		meta["kernel_version"] = hostInfo.KernelVersion
	}

	if memInfo, err := mem.VirtualMemory(); err == nil {
		meta["memory"] = strconv.Itoa(int(memInfo.Total))
	}

	return meta
}
