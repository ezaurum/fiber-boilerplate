package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func New() *Config {
	v := Config{
		Viper: viper.New(),
	}
	v.SetDefault("Port", 3000)
	v.SetConfigName(".env")
	v.SetConfigType("dotenv")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to read config file %v", err))
	}
	return &v
}

func (c *Config) setDefaults() {
	c.SetDefault("PORT", 3000)
	c.SetDefault("PREFORK", false)
}

func (c *Config) Port() int {
	return c.GetInt("PORT")
}

func (c *Config) ListenString() string {
	return fmt.Sprintf(":%d", c.Port())
}