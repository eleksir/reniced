//go:build !linux
// +build !linux

package main

import proc "github.com/shirou/gopsutil/v3/process"

func GetIOPriority(which int, who int) (prio uint32, err error) {
	return 0, NotImplementedError
}

func SetIOPriority(which int, who int, prio uint32) (err error) {
	return NotImplementedError
}

func IORenice(cnf Config, p *proc.Process, processName string) {}
