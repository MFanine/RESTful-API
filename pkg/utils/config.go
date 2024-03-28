package utils

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    EmailSenderName     string `mapstructure:"EMAIL_SENDER_NAME"`
    EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDRESS"`
    EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    if err = viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file, %s", err)
    }

    err = viper.Unmarshal(&config)
    return config, err
}
