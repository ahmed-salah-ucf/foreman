package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (foreman *Foreman) signal () {
	signal.Notify(foreman.signalsChannel, syscall.SIGCHLD)
	go foreman.receiveSignals(foreman.signalsChannel)
}

func (foreman *Foreman) receiveSignals(sigChannel <-chan os.Signal) {
	for sig := range sigChannel {
		switch sig {
		case syscall.SIGCHLD:
			foreman.sigchldHandler()
		}
	}
}

func (foreman *Foreman) sigchldHandler() {
	for _, service := range foreman.services {
		p, _ := os.FindProcess(service.pid)
		stat, _ := p.Wait()
		if stat.ExitCode() != -1 {
			fmt.Printf("[%d] %s process terminated [%v]\n", service.pid, service.name, time.Now())
			if !service.info.runOnce {
				fmt.Printf("[%d] %s process restarted [%v]\n", service.pid, service.name, time.Now())
				foreman.restartService(service.name)
			}
		}
	}
}

func (foreman *Foreman) restartService(service string) {
	foreman.servicesToRunChannel <- service
}