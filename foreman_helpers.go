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

// creates a new object of Foreman
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

// initForeman initialises a new foreman object and signals, then parses the procfile
// and then builds the dependency graph
func initForeman() *Foreman {
	foreman := new()
	foreman.signal()
	foreman.parseProcfile()
	foreman.buildServicesGraph()

	return foreman
}

// runServices runs services after sorting them topologically by first creating
// a pool of workers threads that recieve services from servicesToRunChannel.
// It spawns a periodic checker thread that run services checks after constant
// duration.
// If the grapth has cycles, the program aborts and prints an error message.
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

// createServiceRunners creates a worker pool by starting up numWorkers workers threads
func (foreman *Foreman) createServiceRunners(services <-chan string, numWorkers int) {
	for w := 0; w < numWorkers; w++ {
		go foreman.serviceRunner(services)
	}
}

// serviceRunner is the worker, of which we’ll run several concurrent instances.
func (foreman *Foreman) serviceRunner(services <-chan string) {
	for serviceName := range services {
		foreman.runService(serviceName)
	}
}

// serviceDepsAreAllActive checks if the all dependences of a service are active.
func (foreman *Foreman) serviceDepsAreAllActive(service Service) bool {
	for _, dep := range service.info.deps {
		if foreman.services[dep].status == inactive {
			foreman.restartService(dep)
			return false
		} 
	}
	return true
}

// runService run service by spawning a new process for this service.
// the new spawned process has a new process group id equals its pid.
func (foreman *Foreman) runService(serviceName string) {
	service := foreman.services[serviceName]
	if (len(service.info.cmd)) > 0 {
		if foreman.serviceDepsAreAllActive(service) {
			cmdName, cmdArgs := parseCmdLine(service.info.cmd)
			serviceCmd := exec.Command(cmdName, cmdArgs...)
			serviceCmd.Start()
			service.status = active
			service.pid = serviceCmd.Process.Pid
			syscall.Setpgid(service.pid, service.pid)
			fmt.Printf("[%d] %s process started [%v]\n", service.pid, service.name, time.Now())
			foreman.services[serviceName] = service
		}
	}
}

// sendServicesOnChannel helper function sends a list of services to a service channel.
func sendServicesOnChannel(servicesList []string, servicesChannel chan<- string) {
	for _, service := range servicesList {
		servicesChannel <- service
	}
}

// runPeriodicChecker runs a new foreman thread at every tick from the ticker.
func (foreman *Foreman) runPeriodicChecker(ticker *time.Ticker) {
	for range ticker.C {
		go foreman.checker()
	}
}

// checker the checker process that runs all the checks of all services.
func (foreman *Foreman) checker() {
	for _, service := range foreman.services {
		foreman.runServiceChecks(service)
	}
}

// runServiceChecks helper function runs the checks of a service.
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


// parseCmdLine helper function parses command line string
// into command name and list of args.
func parseCmdLine(cmd string) (name string, arg []string) {
	cmdLine := strings.Split(cmd, " ")
	cmdName := cmdLine[0]
	cmdArgs := cmdLine[1:]

	return cmdName, cmdArgs
}