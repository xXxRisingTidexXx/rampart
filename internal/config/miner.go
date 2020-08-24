package config

type Miner interface {
	Name() string
	Schedule() string
	Metrics() *Server
}
