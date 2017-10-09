package main

type TCPService struct {
	Name             string
	IPAddress        string
	Port             string
	MaxConnections   string
	ImportanceFactor string
}

type CPU struct {
	ImportanceFactor string
	ThresholdValue   string
}

type RAM struct {
	ImportanceFactor string
	ThresholdValue   string
}

type XMLConfig struct {
	CPU                               CPU
	RAM                               RAM
	TCPService                        TCPService
	ReadAgentStatusFromConfig         string
	ReadAgentStatusFromConfigInterval string
	AgentStatus                       string
	Interval                          string
}

func InitConfig()  {
	
}
