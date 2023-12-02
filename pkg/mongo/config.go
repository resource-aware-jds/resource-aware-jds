package mongo

import (
	"fmt"
	"strings"
)

type Config struct {
	DatabaseName string   `envconfig:"DATABASE_NAME" default:"RAJDS"`
	Hosts        []string `envconfig:"HOSTS"`
	Username     string   `envconfig:"USERNAME"`
	Password     string   `envconfig:"PASSWORD"`
}

func (m Config) BuildConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, strings.Join(m.Hosts, ","))
}
