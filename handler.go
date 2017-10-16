package main

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net"
	"strconv"
)

const (
	Normal = "Normal"
	Halt   = "Halt"
	Drain  = "Drain"
	Down   = "Down"
)

var (
	initialRun = true
)

const (
	returnIdle = true
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	conn.Write(GetResponseForMode())
	conn.Close()
}

func GetResponseForMode() (response []byte) {
	ramThresholdValue := GlobalConfig.Ram.ThresholdValue.ToFloat()
	cpuThresholdValue := GlobalConfig.Cpu.ThresholdValue.ToFloat()
	cpuImportance := GlobalConfig.Cpu.ImportanceFactor.ToFloat()
	ramImportance := GlobalConfig.Ram.ImportanceFactor.ToFloat()

	switch GlobalConfig.AgentStatus.Value {
	case Normal:
		cpuLoad, err := cpu.Percent(0, false)
		if err != nil {
			return []byte("0%\n")
		}
		v, err := mem.VirtualMemory()
		if err != nil {
			return []byte("0%\n")
		}
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

		for _, tcpService := range GlobalConfig.TCPService {
			sessionOccupied := GetSessionUtilized(tcpService.IPAddress.Value, tcpService.Port.Value, tcpService.MaxConnections.ToInt())

			utilization = utilization + sessionOccupied*tcpService.ImportanceFactor.ToFloat()
			if tcpService.ImportanceFactor.ToFloat() > 0 {
				divider++
			}

			if sessionOccupied > 99 && tcpService.ImportanceFactor.ToFloat() == 1 {
				response = []byte("0%\n")
				break
			}
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

func GetSessionUtilized(IPAddress, servicePort string, maxNumberOfSessionsPerService int) (result float64) {
	numberOfEstablishedConnections := getNumberOfLocalEstablishedConnections(IPAddress, servicePort)
	if numberOfEstablishedConnections > 0 && maxNumberOfSessionsPerService > 0 {
		result = float64(maxNumberOfSessionsPerService) / float64(numberOfEstablishedConnections)
	}
	return
}

func getNumberOfLocalEstablishedConnections(ipAddress string, port string) int {
	if ipAddress == "*" {
		ipAddress = ""
	}
	result := runcmd("netstat -nlt | grep -w " + ipAddress + ":" + port + "  | grep ESTABLISHED | wc -l")
	count, err := strconv.Atoi(string(bytes.TrimSpace(result)))
	if err != nil {
		return 0
	}
	return count
}
