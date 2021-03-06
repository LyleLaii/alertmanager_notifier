package utils

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config config file info
type Config struct {
	ConfigFileName string
}

// InitConfig initiate config instance
func InitConfig(cfg string) error {
	c := Config{
		ConfigFileName: cfg,
	}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()
	return nil
}

func (c *Config) initConfig() error {
	viper.New()
	if c.ConfigFileName != "" {
		viper.SetConfigFile(c.ConfigFileName) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("settings")
	}
	viper.SetConfigType("yaml")         // 设置配置文件格式为YAML
	viper.AutomaticEnv()                // 读取匹配的环境变量
	viper.SetEnvPrefix("alertmanager_notifier") // 读取环境变量的前缀为alertnotifier
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}
	// a := viper.AllSettings()
	// for key := range a {
	// 	fmt.Println(key, a[key])
	// }
	return nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s", e.Name)
	})
}
