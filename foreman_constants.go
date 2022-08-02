package main

import "time"

const (
	MaxNumServices = 10
	MaxSizeChannel = 30
	NumWorkersThreads = 10
	TickInterval = 100 * time.Millisecond
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
)