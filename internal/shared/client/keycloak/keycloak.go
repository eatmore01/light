package keycloak

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/eatmore01/light/internal/config"
)

type KeycloakClient struct {
	KeycloakClient *gocloak.GoCloak
	Config         *config.Config
}

func NewKeycloakCLient(config *config.Config) *KeycloakClient {
	keycloakClient := gocloak.NewClient(config.KeycloakHost)

	return &KeycloakClient{
		KeycloakClient: keycloakClient,
		Config:         config,
	}
}
