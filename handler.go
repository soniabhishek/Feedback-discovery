package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net"
)

type Mode uint8

const (
	Normal Mode = iota
	Halt
	Drain
	Down
)

const (
	defaultMode       = Normal
	cpuThresholdValue = 100
	ramThresholdValue = 100
	cpuImportance     = 1
	ramImportance     = 0
	returnIdle        = true
	initialRun        = true
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	conn.Write(GetResponseForMode())
	conn.Close()
}

func GetResponseForMode() (response []byte) {
	switch defaultMode {
	case Normal:
		cpuLoad, _ := cpu.Percent(0, false)
		v, _ := mem.VirtualMemory()

		// If any resource is important and utilized 100% then everything else is not important
		if cpuLoad[0] > cpuThresholdValue && cpuThresholdValue > 0 || (v.UsedPercent > ramThresholdValue && ramThresholdValue > 0) {
			response = []byte("0%\n")
		}
		utilization := 0.0
		divider := 0.0
		utilization = utilization + cpuLoad[0]*cpuImportance
		if cpuImportance > 0 {
			divider++
		}

		utilization = utilization + v.UsedPercent*ramImportance
		if ramImportance > 0 {
			divider++
		}
		utilization = utilization / divider

		if utilization < 0 {
			utilization = 0
		}
		if utilization > 100 {
			utilization = 100
		}
		if returnIdle {
			response = []byte(fmt.Sprintf("%v\n", 100-utilization))
		} else {
			response = []byte(fmt.Sprintf("%v\n", utilization))
		}
		if initialRun {
			response = append([]byte("up ready "), response...)
		}
	case Drain:
		response = []byte("drain\n")
	case Halt:
		response = []byte("down\n")
	default:
		response = []byte("error\n")
	}
	return
}
