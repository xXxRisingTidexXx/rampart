package config

import (
	"fmt"
)

type Migrations struct {
	DSNParams map[string]string `yaml:"dsnParams"`
}

func (migrations *Migrations) String() string {
	return fmt.Sprintf("{%v}", migrations.DSNParams)
}
