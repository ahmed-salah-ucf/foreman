package main

import (
	"fmt"
	"os"
	"strings"
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

	for _, serviceName := range topologicallySortedServices {
		foreman.runService(serviceName)
	}
}


func (foreman *Foreman) runService(serviceName string) {
	if !servicePassChecks(foreman.services[serviceName].info.checks) {
		fmt.Println("service {service name} doesn't pass check: " + "check value")
		os.Exit(1)
	}

	// To-Do: actual run of service...
}

func servicePassChecks(serviceChecks Check) bool {
	// To-Do
	return false
}