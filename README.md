# k8s_training

k8s getting started. just an example

## Directory

```
k8s-demo/
├── app/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── mysql/
│   ├── statefulset.yaml
│   ├── service.yaml
│   └── kustomization.yaml
├── traefik/
│   └── values.yaml
├── deploy/
│   ├── base/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── gateway.yaml
│   │   ├── httproute.yaml
│   │   └── kustomization.yaml
│   └── overlays/
│       └── local/
│           ├── kustomization.yaml
│           └── patch-image.yaml
└── README.md
```

## Docker

docker build -t demo-app:local ./app


## Traefik

`cd traefik`

⚠️ DEPRECATION WARNING: Gateway API CRDs will no longer be shipped with this chart in a future major version.
You will need to install them yourself before deploying Traefik v3.7:

`kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.5.1/standard-install.yaml`

- install
`helm install traefik traefik/traefik -n traefik -f values.yaml`

- update
`helm upgrade traefik traefik/traefik -n traefik -f values.yaml`

- uninstall
`helm uninstall traefik -n traefik`

- check default in traefik installer
`helm show values traefik/traefik > values-default.yaml`

## K8S

kubectl apply -k deploy/overlays/local
kubectl delete -k deploy/overlays/local
kubectl rollout restart deployment/demo-app

## Gateway

kubectl get gatewayclass -A
kubectl get gateway -A

## Test

### Local HTTPS

#### install tool

```bash
choco install mkcert
```

#### create local ca

```bash
mkcert -install
```

#### check root ca


```bash
mkcert -CAROOT

```

#### create certificate

```bash
mkcert localhost demo.example.com 127.0.0.1 ::1
```

export: `localhost+3.pem` and `localhost+3-key.pem`

#### create k8s secret

```bash
kubectl create secret tls local-cert-tls \
    --cert=localhost+3.pem \
    --key=localhost+3-key.pem \
    -n traefik
```

### Run

check server `curl --ssl-revoke-best-effort https://localhost/app/healthz`
test queue work `curl --ssl-revoke-best-effort https://localhost/app/work`
test entity check `curl --ssl-revoke-best-effort https://localhost/app/items`
test entity write `curl -X POST --ssl-revoke-best-effort https://localhost/app/items`