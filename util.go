package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// readConf парсит даденный конфиг в глобальную структуру, с которой работает программа.
func readConf(path string) (Config, error) {
	var cfg Config

	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		return Config{}, err
	}

	// Конфиг-файл длинноват для конфига, попробуем следующего кандидата.
	if fileInfo.Size() > 65535 {
		err := fmt.Errorf("config file %s is too long for config, skipping", path)

		return Config{}, err
	}

	buf, err := os.ReadFile(path)

	// Не удалось прочитать.
	if err != nil {
		return Config{}, err
	}

	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.CmdDelay == 0 {
		cfg.CmdDelay = 20
		log.Printf("cmd_delay can not be 0, setting to %d milliseconds", cfg.CmdDelay)
	}

	if cfg.LoopDelay == 0 {
		cfg.LoopDelay = 2000
		log.Printf("loop_delay set to %d milliseconds", cfg.LoopDelay)
	}

	if cfg.Debug {
		log.Println("debug set to true")
	}

	if cfg.MaxWorkers == 0 {
		cfg.MaxWorkers = 5
		log.Printf("max_workers set to %d", cfg.MaxWorkers)
	}

	if cfg.NbCapacity == 0 {
		cfg.NbCapacity = 100
		log.Printf("nb_capacity set to %d", cfg.NbCapacity)
	}

	if cfg.NbCapacity < cfg.MaxWorkers {
		cfg.NbCapacity = cfg.MaxWorkers
		log.Printf("nb_capacity set to same value as max_workers: %d", cfg.MaxWorkers)
	}

	// Обычный юзер не может задрать nice больше 0, только руту такое можно.
	var maxNice int

	if os.Getuid() == 0 {
		maxNice = -20
	}

	for _, v := range cfg.Prio {
		switch {
		case len(v.Name) != 0 && (v.Nice < 20 && v.Nice > maxNice):
			for _, procName := range v.Name {
				renice[procName] = v.Nice

				if cfg.Debug {
					log.Printf("Add %s to list of nicelevel %d processes.", procName, v.Nice)
				}
			}
		case len(v.Name) == 0:
			log.Println("Skipping empty entry in prio config block")
		case v.Nice > 20 || v.Nice < maxNice:
			log.Printf("Niceness value out of range %d < nice < 20 in prio config block", maxNice)
		}
	}

	for _, v := range cfg.Kill {
		if len(v.Name) != 0 {
			switch v.Sig {
			case "kill", "stop", "term", "int", "quit", "abrt", "hup", "usr1", "usr2":
				for _, procName := range v.Name {
					kill[procName] = v.Sig

					if cfg.Debug {
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

	return cfg, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
