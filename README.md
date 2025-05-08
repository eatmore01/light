![image](./image.png)

# Light

This is project needed for deploting to cluster, after login via web interface with keycloak creds and generate personal kubeconfig for current cluster 

- [Prerequists](#prerequists)
- [Image](#image)
- [Deploy](#deploy)
- [Deploy via helmchart](#deploy-via-helmchart)


## Prerequists 

- Modifying kube-apiserver for login with keycloak OIDC


## Image

- `yazhivotnoe/light:0.0.1`

## Deploy

- Create `config.yml` on `config/config.yml` path witch can looks like

```yaml
# base server config
host: "0.0.0.0"
port: "9999"

# Need to change for every cluster
## deployed cluster name
clusterName: "commanda-ahuenaya.dev"
## deployed apiAddress
## template variables https://<ip-address-master-node>:6433 OR ur LB address with dicrover master node
## or other address witch poin in worked alredy exist kubeconfig
cluesterApiAddress: "https://cluster.ru:6443"


# Dont need to change
## default path to CA cluster cert if it placed in pods with serviceaccount 
clusterCAPath: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

## keycloak vars
keycloakHost: "https://keycloak.com"
keycloakRealm: "master"

idpIssuerUrl: "https://keycloak.com/realms/master"
clientID: "kubernetes"
clientSecret: ""

usernameClaim: "preferred_username"

# App vars
TemplatesDir: "templates"
jwtsecret: "SECRET"
cookieSecure: false
```

### Deploy via helmchart

https://github.com/yazhivotnoe/tf-mod-kube-light

