package main

import (
	"fmt"
	"os"
)

func new() *Foreman {
	foreman := Foreman{
		procfile: procfile,
		services: map[string]Service{},
		servicesGraph: map[string][]string{},
	}
	return &foreman
}

func (foreman *Foreman) buildServicesGraph() {
	for _, service := range foreman.services {
		foreman.servicesGraph[service.name] = service.info.deps
	}
}

func initForeman() *Foreman {
	foreman := new()
	foreman.parseProcfile()
	foreman.buildServicesGraph()

	return foreman
}

func graphHasCycle(servicesGraph map[string][]string) (bool, map[string]string) {
	// To-Do
	return false, nil
}

func (foreman *Foreman) topoSortServices() []string {
	// To-Do
	return nil
}

func servicePassChecks(serviceChecks Check) bool {
	// To-Do
	return false
}

func (foreman *Foreman) runService(serviceName string) {
	if !servicePassChecks(foreman.services[serviceName].info.checks) {
		fmt.Println("service {service name} doesn't pass check: " + "check value")
		os.Exit(1)
	}

	// To-Do: actual run of service...
}

func (foreman *Foreman) runServices() {
	if cycleExist, _ := graphHasCycle(foreman.servicesGraph); cycleExist {
		// To-Do
		fmt.Println("found cycle please fix: " + "cycle elements")
		os.Exit(1)
	}

	topologicallySortedServices := foreman.topoSortServices()

	for _, serviceName := range topologicallySortedServices {
		foreman.runService(serviceName)
	}
}