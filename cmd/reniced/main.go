package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	lib "reniced/internal/lib"

	"github.com/alitto/pond"
	"github.com/sevlyar/go-daemon"
)

// main основная функция программы.
func main() {
	var (
		err error
		cnf = lib.Config{
			SigChan:             make(chan os.Signal, 1),
			KillSignal:          make(map[string]syscall.Signal),
			Renice:              make(map[string]int),
			IoreniceClass:       make(map[string]uint32),
			IoreniceClassdata:   make(map[string]uint32),
			NotImplementedError: fmt.Errorf("SYSCALL Not Implemented on current platform"),
			ClassToString: map[uint32]string{
				lib.IOPRIO_CLASS_NONE: "none",
				lib.IOPRIO_CLASS_RT:   "realtime",
				lib.IOPRIO_CLASS_BE:   "best-effort",
				lib.IOPRIO_CLASS_IDLE: "idle",
			},
			StringToClass: map[string]uint32{
				"none":        lib.IOPRIO_CLASS_NONE,
				"0":           lib.IOPRIO_CLASS_NONE,
				"realtime":    lib.IOPRIO_CLASS_RT,
				"1":           lib.IOPRIO_CLASS_RT,
				"best-effort": lib.IOPRIO_CLASS_BE,
				"2":           lib.IOPRIO_CLASS_BE,
				"idle":        lib.IOPRIO_CLASS_IDLE,
				"3":           lib.IOPRIO_CLASS_IDLE,
			},
		}
	)

	// Самое время поставить траппер сигналов.
	signal.Notify(cnf.SigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go cnf.SigHandler()

	err = cnf.ReadConf()

	if err != nil {
		log.Fatalf("unable to parse config: %s", err)
	}

	// Создаём пул со статическим количеством воркеров и длиной очереди неблокируемых задач, равной NbCapacity.
	cnf.Pool = pond.New(
		cnf.MaxWorkers,
		cnf.NbCapacity,
		pond.Strategy(pond.Lazy()),
	)

	defer cnf.Pool.StopAndWait()

	if cnf.Daemon {
		cntxt := &daemon.Context{
			WorkDir: "/",
			Args:    []string{"reniced"},
		}

		if cnf.Pidfile != "" {
			cntxt.PidFileName = cnf.Pidfile
		}

		if d, err := cntxt.Reborn(); err != nil {
			log.Fatal("Unable to run: ", err)

			return
		} else if d != nil {
			return
		}

		defer cntxt.Release() //nolint: errcheck
	}

	// То, ради чего всё затевалось.
	cnf.Reniced()
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
