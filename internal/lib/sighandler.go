package lib

import (
	"log"
	"os"
	"syscall"
)

// sigHandler хэндлер сигналов. Работает на выходе приложения. Держит INT, TERM, QUIT.
func (cnf Config) SigHandler() {
	for {
		var s = <-cnf.SigChan
		switch s {
		case syscall.SIGINT:
			log.Print("Got SIGINT, quitting")
		case syscall.SIGTERM:
			log.Print("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			log.Print("Got SIGQUIT, quitting")

		// Заходим на новую итерацию, если у нас "неинтересный" сигнал
		default:
			continue
		}

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
