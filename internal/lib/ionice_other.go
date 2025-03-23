//go:build !linux
// +build !linux

package lib

import proc "github.com/shirou/gopsutil/v3/process"

func (cnf Config) GetIOPriority(which int, who int) (prio uint32, err error) {
	return 0, cnf.NotImplementedError
}

func (cnf Config) SetIOPriority(which int, who int, prio uint32) (err error) {
	return cnf.NotImplementedError
}

func (cnf Config) IORenice(p *proc.Process, processName string) {}
