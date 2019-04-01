package skelconfig

import "errors"

var (
	ConfJWTSecretErr = errors.New("AppConfig.JWTSecret is not populated")

	//ConfDBHostErr     = errors.New("AppConfig.DBHost is not populated")
	//ConfDBUserErr     = errors.New("AppConfig.DBUser is not populated")
	//ConfDBPasswordErr = errors.New("AppConfig.DBPassword is not populated")
	//ConfDBNameErr     = errors.New("AppConfig.DBName is not populated")
	//ConfDBPortErr     = errors.New("AppConfig.DBPort is not populated")
	//ConfHashPepperErr = errors.New("AppConfig.HashPepper is not populated")
	//ConfConfigURLErr  = errors.New("AppConfig.ConfigUrl is not populated")
)
