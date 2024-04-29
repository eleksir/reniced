//go:build linux
// +build linux

package main

import (
	"fmt"
	"log"
	"syscall"

	proc "github.com/shirou/gopsutil/v3/process"
)

// GetIOPriority достаёт приоритет вводв-вывода для укаанного процесса.
func GetIOPriority(which int, who int) (uint32, error) {
	var (
		err       error
		r0, _, e1 = syscall.Syscall(syscall.SYS_IOPRIO_GET, uintptr(which), uintptr(who), 0)
		prio      = uint32(r0)
	)

	if e1 != 0 {
		err = fmt.Errorf(fmt.Sprintf("Received error number %v", e1))
	}

	return prio, err
}

// GetIOPriority ставит приоритет вводв-вывода для укаанного процесса.
func SetIOPriority(which int, who int, prio uint32) error {
	var (
		err      error
		_, _, e1 = syscall.Syscall(syscall.SYS_IOPRIO_SET, uintptr(which), uintptr(who), uintptr(prio))
	)

	if e1 != 0 {
		err = fmt.Errorf(fmt.Sprintf("Received error number %v", e1))
	}

	return err
}

// IORenice выставляет процессу приоритет для операций ввода-вывода.
func IORenice(cnf Config, p *proc.Process, processName string) {
	if currentIOPrioLevel, err := GetIOPriority(IOPRIO_WHO_PROCESS, int(p.Pid)); err == nil {
		currentClass, currentClassdata := PrioToClassAndClassdata(currentIOPrioLevel)
		class := ioreniceClass[processName]
		classdata := ioreniceClassdata[processName]

		if class == currentClass && classdata == currentClassdata {
			if cnf.Debug && currentIOPrioLevel != 0 { // 0 - это дефолтный уровень io niceness.
				log.Printf(
					"IONiceness for %s(%d) already set to %d %d",
					processName,
					p.Pid,
					currentClass,
					currentClassdata,
				)
			}
		} else {
			iopriority := ClassAndClassdataToPrio(class, classdata)

			_ = SetIOPriority(IOPRIO_WHO_PROCESS, int(p.Pid), iopriority)

			if cnf.Debug {
				log.Printf(
					"Set ioniceness for %s(%d) to %d %d (was %d %d)",
					processName,
					p.Pid,
					class,
					classdata,
					currentClass,
					currentClassdata,
				)
			}
		}
	} else {
		log.Printf("Unable to get ioprio for %s(%d)", processName, p.Pid)
	}
}
