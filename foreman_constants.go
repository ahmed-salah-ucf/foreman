package main

import "time"

const (
	MaxNumServices = 1e6
	MaxSizeChannel = 1e6
	NumWorkersThreads = 10
	TickInterval = 10 * time.Second
)

const (
	procfile = "procfile.yaml"
)

const (
	notVisited NodeStatus = 0
	currentlyVisiting NodeStatus = 1
	visited NodeStatus = 2

	cycleStart string = "__CYCLESTART__"
	null string = "__NULL__"

	active ServiceStatus = "ACTIVE"
	inactive ServiceStatus = "INACTIVE"
)