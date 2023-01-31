// Package metadata implementation of diagnostic metadata
package metadata

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/sentinelos/tasker/internal/diagnostic/labels"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func NewMetadata(name, description string, detector Detector) *Metadata {
	return &Metadata{
		Name:        name,
		Description: description,
		Labels:      detector(),
	}
}

// GetLabel returns the value in the map for the provided label
func (md *Metadata) GetLabel(name string) string {
	return md.Labels.Get(name)
}

func NewHostMetadata() *Metadata {
	return NewMetadata("host", "Host metadata", func() labels.Labels {
		workdir, _ := os.Getwd()
		l := labels.Labels{
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"platform":   platforms.DefaultString(),
			"go_version": strings.Replace(runtime.Version(), "go", "", -1),
			"user":       os.Getenv("USER"),
			"shell":      os.Getenv("SHELL"),
			"workdir":    workdir,
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
			l["id"] = hostInfo.HostID
			l["hostname"] = hostInfo.Hostname
			l["distro_name"] = hostInfo.Platform
			l["distro_family"] = hostInfo.PlatformFamily
			l["distro_version"] = hostInfo.PlatformVersion
			l["kernel_version"] = hostInfo.KernelVersion
		}

		if memInfo, err := mem.VirtualMemory(); err == nil {
			l["memory_total"] = strconv.Itoa(int(memInfo.Total))
		}

		return l
	})
}

func init() {
	Add(NewHostMetadata())
}
