package main

import (
	"os"
	"time"
)

type Check struct {
	cmd string
	tcpPorts []string
	udpPorts []string
}

type ServiceInfo struct {
	cmd string
	runOnce bool
	checks Check
	deps []string
}

type Service struct {
	name string
	pid int
	info ServiceInfo
}

type Foreman struct {
	procfile string
	signalsChannel chan os.Signal
	servicesToRunChannel chan string
	checksTicker *time.Ticker
	services map[string]Service
	servicesGraph map[string][]string
}

type SubCommand string
type NodeStatus int