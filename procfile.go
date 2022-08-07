package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// parseProcfile parses procfile file
func (foreman *Foreman) parseProcfile() error {
	yamlContentMap, err := foreman.yamlUnmarshal()
	if err != nil {
		return err
	}
	foreman.parseProcfileHelper(yamlContentMap)

	return nil
}

func (foreman *Foreman) parseProcfileHelper(yamlContentMap map[string]map[string]any) {
	for serviceName, info := range yamlContentMap {
		foreman.parseService(serviceName, info)
	}
}

// parseService helper function parses service in procfile
func (foreman *Foreman) parseService(serviceName string, serviceInfo map[string]any) {
	info := foreman.parseServiceInfo(serviceName, serviceInfo)

	foreman.services[serviceName] = Service{name: serviceName, info: info}
}

// parseServiceInfo helper function parses service info for service
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

// parseServiceInfoChecks helper function parses service info checks for service
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

// yamlUnmarshal unmarshals yaml procfile contents into go object
func (foreman *Foreman)yamlUnmarshal() (map[string]map[string]any, error) {
	yamlMap := make(map[string]map[string]any)

    data, _ := os.ReadFile(foreman.procfile)

    err := yaml.Unmarshal([]byte(data), yamlMap)

	return yamlMap, err
}