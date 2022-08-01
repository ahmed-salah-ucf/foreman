package main

import (
	"fmt"
	"os"
	"strings"
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
	ticker := time.NewTicker(500 * time.Millisecond)

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
		fmt.Println("run service " + serviceName)
	}
}

func (foreman *Foreman) checker() {
	for _, service := range foreman.services {
		runServiceChecks(service)
	}
}

func runServiceChecks(service Service) {
	
}