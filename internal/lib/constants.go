package lib

const IOPRIO_CLASS_NONE = uint32(0) //nolint:revive,stylecheck
const IOPRIO_CLASS_RT = uint32(1)   //nolint:revive,stylecheck
const IOPRIO_CLASS_BE = uint32(2)   //nolint:revive,stylecheck
const IOPRIO_CLASS_IDLE = uint32(3) //nolint:revive,stylecheck

const IOPRIO_WHO_PROCESS = 1 //nolint:revive,stylecheck
const IOPRIO_WHO_PGRP = 2    //nolint:revive,stylecheck
const IOPRIO_WHO_USER = 3    //nolint:revive,stylecheck

const IOPRIO_CLASS_SHIFT = uint32(13)                            //nolint:revive,stylecheck
const IOPRIO_PRIO_MASK = ((uint32(1) << IOPRIO_CLASS_SHIFT) - 1) //nolint:revive,stylecheck

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
