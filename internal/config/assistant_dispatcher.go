package config

type AssistantDispatcher struct {
	Timeout      int              `yaml:"timeout"`
	WorkerNumber int              `yaml:"worker-number"`
	Handler      AssistantHandler `yaml:"handler"`
}
