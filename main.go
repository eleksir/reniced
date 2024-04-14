package main

import (
	"log"
	"math"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/alitto/pond"
	"github.com/sevlyar/go-daemon"
	proc "github.com/shirou/gopsutil/v3/process"
)

// Пул воркеров.
var pool *pond.WorkerPool

// Канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
var sigChan = make(chan os.Signal, 1)

var kill = make(map[string]syscall.Signal)
var renice = make(map[string]int)

// main основная функция программы.
func main() {
	// Самое время поставить траппер сигналов.
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go sigHandler()

	var (
		cnf Config
		err error
	)

	u, err := user.Current()

	if err != nil {
		log.Fatalf("Unable to get user info for current user: %s", err)
	}

	// Для рута и для других пользователей конфиги находятся в разных местах.
	if u.Uid != "0" || runtime.GOOS == "darwin" {
		if os.Getenv("HOME") == "" {
			log.Fatalln("Unable to get HOME environment variable value")
		}

		cnf, err = readConf(os.Getenv("HOME") + "/.reniced.yaml")
	} else {
		cnf, err = readConf("/etc/reniced.yaml")
	}

	if err != nil {
		log.Fatalf("unable to parse config: %s", err)
	}

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
	reniced(cnf)
}

// reniced основная логика.
func reniced(cnf Config) {
	pool = pond.New(cnf.MaxWorkers, cnf.NbCapacity)

	for {
		processList, err := proc.Processes()

		if err != nil {
			log.Fatalf("Unable to get process list: %s", err)
		}

		// Не самый эффективный, зато работающий способ - простой перебор массивов.
		for _, p := range processList {
			ok := pool.TrySubmit(func() {
				if processName, err := p.Name(); err == nil {
					if killSignal := kill[processName]; killSignal != 0 {
						switch killSignal {
						case syscall.SIGSTOP:
							_ = p.SendSignal(killSignal)

							if cnf.Debug {
								log.Printf(
									"Matching processName with killProcessname (%s), SIG%s sent",
									processName,
									strings.ToUpper(killSignal.String()),
								)
							}

						default:
							_ = p.SendSignal(killSignal)

							if cnf.Debug {
								log.Printf(
									"Matching processName with killProcessname (%s), SIG%s sent",
									processName,
									strings.ToUpper(killSignal.String()),
								)
							}
						}
					}

					// Для каждого процесса извлекаем его текущий priority.
					if currentPrioLevel, err := syscall.Getpriority(syscall.PRIO_PROCESS, int(p.Pid)); err == nil {
						if niceLevel := renice[processName]; niceLevel != 0 {
							prioLevel := niceLevel
							currentNiceLevel := currentPrioLevel

							if runtime.GOOS == "linux" {
								// Value of PR = 20 + (-20 to +19) is 0 to 39
								switch {
								case niceLevel > 0:
									prioLevel = 20 - niceLevel
								case niceLevel < 0:
									prioLevel = 20 + int(math.Abs(float64(niceLevel)))
								default:
									prioLevel = 20
								}

								currentNiceLevel = 20 - currentPrioLevel
							}

							if currentPrioLevel != prioLevel {
								// on Linux Setpriority actually operates niceness value.
								_ = syscall.Setpriority(syscall.PRIO_PROCESS, int(p.Pid), niceLevel)

								if cnf.Debug {
									log.Printf(
										"Set niceness for %s(%d) to %d (was %d)",
										processName,
										p.Pid,
										niceLevel,
										currentNiceLevel,
									)
								}
							} else if cnf.Debug {
								log.Printf("Niceness for %d already set to %d", p.Pid, niceLevel)
							}
						}
					}
				}
			})

			if !ok {
				log.Printf("Unable to add task to pool")
			}

			time.Sleep(time.Millisecond * time.Duration(cnf.CmdDelay))
		}

		time.Sleep(time.Millisecond * time.Duration(cnf.LoopDelay))
	}
}

// sigHandler хэндлер сигналов. Работает на выходе приложения. Держит INT, TERM, QUIT.
func sigHandler() {
	for {
		var s = <-sigChan
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

		pool.StopAndWait()

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
