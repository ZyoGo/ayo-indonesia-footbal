package config

import (
	"crypto/rsa"
	"fmt"
	"os"
	"sync"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var (
	appConfig *AppConfig
	lock      = &sync.Mutex{}
	jwtKeys   *JWTKeys
)

type AppConfig struct {
	App struct {
		Port           uint16   `toml:"port"`
		AllowedOrigins []string `toml:"allowedOrigins"`
	} `toml:"app"`
	Database struct {
		Name     string `toml:"name"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Port     uint16 `toml:"port"`
		Address  string `toml:"address"`
		Driver   string `toml:"driver"`
	} `toml:"database"`
	JWT struct {
		PrivateKeyPath string `toml:"private_key_path"`
		PublicKeyPath  string `toml:"public_key_path"`
		Issuer         string `toml:"issuer"`
		Subject        string `toml:"subject"`
	} `toml:"jwt"`
}

type JWTKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		var err error
		appConfig, err = loadConfig()
		if err != nil {
			logger.Get().Error(fmt.Sprintf("Failed to load config: %v", err))
			os.Exit(1)
		}
	}

	return appConfig
}

func GetJWTKeys() *JWTKeys {
	lock.Lock()
	defer lock.Unlock()

	if jwtKeys == nil {
		var err error
		jwtKeys, err = loadJWTKeys()
		if err != nil {
			logger.Get().Error(fmt.Sprintf("Failed to load JWT keys: %v", err))
			os.Exit(1)
		}
	}

	return jwtKeys
}

func loadConfig() (*AppConfig, error) {
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")
	viper.SetConfigName("app")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
		return nil, err
	}

	var config AppConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	logger.Get().Info("Configuration successfully loaded")

	return &config, nil
}

func loadJWTKeys() (*JWTKeys, error) {
	cfg := GetConfig()

	privateKeyPEM, err := os.ReadFile(cfg.JWT.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	publicKeyPEM, err := os.ReadFile(cfg.JWT.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	logger.Get().Info("JWT keys loaded successfully")

	return &JWTKeys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
