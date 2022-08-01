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

func graphHasCycle(servicesGraph map[string][]string) (bool, map[string]string) {
	var hasCycle bool = false
	var parentMap = map[string]string{}
	var visitingStatus = map[string]NodeStatus{}
	for node := range servicesGraph {
		visitingStatus[node] = notVisited
	}

	var hasCycleDFS func (string)
	hasCycleDFS = func (node string)  {
		if visitingStatus[node] == visited {
			return
		} else if visitingStatus[node] == currentlyVisiting {
			hasCycle = true
			parentMap[cycleStart] = node
			return
		}

		visitingStatus[node] = currentlyVisiting
		for _, dep := range servicesGraph[node] {
			parentMap[dep] = node
			hasCycleDFS(dep)
		}

		visitingStatus[node] = visited
	}

	for node := range servicesGraph {
		parentMap[node] = null
		hasCycleDFS(node)

		if hasCycle {
			return true, parentMap
		}
	}

	return false, parentMap
}

func getCycleElements(parentMap map[string]string) []string {
	cycleElements := make([]string, 0)
	start := parentMap[cycleStart]
	cycleElements = append(cycleElements, start)

	next := parentMap[start]
	for start != next {
		cycleElements = append(cycleElements, next)
		next = parentMap[next]
	}

	return cycleElements
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