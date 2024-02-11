package remind

import "time"

type Config struct {
	InitDelay    time.Duration `env:"INIT_DELAY" env-default:"3d"`
	Frequency    time.Duration `env:"FREQUENCY" env-default:"7d"`
	RunFrequency time.Duration `env:"RUN_FREQUENCY" env-default:"1h"`
}
