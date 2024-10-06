package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"
)

// readConf вычиляет, от какого пользователя запущен процесс и в зависимости от результата, скармливает readConfFile тот
// или иной путь, по которому предположительно лежит конфиг.
func (cnf *Config) ReadConf() error {
	var (
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

		err = cnf.readConfFile(os.Getenv("HOME") + "/.reniced.yaml")
	} else {
		if runtime.GOOS == "freebsd" {
			err = cnf.readConfFile("/usr/local/etc/reniced.yaml")
		} else {
			err = cnf.readConfFile("/etc/reniced.yaml")
		}
	}

	return err
}

// readConfFile парсит даденный конфиг в глобальную структуру, с которой работает программа.
func (cnf *Config) readConfFile(path string) error {
	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		return err
	}

	// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
	if fileInfo.Size() > 65535 {
		err = fmt.Errorf("config file %s is too long for config, skipping", path)

		return err
	}

	buf, err := os.ReadFile(path)

	// Не удалось прочитать.
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(buf, &cnf); err != nil {
		return err
	}

	if cnf.LoopDelay == 0 {
		cnf.LoopDelay = 2000
		log.Printf("loop_delay set to %d milliseconds", cnf.LoopDelay)
	}

	if cnf.Debug {
		log.Println("debug set to true")
	}

	if cnf.MaxWorkers == 0 {
		cnf.MaxWorkers = 5
		log.Printf("max_workers set to %d", cnf.MaxWorkers)
	}

	if cnf.NbCapacity == 0 {
		cnf.NbCapacity = 100
		log.Printf("nb_capacity set to %d", cnf.NbCapacity)
	}

	if cnf.NbCapacity < cnf.MaxWorkers {
		cnf.NbCapacity = cnf.MaxWorkers
		log.Printf("nb_capacity set to same value as max_workers: %d", cnf.MaxWorkers)
	}

	// Обычный юзер не может задрать nice больше 0, только руту такое можно.
	var maxNice int

	if os.Getuid() == 0 {
		maxNice = -20
	}

	for _, v := range cnf.Prio {
		switch {
		case len(v.Name) != 0 && (v.Nice < 20 && v.Nice > maxNice):
			for _, procName := range v.Name {
				cnf.Renice[procName] = v.Nice

				if cnf.Debug {
					log.Printf("Add %s to list of nicelevel %d processes.", procName, v.Nice)
				}
			}
		case len(v.Name) == 0:
			log.Println("Skipping empty entry in prio config block")
		case v.Nice > 20 || v.Nice < maxNice:
			log.Printf("Niceness value out of range %d < nice < 20 in prio config block", maxNice)
		}
	}

	for _, v := range cnf.IOPrio {
		switch {
		case len(v.Name) == 0:
			log.Println("Skipping empty name entry in ioprio config block")
		case v.Class > 3:
			log.Printf("Skipping entry %s with class > 3 in ioprio config block", v.Name)
		case v.Class == 0 && v.Prio != 0:
			log.Printf("Skipping entry %s with class 0 and prio not 0 in ioprio config block", v.Name)
		case v.Class == 3 && v.Prio != 0:
			log.Printf("Skipping entry %s with class 3 and prio not 0 in ioprio config block", v.Name)
		case v.Class == 1 && v.Prio > 7:
			log.Printf("Skipping entry %s with class 1 and prio > 7 in ioprio config block", v.Name)
		case v.Class == 2 && v.Prio > 7:
			log.Printf("Skipping entry %s with class 2 and prio > 7 in ioprio config block", v.Name)
		default:
			for _, processName := range v.Name {
				cnf.IoreniceClass[processName] = v.Class
				cnf.IoreniceClassdata[processName] = v.Prio

				if cnf.Debug {
					log.Printf("Add %s to list of ionice class %d processes.", processName, v.Class)
					log.Printf("Add %s to list of ionice prio %d processes.", processName, v.Prio)
				}
			}
		}
	}

	for _, v := range cnf.Kill {
		if len(v.Name) != 0 {
			switch v.Sig {
			case "kill", "stop", "term", "int", "quit", "abrt", "hup", "usr1", "usr2":
				for _, procName := range v.Name {
					switch v.Sig {
					case "kill":
						cnf.KillSignal[procName] = syscall.SIGKILL
					case "stop":
						cnf.KillSignal[procName] = syscall.SIGSTOP
					case "term":
						cnf.KillSignal[procName] = syscall.SIGTERM
					case "int":
						cnf.KillSignal[procName] = syscall.SIGINT
					case "quit":
						cnf.KillSignal[procName] = syscall.SIGQUIT
					case "abrt":
						cnf.KillSignal[procName] = syscall.SIGABRT
					case "hup":
						cnf.KillSignal[procName] = syscall.SIGHUP
					case "usr1":
						cnf.KillSignal[procName] = syscall.SIGUSR1
					case "usr2":
						cnf.KillSignal[procName] = syscall.SIGUSR2
					default:
						continue
					}

					if cnf.Debug {
						log.Printf(
							"Add %s to list of processes that should be killed with SIG%s signal.",
							procName,
							strings.ToUpper(v.Sig),
						)
					}
				}
			default:
				log.Println("Skipping unsupported signal entry in kill config block")
			}
		} else {
			log.Println("Skipping empty entry in kill config block")
		}
	}

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
