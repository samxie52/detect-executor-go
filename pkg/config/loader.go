package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Loader 配置加载器 loader of config loader
type Loader struct {
	v      *viper.Viper
	config *Config
}

func NewLoader() *Loader {
	v := viper.New()

	// setting file type of config file
	v.SetConfigType("yaml")

	// setting search path of config file
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/detect-executor-go")
	v.AddConfigPath("$HOME/detect-executor-go")

	//set prefix of environment variable
	v.SetEnvPrefix("DETECT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return &Loader{
		v:      v,
		config: &Config{},
	}
}

// Load 加载配置 load config

func (l *Loader) Load(configFile string) (*Config, error) {
	// set default config values
	l.setDefaults()

	//if config file is specified, use it
	if configFile != "" {
		l.v.SetConfigFile(configFile)
	} else {
		//according to the environment variable to load config file
		env := os.Getenv("GO_ENV")
		if env == "" {
			env = "dev"
		}
		l.v.SetConfigName(fmt.Sprintf("config.%s", env))
	}

	//read config file
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found")
		} else {
			return nil, fmt.Errorf("read config file failed: %v", err)
		}
	}

	//unmarshal config file
	if err := l.v.Unmarshal(l.config); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %v", err)
	}

	return l.config, nil
}

func (l *Loader) Watch(callback func(config *Config)) error {
	l.v.WatchConfig()
	l.v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)

		//reload config
		if err := l.v.Unmarshal(l.config); err != nil {
			fmt.Printf("unmarshal config failed: %v\n", err)
			return
		}

		//call callback
		callback(l.config)
	})

	return nil
}

func (l *Loader) setDefaults() {
	l.v.SetDefault("server.host", "0.0.0.0")
	l.v.SetDefault("server.port", 8080)
	l.v.SetDefault("server.read_timeout", 30*time.Second)
	l.v.SetDefault("server.write_timeout", 30*time.Second)
	l.v.SetDefault("server.idle_timeout", 5*time.Minute)

	l.v.SetDefault("database.host", "localhost")
	l.v.SetDefault("database.port", 3306)
	l.v.SetDefault("database.username", "root")
	l.v.SetDefault("database.password", "")
	l.v.SetDefault("database.database", "detect_executor")
	l.v.SetDefault("database.charset", "utf8mb4")
	l.v.SetDefault("database.max_idle_conns", 10)
	l.v.SetDefault("database.max_open_conns", 100)
	l.v.SetDefault("database.conn_max_lifetime", 1*time.Hour)

	l.v.SetDefault("redis.host", "localhost")
	l.v.SetDefault("redis.port", 6379)
	l.v.SetDefault("redis.password", "")
	l.v.SetDefault("redis.db", 0)
	l.v.SetDefault("redis.pool_size", 10)
	l.v.SetDefault("redis.min_idle_conns", 10)
	l.v.SetDefault("redis.dial_timeout", 5*time.Second)
	l.v.SetDefault("redis.read_timeout", 5*time.Second)
	l.v.SetDefault("redis.write_timeout", 5*time.Second)

	l.v.SetDefault("detect.max_workers", 10)
	l.v.SetDefault("detect.batch_size", 10)
	l.v.SetDefault("detect.queue_size", 100)
	l.v.SetDefault("detect.timeout", 10*time.Second)
	l.v.SetDefault("detect.retry_count", 3)
	l.v.SetDefault("detect.retry_delay", 1*time.Second)
	l.v.SetDefault("detect.cpp_engine_url", "http://localhost:8081")
	l.v.SetDefault("detect.video_stream_url", "http://localhost:8082")

	l.v.SetDefault("log.level", "info")
	l.v.SetDefault("log.format", "json")
	l.v.SetDefault("log.output", "stdout")
	l.v.SetDefault("log.filename", "logs/detect-executor-go.log")
	l.v.SetDefault("log.max_size", 100)
	l.v.SetDefault("log.max_age", 28)
	l.v.SetDefault("log.max_backups", 3)
	l.v.SetDefault("log.compress", true)

	l.v.SetDefault("metrics.enabled", false)
	l.v.SetDefault("metrics.path", "/metrics")
	l.v.SetDefault("metrics.port", 9090)
}
