package config

type CommonConfigModel struct {
	DebugMode bool `envconfig:"DEBUG_MODE" default:"false"`
}
