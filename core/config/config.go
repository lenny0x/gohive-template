package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Environment 环境类型
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvTest        Environment = "test"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
)

// BaseConfig 基础配置
type BaseConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Log      LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	GroupID string   `mapstructure:"group_id"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LoadConfig 加载配置
// configPath: 配置文件路径 (如 ./config/config.toml)
// 支持环境变量覆盖，格式: APP_DATABASE_HOST
func LoadConfig[T any](configPath string) (*T, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)

	// 设置环境变量前缀和绑定
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值
	setDefaults(v)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// 加载环境特定配置
	if err := loadEnvConfig(v, configPath); err != nil {
		// 环境配置不存在不报错，只是用基础配置
		fmt.Printf("No environment config found: %v\n", err)
	}

	// 解析配置到结构体
	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// loadEnvConfig 加载环境特定配置
func loadEnvConfig(v *viper.Viper, basePath string) error {
	// 获取当前环境
	env := v.GetString("app.env")
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = string(EnvDevelopment)
	}

	// 构建环境配置文件路径
	// config.toml -> config.development.toml
	ext := getFileExt(basePath)
	envPath := strings.TrimSuffix(basePath, ext) + "." + env + ext

	// 检查文件是否存在
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return fmt.Errorf("env config not found: %s", envPath)
	}

	// 合并环境配置
	v.SetConfigFile(envPath)
	if err := v.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to merge env config: %w", err)
	}

	return nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// App
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)

	// Database
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.max_idle_conns", 10)

	// Redis
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)

	// Kafka
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})

	// Log
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")
}

// getFileExt 获取文件扩展名
func getFileExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// GetEnv 获取当前环境
func GetEnv() Environment {
	env := os.Getenv("APP_ENV")
	switch env {
	case "test":
		return EnvTest
	case "staging":
		return EnvStaging
	case "production":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}

// IsDevelopment 是否开发环境
func IsDevelopment() bool {
	return GetEnv() == EnvDevelopment
}

// IsProduction 是否生产环境
func IsProduction() bool {
	return GetEnv() == EnvProduction
}
