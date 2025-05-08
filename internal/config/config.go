package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host string `yaml:"host" env-defaults:"0.0.0.0"`
	Port string `yaml:"port" env-defaults:"9999"`

	ClusterName   string `yaml:"clusterName" env-default:"cluster_name"`
	IssuerUrl     string `yaml:"idpIssuerUrl" env-default:"idpIssuerUrl"`
	ClientID      string `yaml:"clientID" env-default:"client_id"`
	ClientSecret  string `yaml:"clientSecret" env-default:"client_secret"`
	UsernameClaim string `yaml:"usernameClaim" env-default:"username_claim"`

	// template variables https://<ip-address-master-node>:16433
	// or other address witch poin in worked alredy exist kubeconfig
	CluesterApiAddress string `yaml:"cluesterApiAddress"`

	KeycloakHost  string `yaml:"keycloakHost"`
	KeycloakRealm string `yaml:"keycloakRealm"`

	CookieSecure bool `yaml:"cookieSecure" env-default:"false"`

	JWTSecret string `yaml:"jwtsecret" env-default:"SECRET"`

	// cluster ca path in pod via load kubernetes client via SA incluster method
	CLusterCAPath string `yaml:"clusterCAPath"`

	TemplatesDir string `yaml:"TemplatesDir" env-default:"templates"`
}

func MustLoad() *Config {
	var cfg Config
	path := "config/config.yml"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found on path:" + path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
