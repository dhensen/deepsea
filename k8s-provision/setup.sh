#!/bin/bash

minikube ip
if [ $? -ne 0 ]; then
  minikube start --insecure-registry "127.0.0.1:5000"
fi

kubectl create -f kube-registry.yaml

# See https://github.com/kubernetes/ingress/tree/master/examples/deployment/haproxy for more information

# if folder secret does not exist create TLS secret
if [ ! -d secret ]; then
  mkdir -p secret
  openssl req \
    -x509 -newkey rsa:2048 -nodes -days 365 \
    -keyout tls.key -out tls.crt -subj '/CN=localhost'
fi

kubectl create secret tls tls-secret --cert=secret/tls.crt --key=secret/tls.key

kubectl run ingress-default-backend \
  --image=gcr.io/google_containers/defaultbackend:1.0 \
  --port=8080 \
  --limits=cpu=10m,memory=20Mi \
  --expose

kubectl create configmap haproxy-ingress

# wget https://github.com/kubernetes/ingress/raw/master/examples/deployment/haproxy/haproxy-ingress.yaml
# modification: added hostPort for 80,443 and 1936
kubectl create -f haproxy-ingress.yaml

kubectl expose deploy/haproxy-ingress --type=NodePort

kubectl create -f mysql.yaml

# For the PV we need to create the /data/volumes folder and own it
# Everything inside /data is persisted across minikube reboot.
# Beware: ofcourse this does not scale and PV uses hostPath is only for development purposes
minikube ssh "sudo mkdir -p /data/volumes && sudo chown -R docker:docker /data/volumes"

docker build -t deepsea:v1 $PWD/../
docker tag deepsea:v1 192.168.99.100:5000/deepsea:v1

status=$(kubectl get pods | grep registry | awk '{print $3;}')
while [ "$status" != 'Running' ];
do
  echo "waiting for registry to be ready before pushing image"
  sleep 2
  status=$(kubectl get pods | grep registry | awk '{print $3;}')
done

docker push 192.168.99.100:5000/deepsea:v1

kubectl create -f $PWD/deepsea.yaml