package config

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Conf        Config
	Env         string
	RedisClient *redis.Client

	EnvironmentLocal = "LOCAL"
	EnvironmentDev   = "DEV"
	EnvironmentUAT   = "UAT"
	EnvironmentProd  = "PROD"

	searchPath = []string{
		"/etc/gin-boilerplate",
		"$HOME/.gin-boilerplate",
		".",
	}
	configDefaults = map[string]interface{}{
		"port":                  8080,
		// INFO, bukan DEBUG. Bawaan berlaku justru pada pemasangan yang confignya
		// paling ringkas — biasanya produksi — sedangkan pengembangan lokal
		// menyetelnya eksplisit. DEBUG di sini membuat jejak SQL GORM mencatat
		// setiap kueri lengkap dengan nilai yang diikat. Lihat database.levelJejakSQL.
		"logLevel":              "INFO",
		"logFormat":             "text",
		"secretKey":             "supersecret",
		"accessTokenExpiry":     24,
		"accessTokenExpiryUnit": "hour",
		"httpTimeout":           "30s",
		"rate":                  100,
		"interval":              "1m",
		"redis.enableRedis":     false,
		"postgres.host":         "localhost",
		"postgres.port":         5432,
		"postgres.user":         "postgres",
		"postgres.password":     "password",
		"postgres.dbname":       "gin_boilerplate",
		"postgres.sslmode":      "disable",
		"postgres.timezone":     "UTC",
		// Security configurations
		"security.enableSanitization":          true,
		"security.enableXSSDetection":          true,
		"security.enableSQLInjectionDetection": true,
		"security.maxStringLength":             1000,
		"security.strictMode":                  true,
		"security.maxBodySize":                 10485760, // 10MB
		"security.maxFormFields":               100,
		"security.maxFormMemory":               33554432, // 32MB
		"security.maxFormFiles":                10,
		"security.maxHeaderSize":               1048576, // 1MB
		"security.maxURLLength":                2048,
		"security.maxQueryParams":              50,
		"security.maxFileSize":                 52428800,  // 50MB
		"security.maxResponseSize":             104857600, // 100MB
	}
	configName = map[string]string{
		"local": "config.local",
		"dev":   "config.dev",
		"uat":   "config.uat",
		"prod":  "config.prod",
		"test":  "config.test",
	}
)

type Config struct {
	Env                   string         `mapstructure:"env"`
	Port                  int            `mapstructure:"port"`
	LogLevel              string         `mapstructure:"logLevel"`
	LogMode               bool           `mapstructure:"logMode"`
	LogFormat             string         `mapstructure:"logFormat"`
	Postgres              PostgresConfig `mapstructure:"postgres"`
	Redis                 RedisConfig    `mapstructure:"redis"`
	Rate                  int64          `mapstructure:"rate"`
	Interval              string         `mapstructure:"interval"`
	SecretKey             string         `mapstructure:"secretKey"`
	AccessTokenExpiry     int            `mapstructure:"accessTokenExpiry"`
	AccessTokenExpiryUnit string         `mapstructure:"accessTokenExpiryUnit"`
	HTTPTimeout           time.Duration  `mapstructure:"httpTimeout"`
	Security              SecurityConfig `mapstructure:"security"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type RedisConfig struct {
	EnableRedis bool   `mapstructure:"enableRedis"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db"`
}

type SecurityConfig struct {
	// Input Sanitization
	EnableSanitization          bool `mapstructure:"enableSanitization"`
	EnableXSSDetection          bool `mapstructure:"enableXSSDetection"`
	EnableSQLInjectionDetection bool `mapstructure:"enableSQLInjectionDetection"`
	MaxStringLength             int  `mapstructure:"maxStringLength"`
	StrictMode                  bool `mapstructure:"strictMode"`

	// Request Size Limiting
	MaxBodySize     int64 `mapstructure:"maxBodySize"`
	MaxFormFields   int   `mapstructure:"maxFormFields"`
	MaxFormMemory   int64 `mapstructure:"maxFormMemory"`
	MaxFormFiles    int   `mapstructure:"maxFormFiles"`
	MaxHeaderSize   int64 `mapstructure:"maxHeaderSize"`
	MaxURLLength    int   `mapstructure:"maxURLLength"`
	MaxQueryParams  int   `mapstructure:"maxQueryParams"`
	MaxFileSize     int64 `mapstructure:"maxFileSize"`
	MaxResponseSize int64 `mapstructure:"maxResponseSize"`
}

func initialiseFileAndEnv(v *viper.Viper, env string) error {
	v.SetConfigName(configName[env])
	for _, path := range searchPath {
		v.AddConfigPath(path)
	}
	v.SetEnvPrefix("GIN_BOILERPLATE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	return v.ReadInConfig()
}

func initialiseDefaults(v *viper.Viper) {
	for key, value := range configDefaults {
		v.SetDefault(key, value)
	}
}

func Initialize() {
	v := viper.New()
	initialiseDefaults(v)

	if err := initialiseFileAndEnv(v, Env); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Warning("No config file found. Using environment variables and defaults")
		} else {
			log.Warningf("Error reading config file: %v", err)
		}
	}

	err := v.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("Error unmarshalling configuration: %s", err.Error())
	}

	log.Infof("Configuration loaded for environment: %s", Conf.Env)
}
