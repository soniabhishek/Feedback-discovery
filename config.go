package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var GlobalConfig *XMLConfig

type ValueAttr struct {
	Value string `xml:"value,attr"`
}

func (va ValueAttr) ToInt() int {
	val, err := strconv.Atoi(va.Value)
	if err != nil {
		panic(err)
	}
	return val
}

func (va ValueAttr) ToFloat() float64 {
	val, err := strconv.ParseFloat(va.Value, 64)
	if err != nil {
		panic(err)
	}
	return val
}
func (va ValueAttr) ToString() string {
	return va.Value
}

type TCPService struct {
	Name             ValueAttr
	IPAddress        ValueAttr
	Port             ValueAttr
	MaxConnections   ValueAttr
	ImportanceFactor ValueAttr
}

type CPU struct {
	ImportanceFactor ValueAttr
	ThresholdValue   ValueAttr
}

type RAM struct {
	ImportanceFactor ValueAttr
	ThresholdValue   ValueAttr
}

type XMLConfig struct {
	XMLName                           xml.Name `xml:"xml"`
	Cpu                               CPU
	Ram                               RAM
	TCPService                        []TCPService
	ReadAgentStatusFromConfig         ValueAttr
	ReadAgentStatusFromConfigInterval ValueAttr
	AgentStatus                       ValueAttr
	Interval                          ValueAttr
}

func readConfig() {
	xmlFile, err := os.Open("config.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	content, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = xml.Unmarshal(content, &GlobalConfig)
	if err != nil {
		fmt.Println("Error xmling file:", err)
		return
	}
}

func InitConfig() {
	readConfig()
	if GlobalConfig.ReadAgentStatusFromConfig.Value == "true" {
		interval, err := strconv.Atoi(GlobalConfig.ReadAgentStatusFromConfigInterval.Value)
		if err != nil {
			fmt.Println("Error format error:", err)
			return
		}
		tick := time.Tick(time.Duration(interval) * time.Second)
		go func() {
			for _ = range tick {
				readConfig()
			}
		}()
	}
}
