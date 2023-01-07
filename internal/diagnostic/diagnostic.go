// Package diagnostic implementation of error handler
package diagnostic

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	DefaultSeverity  = DiagInfo
	DefaultTimeStamp = time.RFC3339
)

func NewDefaultDiagnostic() *Diagnostic {
	return &Diagnostic{
		Name:        "Default",
		Description: "",
		Severity:    DefaultSeverity,
		Writers: []Writer{
			NewConsoleWriter(),
		},
		Meta: HostMetadata(),
	}
}

func (d *Diagnostic) Trace(message string, fields ...Field) {
	d.Write(DiagTrace, message, fields...)
}

func (d *Diagnostic) Debug(message string, fields ...Field) {
	d.Write(DiagDebug, message, fields...)
}

func (d *Diagnostic) Info(message string, fields ...Field) {
	d.Write(DiagInfo, message, fields...)
}

func (d *Diagnostic) Warn(message string, fields ...Field) {
	d.Write(DiagWarn, message, fields...)
}

func (d *Diagnostic) Error(message string, fields ...Field) {
	d.Write(DiagError, message, fields...)
}

func (d *Diagnostic) Fatal(message string, fields ...Field) {
	d.Write(DiagFatal, message, fields...)
}

func (d *Diagnostic) Write(severity Severity, message string, fields ...Field) {
	if d.Severity >= severity {
		for _, writer := range d.Writers {
			writer.WriteEntry(&Entry{
				Severity: severity,
				Message:  message,
				Fields:   fields,
				Time:     time.Now(),
			})
		}
	}
}

func HostMetadata() map[string]string {
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
