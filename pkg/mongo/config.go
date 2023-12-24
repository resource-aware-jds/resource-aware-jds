package mongo

import (
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	DatabaseName string   `envconfig:"DATABASE_NAME" default:"rajds"`
	Hosts        []string `envconfig:"HOSTS"`
	Username     string   `envconfig:"USERNAME"`
	Password     string   `envconfig:"PASSWORD"`
}

func (m Config) BuildConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", url.QueryEscape(m.Username), url.QueryEscape(m.Password), strings.Join(m.Hosts, ","))
}
