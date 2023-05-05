package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	Filename  string `mapstructure:"filename"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxAge    int    `mapstructure:"max_age"`
	MaxBackup int    `mapstructure:"max_backup"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxidelConns int    `mapstructure:"max_idel_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filePath string) (err error) {
	viper.SetConfigFile(filePath)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadConfig failed, err: %v\n", err)
		return
	}
	// Read the configuration information and send the sequence to the Conf variable
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err: %v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed ~")
		// Synchronized deserialization to Conf variables when profile changes occur
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err: %v\n", err)
		}
	})
	return
}
