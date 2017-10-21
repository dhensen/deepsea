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

CERT_NAME="nginx-tls"
KEY_FILE="nginx-tls.key"
CERT_FILE="nginx-tls.crt"
HOST="localhost"

# Generate keys
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
# Create a TLS secret
kubectl create secret tls ${CERT_NAME} --key ${KEY_FILE} --cert ${CERT_FILE}
# Remove locally generated key and certificate
rm ${KEY_FILE} ${CERT_FILE}

# Common
kubectl run echoheaders --image=gcr.io/google_containers/echoserver:1.8 --replicas=1 --port=8080
kubectl expose deployment echoheaders --port=80 --target-port=8080 --name=echoheaders-x
kubectl expose deployment echoheaders --port=80 --target-port=8080 --name=echoheaders-y
kubectl create -f common/echo-headers-ingress.yaml
kubectl create -f nginx/default-backend.yaml
# why does this need to be exposed? also the .yaml already contains the service definition
kubectl expose deployment default-http-backend --port=80 --target-port=8080 --name=default-http-backend --namespace=kube-system


# Create the nginx ingress controller
kubectl apply -f nginx/nginx-ingress-controller.yaml
# Expose the nginx ingress controller
kubectl expose deployment nginx-ingress-controller --type=NodePort --namespace=kube-system