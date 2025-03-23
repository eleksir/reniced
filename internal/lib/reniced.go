package lib

import (
	"log"
	"math"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	proc "github.com/shirou/gopsutil/v3/process"
)

// reniced основная логика.
func (cnf Config) Reniced() {
	var (
		processListUpdateTimer *time.Ticker
		processKillTimer       *time.Ticker
		processReniceTimer     *time.Ticker
		processIONiceTimer     *time.Ticker
		processListTimer       *time.Ticker
		processList            []*proc.Process
		err                    error
	)

	// Прежде всего настроим таймеры
	processListUpdateTimer = time.NewTicker(time.Duration(cnf.ProcessListUpdatePeriod) * time.Millisecond)
	processReniceTimer = time.NewTicker(time.Duration(cnf.ProcessNicePeriod) * time.Millisecond)
	processKillTimer = time.NewTicker(time.Duration(cnf.ProcessKillPeriod) * time.Millisecond)
	processIONiceTimer = time.NewTicker(time.Duration(cnf.ProcessIONicePeriod) * time.Millisecond)
	processListTimer = time.NewTicker(time.Duration(cnf.PeriodicProcessInvocationPeriod) * time.Millisecond)

	for {
		select {
		// Обновляем список процессов.
		case <-processListUpdateTimer.C:
			processList, err = proc.Processes()

			if err != nil {
				log.Fatalf("Unable to get process list: %s", err)
			}
		// Проходим по списку процессов и посылаем сигналы процессам.
		case <-processKillTimer.C:
			for _, p := range processList {
				ok := cnf.Pool.TrySubmit(func() {
					if processName, err := p.Name(); err == nil {
						// Посылаем процессу сигналы, если таковые есть в конфиге.
						if killSignal := cnf.KillSignal[processName]; killSignal != 0 {
							switch killSignal { //nolint:exhaustive
							case syscall.SIGSTOP:
								_ = p.SendSignal(killSignal)

								if cnf.Debug {
									log.Printf(
										"Matching processName with killProcessname %s(%d), SIG%s sent",
										processName,
										p.Pid,
										strings.ToUpper(killSignal.String()),
									)
								}

							default:
								_ = p.SendSignal(killSignal)

								if cnf.Debug {
									log.Printf(
										"Matching processName with killProcessname %s(%d), SIG%s sent",
										processName,
										p.Pid,
										strings.ToUpper(killSignal.String()),
									)
								}
							}
						}
					}
				})

				if !ok {
					log.Printf("Unable to add task to pool. Increase nb_capacity?")
				}
			}
		// Проходим по списку процессов и меняем ionice процессам.
		case <-processIONiceTimer.C:
			for _, p := range processList {
				ok := cnf.Pool.TrySubmit(func() {
					if processName, err := p.Name(); err == nil {
						cnf.IORenice(p, processName)
					}
				})

				if !ok {
					log.Printf("Unable to add task to pool. Increase nb_capacity?")
				}
			}
		// Проходим по списку того, что надо запустить и запускаем.
		case <-processListTimer.C:
			for _, v := range cnf.PeriodicProcesses {
				log.Printf("Invoking command: %s", v)
				exec.Command(v)
				time.Sleep(time.Duration(cnf.PeriodicProcessInvocationDelay) * time.Millisecond)
			}
		// Проходим по списку процессов и меняем ionice процессам.
		case <-processReniceTimer.C:
			// Не самый эффективный, зато работающий способ - простой перебор массивов.
			for _, p := range processList {
				ok := cnf.Pool.TrySubmit(func() {
					if processName, err := p.Name(); err == nil {
						// Для каждого процесса извлекаем его текущий priority.
						if currentPrioLevel, err := syscall.Getpriority(syscall.PRIO_PROCESS, int(p.Pid)); err == nil {
							if niceLevel := cnf.Renice[processName]; niceLevel != 0 {
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
									log.Printf(
										"Niceness for %s(%d) already set to %d",
										processName,
										p.Pid,
										niceLevel,
									)
								}
							}
						}
					}
				})

				if !ok {
					log.Printf("Unable to add task to pool. Increase nb_capacity?")
				}
			}
		}
	}

	// Мы сюда никогда не попадём.
	// pool.StopAndWait()
} //nolint:wsl

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
