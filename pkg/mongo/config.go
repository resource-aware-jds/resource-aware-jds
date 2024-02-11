package mongo

type Config struct {
	DatabaseName     string `envconfig:"DATABASE_NAME" default:"rajds"`
	ConnectionString string `envconfig:"CONNECTION_STRING"`
}
