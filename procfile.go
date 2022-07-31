package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func (foreman *Foreman) parseProcfile() {
	yamlContentMap := foreman.yamlUnmarshal()
	
	foreman.parseProcfileHelper(yamlContentMap)
}

func (foreman *Foreman) parseProcfileHelper(yamlContentMap map[string]map[string]any) {
	for serviceName, info := range yamlContentMap {
		foreman.parseService(serviceName, info)
	}
}

func (foreman *Foreman) parseService(serviceName string, serviceInfo map[string]any) {
	info := foreman.parseServiceInfo(serviceName, serviceInfo)

	foreman.services[serviceName] = Service{name: serviceName, info: info}
}

func (foreman *Foreman) parseServiceInfo(serviceName string, serviceInfo map[string]any) ServiceInfo {
	info := ServiceInfo{}
	for key, value := range serviceInfo {
		switch key {
		case "cmd":
			info.cmd = value.(string)
		case "run_once":
			info.runOnce = value.(bool)
		case "deps":
			for _, dep := range value.([]any) {
				info.deps = append(info.deps, dep.(string))
			}
		case "checks":
			info.checks = foreman.parseServiceInfoChecks(value)
		}
	}

	return info
}

func (foreman *Foreman) parseServiceInfoChecks(value any) Check {
	checks := Check{}
	for checkKey, checkValue := range value.(map[string]any) {
		switch checkKey {
		case "cmd":
			checks.cmd = checkValue.(string)
		case "tcp_ports":
			for _, port := range checkValue.([]any) {
				checks.tcpPorts = append(checks.tcpPorts, fmt.Sprint(port.(int)))
			}
		case "udp_ports":
			for _, port := range checkValue.([]any) {
				checks.udpPorts = append(checks.udpPorts, fmt.Sprint(port.(int)))
			}
		}
	}
	return checks
}

func (foreman *Foreman)yamlUnmarshal() map[string]map[string]any {
	yamlMap := make(map[string]map[string]any)

    data, _ := os.ReadFile(foreman.procfile)

    err := yaml.Unmarshal([]byte(data), yamlMap)
    if err != nil {
        panic(err)
    }

	return yamlMap
}