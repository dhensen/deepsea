#!/usr/bin/bash

export KUBECONFIG=/home/dino/.secrets/clusters/deepsea/auth/kubeconfig

openssl req \
  -x509 -newkey rsa:2048 -nodes -days 365 \
  -keyout tls.key -out tls.crt -subj '/CN=localhost'

kubectl create secret tls tls-secret --cert=tls.crt --key=tls.key

rm -v tls.crt tls.key

kubectl run http-svc \
  --image=gcr.io/google_containers/echoserver:1.3 \
  --port=8080 \
  --replicas=1 \
  --expose

kubectl run ingress-default-backend \
  --image=gcr.io/google_containers/defaultbackend:1.0 \
  --port=8080 \
  --limits=cpu=10m,memory=20Mi \
  --expose

kubectl create configmap haproxy-ingress

# WARNING: setting this permissing clusterrolebinding is a security risk
# Read and understand this whole document and implement a better solution:
# https://kubernetes.io/docs/admin/authorization/rbac/#permissive-rbac-permissions
kubectl create clusterrolebinding permissive-binding \
  --clusterrole=cluster-admin \
  --user=admin \
  --user=kubelet \
  --group=system:serviceaccounts

kubectl create -f haproxy-ingress.yaml

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

kubectl expose deploy/haproxy-ingress --type=NodePort
kubectl get svc/haproxy-ingress -oyaml

