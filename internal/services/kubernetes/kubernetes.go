package kubernetes

import (
	"encoding/base64"
	"io/ioutil"
	"log"

	// for out-of-cluster testing, include:
	// "k8s.io/client-go/tools/clientcmd"

	"github.com/eatmore01/light/internal/config"
	"github.com/eatmore01/light/internal/controllers/auth"
	kubernetesClient "github.com/eatmore01/light/internal/shared/client/kubernetes" // for in cluster access
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
)

type KubeService struct {
	k8sClient *kubernetes.Clientset
	k8sCfg    *rest.Config
	AppCfg    *config.Config
}

type KubeConfigInfo struct {
	ClusterName    string
	APIServerUrl   string
	ClusterCAData  string
	KubeConfigUser string
	CurrentCtx     string
	ClientID       string
	ClientSecret   string
	IDToken        string
	RefreshToken   string
	IssuerUrl      string
}

func NewKubeService() *KubeService {
	client, config := kubernetesClient.NewClient()

	KubeService := &KubeService{
		k8sClient: client,
		k8sCfg:    config,
	}

	return KubeService
}

func (ks *KubeService) GenerateInfo(cfg *config.Config, claims *auth.CustomClaims) KubeConfigInfo {
	caCert, err := ioutil.ReadFile(cfg.CLusterCAPath)
	if err != nil {
		log.Fatalf("Error reading CA cert: %v", err)
	}

	CAData := base64.StdEncoding.EncodeToString(caCert)

	return KubeConfigInfo{
		ClusterName:    cfg.ClusterName,
		APIServerUrl:   cfg.CluesterApiAddress,
		ClusterCAData:  CAData,
		KubeConfigUser: claims.Username,
		CurrentCtx:     cfg.ClusterName,
		IDToken:        claims.IDToken,
		RefreshToken:   claims.RefreshToken,
		ClientID:       cfg.ClientID,
		ClientSecret:   cfg.ClientSecret,
		IssuerUrl:      cfg.IssuerUrl,
	}
}

func (ks *KubeService) GenerateKubeConfig(info KubeConfigInfo) *clientcmdapi.Config {

	caData, err := base64.StdEncoding.DecodeString(info.ClusterCAData)
	if err != nil {
		log.Fatalf("Error decoding CA data: %v", err)
	}

	kcfg := &clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: []clientcmdapi.NamedCluster{
			{
				Name: info.ClusterName,
				Cluster: clientcmdapi.Cluster{
					Server:                   info.APIServerUrl,
					CertificateAuthorityData: caData,
				},
			},
		},
		Contexts: []clientcmdapi.NamedContext{
			{
				Name: info.CurrentCtx,
				Context: clientcmdapi.Context{
					Cluster:  info.ClusterName,
					AuthInfo: info.KubeConfigUser,
				},
			},
		},
		CurrentContext: info.CurrentCtx,
		Preferences:    clientcmdapi.Preferences{},
		AuthInfos: []clientcmdapi.NamedAuthInfo{
			{
				Name: info.KubeConfigUser,
				AuthInfo: clientcmdapi.AuthInfo{
					AuthProvider: &clientcmdapi.AuthProviderConfig{
						Name: "oidc",
						Config: map[string]string{
							"client-id":      info.ClientID,
							"client-secret":  info.ClientSecret,
							"id-token":       info.IDToken,
							"refresh-token":  info.RefreshToken,
							"idp-issuer-url": info.IssuerUrl,
						},
					},
				},
			},
		},
	}

	return kcfg
}
