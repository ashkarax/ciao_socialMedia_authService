package config_authSvc

import (
	"github.com/spf13/viper"
)

type PortManager struct {
	RunnerPort     string `mapstructure:"PORTNO"`
	PostNrelSvcUrl string `mapstructure:"POSTNREL_SVC_URL"`
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

type AWS struct {
	Region     string `mapstructure:"AWS_REGION"`
	AccessKey  string `mapstructure:"AWS_ACCESS_KEY_ID"`
	SecrectKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Endpoint   string `mapstructure:"AWS_ENDPOINT"`
}

type Config struct {
	PortMngr PortManager
	DB       DataBase
	Token    Token
	Smtp     Smtp
	AwsS3    AWS
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var db DataBase
	var token Token
	var smtp Smtp
	var awsS3 AWS

	viper.AddConfigPath("./")
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
	err = viper.Unmarshal(&awsS3)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmngr, DB: db, Token: token, Smtp: smtp, AwsS3: awsS3}
	return &config, nil

}
