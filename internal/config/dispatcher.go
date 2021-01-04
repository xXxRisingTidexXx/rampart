package config

type Dispatcher struct {
	Token        string `yaml:"-"`
	Timeout      int    `yaml:"timeout"`
	WorkerNumber int    `yaml:"worker-number"`
}
