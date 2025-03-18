package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	App struct {
		Name             string `yaml:"name" env:"NAME" envDefault:"Blog"`
		Version          string `yaml:"version" env:"VERSION" envDefault:"Latest"`
		LogLevel         string `yaml:"log_level" env:"LOG_LEVEL" envDefault:"debug"`
		Env              string `yaml:"env" env:"ENV" envDefault:"development"`
		Locale           string `yaml:"locale" env:"LOCALE" envDefault:"en"`
		FallbackLocale   string `yaml:"fallback_locale" env:"FALLBACK_LOCALE" envDefault:"en"`
		TranslationsPath string `yaml:"translations_path" env:"TRANSLATIONS_PATH" env-default:"assets/locales"`
	} `yaml:"app" env-prefix:"APP_"`
	Router struct {
		Address string `yaml:"address" env:"ADDRESS"`
	} `yaml:"router" env-prefix:"ROUTER_"`
	DB struct {
		Mongo struct {
			HostName   string `yaml:"host_name" env:"HOST_NAME"`
			Port       string `yaml:"port" env:"PORT"`
			Username   string `yaml:"username" env:"USERNAME"`
			Password   string `yaml:"password" env:"PASSWORD"`
			Database   string `yaml:"database" env:"DATABASE"`
			AuthSource string `yaml:"auth_source" env:"AUTH_SOURCE"`
		} `yaml:"mongo" env-prefix:"MONGO_"`
	} `yaml:"db" env-prefix:"DB_"`
	Auth struct {
		JWT struct {
			SecretKey   string `yaml:"secret_key" env:"SECRET_KEY"`
			Issuer      string `yaml:"issuer" env:"ISSUER" env-default:"Blog"`
			ExpiredTime int    `yaml:"expired_time" env:"EXPIRED_TIME" envDefault:"24"` //Hour
		} `yaml:"jwt" env-prefix:"JWT_"`
	} `yaml:"auth" env-prefix:"AUTH_"`
}

func Load(filePath string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		log.Error(err)
		if err = cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
