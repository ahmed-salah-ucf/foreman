package main

func (foreman *Foreman) buildServicesGraph() {
	for _, service := range foreman.services {
		foreman.servicesGraph[service.name] = service.info.deps
	}
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

func (foreman *Foreman) topoSortServices() []string {
	deps := make([]string, 0)
	visitingStatus := make(map[string]NodeStatus, 0)

	var topoSortDFS func (serviceName string)
	topoSortDFS = func (serviceName string) {
		if visitingStatus[serviceName] == visited {
			return
		}

		visitingStatus[serviceName] = visited
		
		for _, dep := range foreman.servicesGraph[serviceName] {
			topoSortDFS(dep)
		}

		deps = append(deps, serviceName)
	}

	for service := range foreman.services {
		topoSortDFS(service)
	}

	return deps
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