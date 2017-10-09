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
		//TODO : How to handle error
		cpuLoad, _ := cpu.Percent(0, false)
		v, _ := mem.VirtualMemory()
		averageCpuLoad := cpuLoad[0]
		usedRam := v.UsedPercent
		// If any resource is important and utilized 100% then everything else is not important
		if averageCpuLoad > cpuThresholdValue && cpuThresholdValue > 0 || (usedRam > ramThresholdValue && ramThresholdValue > 0) {
			response = []byte("0%\n")
		}
		utilization := 0.0
		divider := 0.0
		utilization = utilization + averageCpuLoad*cpuImportance
		if cpuImportance > 0 {
			divider++
		}

		utilization = utilization + usedRam*ramImportance
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
