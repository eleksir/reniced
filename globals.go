package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/alitto/pond"
)

// Пул воркеров.
var pool *pond.WorkerPool

// Канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
var sigChan = make(chan os.Signal, 1)

var kill = make(map[string]syscall.Signal)
var renice = make(map[string]int)
var ioreniceClass = make(map[string]uint32)
var ioreniceClassdata = make(map[string]uint32)

var NotImplementedError = fmt.Errorf("SYSCALL Not Implemented on current platform") //nolint:revive,stylecheck

const IOPRIO_CLASS_NONE = uint32(0) //nolint:revive,stylecheck
const IOPRIO_CLASS_RT = uint32(1)   //nolint:revive,stylecheck
const IOPRIO_CLASS_BE = uint32(2)   //nolint:revive,stylecheck
const IOPRIO_CLASS_IDLE = uint32(3) //nolint:revive,stylecheck

const IOPRIO_WHO_PROCESS = 1 //nolint:revive,stylecheck
const IOPRIO_WHO_PGRP = 2    //nolint:revive,stylecheck
const IOPRIO_WHO_USER = 3    //nolint:revive,stylecheck

const IOPRIO_CLASS_SHIFT = uint32(13)                            //nolint:revive,stylecheck
const IOPRIO_PRIO_MASK = ((uint32(1) << IOPRIO_CLASS_SHIFT) - 1) //nolint:revive,stylecheck

var ClassToString = map[uint32]string{
	IOPRIO_CLASS_NONE: "none",
	IOPRIO_CLASS_RT:   "realtime",
	IOPRIO_CLASS_BE:   "best-effort",
	IOPRIO_CLASS_IDLE: "idle",
}

var StringToClass = map[string]uint32{
	"none":        IOPRIO_CLASS_NONE,
	"0":           IOPRIO_CLASS_NONE,
	"realtime":    IOPRIO_CLASS_RT,
	"1":           IOPRIO_CLASS_RT,
	"best-effort": IOPRIO_CLASS_BE,
	"2":           IOPRIO_CLASS_BE,
	"idle":        IOPRIO_CLASS_IDLE,
	"3":           IOPRIO_CLASS_IDLE,
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
