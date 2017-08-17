#!/bin/bash

openssl req \
  -x509 -newkey rsa:2048 -nodes -days 365 \
  -keyout tls.key -out tls.crt -subj '/CN=localhost'

kubectl create secret tls tls-secret --cert=tls.crt --key=tls.key

kubectl run ingress-default-backend \
  --image=gcr.io/google_containers/defaultbackend:1.0 \
  --port=8080 \
  --limits=cpu=10m,memory=20Mi \
  --expose

kubectl create configmap haproxy-ingress

# wget https://github.com/kubernetes/ingress/raw/master/examples/deployment/haproxy/haproxy-ingress.yaml
kubectl create -f haproxy-ingress.yaml

kubectl expose deploy/haproxy-ingress --type=NodePort

kubectl create -f mysql.yaml