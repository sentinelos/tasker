// Package host implementation of diagnostic host metadata
package host

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
	"github.com/sentinelos/tasker/internal/diagnostic/metadata"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func NewHostMetadata() *metadata.Metadata {
	return metadata.NewMetadata("host", "Host metadata", func() labels.Labels {
		l := labels.Labels{
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"platform":   platforms.DefaultString(),
			"go_version": strings.Replace(runtime.Version(), "go", "", -1),
			"user":       os.Getenv("USER"),
			"shell":      os.Getenv("SHELL"),
		}

		if cpuInfo, err := cpu.Info(); err == nil {
			l["cpu_total"] = strconv.Itoa(len(cpuInfo))
			for _, logical := range cpuInfo {
				if _, found := l["cpu_"+logical.PhysicalID]; !found {
					l["cpu_"+logical.PhysicalID] = logical.ModelName
				}
			}
		}

		if hostInfo, err := host.Info(); err == nil {
			l["hostname"] = hostInfo.Hostname
			l["platform"] = hostInfo.Platform
			l["platform_version"] = hostInfo.PlatformVersion
			l["kernel_version"] = hostInfo.KernelVersion
		}

		if memInfo, err := mem.VirtualMemory(); err == nil {
			l["memory"] = strconv.Itoa(int(memInfo.Total))
		}

		return l
	})
}

func init() {
	metadata.Add(NewHostMetadata())
}
