package http

type Config struct {
	Port       uint       `env:"PORT,notEmpty"`
	HeaderAuth HeaderAuth `envPrefix:"HEADER_AUTH_"`
}

type HeaderAuth struct {
	Key   string `env:"KEY"`
	Value string `env:"VALUE"`
}
