package main

const (
	startCmd SubCommand = "start"
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