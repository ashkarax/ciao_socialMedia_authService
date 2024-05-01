package config_authSvc

import (
	"github.com/spf13/viper"
)

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
}

type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
}

type Token struct {
	AdminSecurityKey      string `mapstructure:"ADMIN_TOKENKEY"`
	RestaurantSecurityKey string `mapstructure:"RESTAURANT_TOKENKEY"`
	UserSecurityKey       string `mapstructure:"USER_TOKENKEY"`
	TempVerificationKey   string `mapstructure:"TEMPERVERY_TOKENKEY"`
}

type Smtp struct {
	SmtpSender   string `mapstructure:"SMTP_SENDER"`
	SmtpPassword string `mapstructure:"SMTP_APPKEY"`
	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtpPort     string `mapstructure:"SMTP_PORT"`
}

type Config struct {
	PortMngr PortManager
	DB       DataBase
	Token    Token
	Smtp     Smtp
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var db DataBase
	var token Token
	var smtp Smtp

	viper.AddConfigPath("./pkg/infrastructure/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&portmngr)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&token)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&smtp)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmngr, DB: db, Token: token, Smtp: smtp}
	return &config, nil

}