package raspik

import (
	"syscall"
)

const (
	SI_LOAD_SHIFT = 16
)

type Getter interface {
	Get() error
}

type Load struct {
	// 1 minute load averages
	One float64
	// 5 minute load averages
	Five float64
	// 15 minute load averages
	Fifteen float64
}

type Uptime struct {
	// Seconds since boot
	Uptime uint64
}

type Mem struct {
	// Total usable main memory size
	TotalRam uint64
	// Available memory size
	FreeRam uint64
	// Amount of shared memory
	SharedRam uint64
	// Memory used by buffers
	BufferRam uint64
}

type Swap struct {
	// Total swap space size
	TotalSwap uint64
	// Used swap space
	UsedSwap uint64
	// Swap space still available
	FreeSwap uint64
}

func (load *Load) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	load.One = float64(sysinfo.Loads[0]) / (1 << SI_LOAD_SHIFT)
	load.Five = float64(sysinfo.Loads[1]) / (1 << SI_LOAD_SHIFT)
	load.Fifteen = float64(sysinfo.Loads[2]) / (1 << SI_LOAD_SHIFT)

	return nil
}

func (uptime *Uptime) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	uptime.Uptime = uint64(sysinfo.Uptime)

	/*
	   days := sysinfo.Uptime / 86400
	   hours := (sysinfo.Uptime / 3600) - (days * 24)
	   mins := (sysinfo.Uptime / 60) - (days * 1440) - (hours * 60)
	*/

	return nil
}

func (mem *Mem) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	mem.TotalRam = uint64(sysinfo.Totalram)
	mem.FreeRam = uint64(sysinfo.Freeram)
	mem.SharedRam = uint64(sysinfo.Sharedram)
	mem.BufferRam = uint64(sysinfo.Bufferram)

	return nil
}

func (swap *Swap) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	swap.TotalSwap = uint64(sysinfo.Totalswap)
	swap.FreeSwap = uint64(sysinfo.Freeswap)
	swap.UsedSwap = swap.TotalSwap - swap.FreeSwap

	return nil
}
