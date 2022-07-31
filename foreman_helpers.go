package main

import "fmt"

type Foreman struct {
	procfile string
	services map[string]Service
	servicesGraph map[string][]string
}

func new() *Foreman{
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

func graphHasCycle() (bool, map[string]string) {
	return false, nil
}

func (foreman *Foreman) runServices() {
	if cycleExist, _ := graphHasCycle(); cycleExist {
		// To-Do
		fmt.Println("found cycle please fix: " + "cycle elements")
	}
}