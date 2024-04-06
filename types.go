package main

type Config struct {
	// Weather we fork to backgoround or stay as-is.
	Daemon bool `yaml:"daemon,omitempty"`

	// Where to put pid-file.
	Pidfile string `yaml:"pidfile,omitempty"`

	// LoopDelay in milliseconds - whole set of action delay period.
	LoopDelay int `yaml:"loop_delay,omitempty"`

	// CmdDelay in milliseconds - delay between each change application.
	CmdDelay int `yaml:"cmd_delay,omitempty"`

	// Whether to print debug info on stdout.
	Debug bool `yaml:"debug,omitempty"`

	// Maximum workers in pool.
	MaxWorkers int `yaml:"max_workers,omitempty"`

	// Capacity of non-blocking tasks. All tasks that exceeding this capacity will block submission of new tasks until
	// queue has vacant place for new task.
	NbCapacity int `yaml:"nb_capacity,omitempty"`

	// В гошке yaml не позволяет сделать удобный парсер для конструкций, аналогичных json-овскому
	// {kill:{"stop":["proc1","proc2"], "term":["proc3","proc4"]}}
	// Поэтому мы просто возьмём массив и каждому элементу добавим тип сигнала.
	Kill Kill `yaml:"kill,omitempty"`

	Prio Prio `yaml:"prio,omitempty"`
}

type Prio []struct {
	Nice int32    `yaml:"nice,omitempty"`
	Name []string `yaml:"name,omitempty"`
}

type Kill []struct {
	Sig  string   `yaml:"sig,omitempty"`
	Name []string `yaml:"name,omitempty"`
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
