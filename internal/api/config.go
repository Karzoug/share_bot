package api

type Config struct {
	Token string `env:"TOKEN,notEmpty"`
}
