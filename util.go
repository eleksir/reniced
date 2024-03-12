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
		err := fmt.Errorf("Config file %s is too long for config, skipping", path)

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
		log.Println("cmd_delay can not be 0, setting to 200 milliseconds")

		cfg.CmdDelay = 200
	}

	if cfg.LoopDelay == 0 {
		log.Println("loop_delay can not be 0, setting to 2000 milliseconds")

		cfg.LoopDelay = 2000
	}

	if cfg.Debug {
		log.Println("debug set to true")
	} else {
		log.Println("debug set to false")
	}

	newPrio := make(Prio, 0)

	// Обычный юзер не может задрать nice больше 0, только руту такое можно.
	var maxNice int32

	if os.Getuid() == 0 {
		maxNice = -20
	}

	for _, v := range cfg.Prio {
		switch {
		case len(v.Name) != 0 && (v.Nice < 20 && v.Nice > maxNice):
			newPrio = append(newPrio, v)

			if cfg.Debug {
				for _, procName := range v.Name {
					log.Printf("Add %s to list of nicelevel %d processes.", procName, v.Nice)
				}
			}
		case len(v.Name) == 0:
			log.Println("Skipping empty entry in prio config block")
		case v.Nice > 20 || v.Nice < maxNice:
			log.Printf("Niceness value out of range %d < nice < 20 in prio config block", maxNice)
		}
	}

	cfg.Prio = newPrio

	newKill := make(Kill, 0)

	for _, v := range cfg.Kill {
		if len(v.Name) != 0 {
			switch v.Sig {
			case "kill", "stop", "term", "int", "quit", "abrt", "hup", "usr1", "usr2":
				newKill = append(newKill, v)

				if cfg.Debug {
					for _, procName := range v.Name {
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

	cfg.Kill = newKill

	return cfg, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
