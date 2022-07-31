package main

type Check struct {
	cmd string
	tcpPorts []string
	udpPorts []string
}

type ServiceInfo struct {
	cmd string
	runOnce bool
	checks Check
	deps []string
}

type Service struct {
	name string
	info ServiceInfo
}

// type Foreman struct {
// 	procfile string
// 	services map[string]Service
// 	servicesGraph map[string][]string
// }

type SubCommand string