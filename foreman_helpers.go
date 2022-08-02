package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func new() *Foreman {
	foreman := Foreman {
		procfile: procfile,
		signalsChannel: make(chan os.Signal, MaxSizeChannel),
		servicesToRunChannel: make(chan string, MaxNumServices),
		checksTicker: time.NewTicker(TickInterval),
		services: map[string]Service{},
		servicesGraph: map[string][]string{},
	}
	return &foreman
}

func initForeman() *Foreman {
	foreman := new()
	foreman.parseProcfile()
	foreman.buildServicesGraph()
	foreman.signal()

	return foreman
}

func (foreman *Foreman) runServices() {
	if cycleExist, parentMap := graphHasCycle(foreman.servicesGraph); cycleExist {
		cycleElementsList := getCycleElements(parentMap)
		fmt.Printf("found cycle please fix: [%v]\n", strings.Join(cycleElementsList, ", "))
		os.Exit(1)
	}

	topologicallySortedServices := foreman.topoSortServices()
	
	

	foreman.createServiceRunners(foreman.servicesToRunChannel, NumWorkersThreads)
	sendServicesOnChannel(topologicallySortedServices, foreman.servicesToRunChannel)

	foreman.runPeriodicChecker(foreman.checksTicker)
}

// create a worker pool by starting up numWorkers workers threads
func (foreman *Foreman) createServiceRunners(services <-chan string, numWorkers int) {
	for w := 0; w < numWorkers; w++ {
		go foreman.serviceRunner(services)
	}
}

// Here’s the worker, of which we’ll run several concurrent instances.
func (foreman *Foreman) serviceRunner(services <-chan string) {
	for serviceName := range services {
		foreman.runService(serviceName)
	}
}

func (foreman *Foreman) runService(serviceName string) {
	service := foreman.services[serviceName]
	if (len(service.info.cmd)) > 0 {
		cmdName, cmdArgs := parseCmdLine(service.info.cmd)
		serviceCmd := exec.Command(cmdName, cmdArgs...)
		serviceCmd.Start()
		syscall.Setpgid(serviceCmd.Process.Pid, serviceCmd.Process.Pid)
		service.pid = serviceCmd.Process.Pid
		fmt.Printf("[%d] %s process started [%v]\n", service.pid, service.name, time.Now())
		foreman.services[serviceName] = service
	}
}

func sendServicesOnChannel(servicesList []string, servicesChannel chan<- string) {
	for _, service := range servicesList {
		servicesChannel <- service
	}
}

func (foreman *Foreman) runPeriodicChecker(ticker *time.Ticker) {
	for range ticker.C {
		go foreman.checker()
	}
}

func (foreman *Foreman) checker() {
	for _, service := range foreman.services {
		foreman.runServiceChecks(service)
	}
}

func (foreman *Foreman) runServiceChecks(service Service) {
	if len(service.info.checks.cmd) > 0 {
		cmdName, cmdArgs := parseCmdLine(service.info.checks.cmd)
		
		checkCmd := exec.Command(cmdName, cmdArgs...)

		if err := checkCmd.Run(); err != nil {
			if syscall.Kill(service.pid, syscall.SIGTERM); err != nil {
				syscall.Kill(service.pid, syscall.SIGKILL)
				return
			}
			fmt.Printf("[%d] %s process terminated as check [%v] failed\n", service.pid, service.name, service.info.checks.cmd)
			return
		}
	}
	if len(service.info.checks.tcpPorts) > 0 {
		for _, port := range service.info.checks.tcpPorts {
			address := "localhost:" + port
			_, err := net.Dial("tcp", address)
			if err != nil {
				if syscall.Kill(service.pid, syscall.SIGTERM); err != nil {
					syscall.Kill(service.pid, syscall.SIGKILL)
					return
				}
				fmt.Printf("[%d] %s process terminated as TCP port [%v] is not listening\n", service.pid, service.name, port)
				return
			}
		}
	}

	if len(service.info.checks.udpPorts) > 0 {
		for _, port := range service.info.checks.udpPorts {
			address := "localhost:" + port
			_, err := net.Dial("udp", address)
			if err != nil {
				if syscall.Kill(service.pid, syscall.SIGTERM); err != nil {
					syscall.Kill(service.pid, syscall.SIGKILL)
					return
				}
				fmt.Printf("[%d] %s process terminated as UDP port [%v] is not listening\n", service.pid, service.name, port)
				return
			}
		}
	}
}

func parseCmdLine(cmd string) (name string, arg []string) {
	cmdLine := strings.Split(cmd, " ")
	cmdName := cmdLine[0]
	cmdArgs := cmdLine[1:]

	return cmdName, cmdArgs
}