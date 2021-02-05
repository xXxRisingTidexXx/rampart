package config

type Dispatcher struct {
	Timeout      int     `yaml:"timeout"`
	WorkerNumber int     `yaml:"worker-number"`
	Handler      Handler `yaml:"handler"`
}
