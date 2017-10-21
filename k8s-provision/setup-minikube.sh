#!/bin/bash

# Start minikube if it is not already running
minikube ip 2>/dev/null
if [ $? -ne 0 ]; then
  echo "==> Starting Minikube"
  minikube start --insecure-registry "127.0.0.1:5000"
fi

# Explicitly use minikube context just for safety:
# in case minikube was already running and another context was switched to
kubectl config use-context minikube
export KUBECONFIG=~/.kube/config

# Deploy Ingress Controller
# ./deploy-haproxy-ingress.sh
./deploy-nginx-ingress.sh

# Deploy Docker Registry
echo "==> Create kube-registry"
kubectl create -f common/kube-registry.yaml

# Deploy MySQL
echo "==> Create mysql"
kubectl create -f common/mysql.yaml

# For the PV we need to create the /data/volumes folder and own it
# Everything inside /data is persisted across minikube reboot.
# Beware: ofcourse this does not scale and PV uses hostPath is only for development purposes
echo "==> Create PV folders on minikube node"
minikube ssh "sudo mkdir -p /data/volumes && sudo chown -R docker:docker /data/volumes"

# Wait before deployed registry is ready, the perform deepsea deploy

echo "==> Waiting untill registry is up"
status=$(kubectl get pods | grep registry | awk '{print $3;}')
while [ "$status" != 'Running' ];
do
  echo "waiting for registry to be ready before pushing image"
  sleep 2
  status=$(kubectl get pods | grep registry | awk '{print $3;}')
done

../bin/deepsea-deploy.sh
