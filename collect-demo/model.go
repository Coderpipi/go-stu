package main

import "github.com/shirou/gopsutil/v4/disk"

const (
	CPU  SysType = "cpu"
	Mem  SysType = "mem"
	Disk SysType = "disk"
	Net  SysType = "net"
)

type (
	SysType string

	Sys interface {
		Type() SysType
	}

	CPUInfo struct {
		Percent float64 `structs:"cpu_percent"`
	}

	MemInfo struct {
		// Total amount of RAM on this system
		Total uint64 `structs:"total"`

		// RAM available for programs to allocate
		//
		// This value is computed from the kernel specific values.
		Available uint64 `structs:"available"`

		// RAM used by programs
		//
		// This value is computed from the kernel specific values.
		Used uint64 `structs:"used"`

		// Percentage of RAM used by programs
		//
		// This value is computed from the kernel specific values.
		UsedPercent float64 `structs:"used_percent"`

		// OS X
		Active   uint64 `structs:"active"`
		Inactive uint64 `structs:"inactive"`
		Wired    uint64 `structs:"wired"`
	}

	NetInfo struct {
		NetName         string  `json:"-"`
		BytesSent       uint64  `structs:"-"`
		BytesRecv       uint64  `structs:"-"`
		PacketsSent     uint64  `structs:"-"`
		PacketsRecv     uint64  `structs:"-"`
		BytesSentRate   float64 `structs:"bytes_sent_rate"`
		BytesRecvRate   float64 `structs:"bytes_recv_rate"`
		PacketsSentRate float64 `structs:"packets_sent_rate"`
		PacketsRecvRate float64 `structs:"packets_recv_rate" `
	}

	DiskInfo struct {
		MountPoint string          `struct:"-"`
		UsageStat  *disk.UsageStat `struct:"-"`
	}
)

func (DiskInfo) Type() SysType {
	return Disk
}

func (NetInfo) Type() SysType {
	return Net
}

func (CPUInfo) Type() SysType {
	return CPU
}

func (MemInfo) Type() SysType {
	return Mem
}
