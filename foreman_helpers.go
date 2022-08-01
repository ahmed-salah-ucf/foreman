package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func new() *Foreman {
	foreman := Foreman{
		procfile: procfile,
		services: map[string]Service{},
		servicesGraph: map[string][]string{},
	}
	return &foreman
}



func initForeman() *Foreman {
	foreman := new()
	foreman.parseProcfile()
	foreman.buildServicesGraph()

	return foreman
}

func (foreman *Foreman) runServices() {
	if cycleExist, parentMap := graphHasCycle(foreman.servicesGraph); cycleExist {
		cycleElementsList := getCycleElements(parentMap)
		fmt.Printf("found cycle please fix: [%v]\n", strings.Join(cycleElementsList, ", "))
		os.Exit(1)
	}

	topologicallySortedServices := foreman.topoSortServices()
	servicesToRunChannel := make(chan string, MaxNumServices)
	ticker := time.NewTicker(5 * time.Second)

	foreman.createServiceRunners(servicesToRunChannel, NumWorkersThreads)
	
	sendServicesOnChannel(topologicallySortedServices, servicesToRunChannel)

	foreman.runPeriodicChecker(ticker)
}

func (foreman *Foreman) runPeriodicChecker(ticker *time.Ticker) {
	for range ticker.C {
		go foreman.checker()
	}
}

func sendServicesOnChannel(servicesList []string, servicesChannel chan<- string) {
	for _, service := range servicesList {
		servicesChannel <- service
	}
}

// create a worker pool by starting up numWorkers workers threads
func (foreman *Foreman) createServiceRunners(services <-chan string, numWorkers int) {
	for w := 0; w < numWorkers; w++ {
		go foreman.serviceRunner(services)
	}
}

// Here’s the worker, of which we’ll run several concurrent instances.
func (foreman *Foreman) serviceRunner(services <-chan string) {
	fmt.Println("enter service runner")
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
		service.pid = serviceCmd.Process.Pid
	}

	foreman.services[serviceName] = service
}

func parseCmdLine(cmd string) (name string, arg []string) {
	cmdLine := strings.Split(cmd, " ")
	cmdName := cmdLine[0]
	cmdArgs := cmdLine[1:]

	return cmdName, cmdArgs
}

func (foreman *Foreman) checker() {
	for _, service := range foreman.services {
		runServiceChecks(service)
	}
}

func runServiceChecks(service Service) {
	if len(service.info.checks.cmd) > 0 {
		cmdName, cmdArgs := parseCmdLine(service.info.checks.cmd)
		
		checkCmd := exec.Command(cmdName, cmdArgs...)

		if err := checkCmd.Run(); err != nil {
			if syscall.Kill(service.pid, syscall.SIGTERM); err != nil {
				syscall.Kill(service.pid, syscall.SIGKILL)
				return
			}
		}
	}

	if len(service.info.checks.tcpPorts) > 0 {
		// To-Do
	}

	if len(service.info.checks.udpPorts) > 0 {
		// To-Do
	}
}