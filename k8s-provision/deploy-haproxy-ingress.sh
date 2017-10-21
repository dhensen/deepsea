#!/usr/bin/bash

if [ -z "$KUBECONFIG" ]; then
  export KUBECONFIG=/home/dino/.secrets/clusters/deepsea/auth/kubeconfig
  echo "KUBECONFIG not set, defaulting to: ${KUBECONFIG}"
else
  echo "Using KUBECONFIG=${KUBECONFIG}"
fi

echo "kubectl's current context is: $(kubectl config current-context)"

# WARNING: setting this permissing clusterrolebinding is a security risk
# Read and understand this whole document and implement a better solution:
# https://kubernetes.io/docs/admin/authorization/rbac/#permissive-rbac-permissions
kubectl create clusterrolebinding permissive-binding \
  --clusterrole=cluster-admin \
  --user=admin \
  --user=kubelet \
  --group=system:serviceaccounts

CERT_NAME="haproxy-tls"
KEY_FILE="haproxy-tls.key"
CERT_FILE="haproxy-tls.crt"
HOST="localhost"

# Generate keys
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
# Create a TLS secret
kubectl create secret tls ${CERT_NAME} --key ${KEY_FILE} --cert ${CERT_FILE}
# Remove locally generated key and certificate
rm ${KEY_FILE} ${CERT_FILE}

# Common
echo "==> Create haproxy-ingress namespace"
kubectl create namespace haproxy-ingress

echo "==> Create ingress-default-backend in haproxy-ingress namespace"
kubectl run ingress-default-backend \
  --image=gcr.io/google_containers/defaultbackend:1.0 \
  --port=8080 \
  --limits=cpu=10m,memory=20Mi \
  --namespace=haproxy-ingress \
  --expose

echo "==> Create haproxy-ingress configmap"
kubectl create configmap haproxy-ingress --namespace=haproxy-ingress

echo "==> Create app http-svc in default namespace"
kubectl run http-svc \
  --image=gcr.io/google_containers/echoserver:1.3 \
  --port=8080 \
  --replicas=1 \
  --expose

echo "==> Create foo.bar ingress"
kubectl create -f - <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: app
spec:
  rules:
  - host: foo.bar
    http:
      paths:
      - path: /
        backend:
          serviceName: http-svc
          servicePort: 8080
EOF

# Create the haproxy ingress controller
echo "==> Create haproxy-ingress controller"
kubectl create -f haproxy/haproxy-ingress.yaml
# Expose the haproxy ingress controller
echo "==> Expose haproxy-ingress"
kubectl expose deploy/haproxy-ingress --type=NodePort --namespace=haproxy-ingress
# kubectl get svc/haproxy-ingress -oyaml

