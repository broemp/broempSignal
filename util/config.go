package util

import "github.com/spf13/viper"

// Config stores all settings
// Values are read by viper from a .env file or enviroment variables
type Config struct {
	// DB Settings
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	// Server Settings
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`

	// API Keys
	DiscordToken  string `mapstructure:"DISCORD_TOKEN"`
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
}

// LoadConfig reads configuration file/enviroment
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	setDefaults()
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	return
}

func setDefaults() {
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://broempSignal:broempSignal@localhost:5432/broempSignal?sslmode=disable")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
}
