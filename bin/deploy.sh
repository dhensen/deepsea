#!/bin/bash

VERSION=$(date +%s)
REGISTRY_HOST=$(minikube ip)
REGISTRY_PORT=5000


docker build -t deepsea:$VERSION .
docker tag deepsea:$VERSION $REGISTRY_HOST:$REGISTRY_PORT/deepsea:$VERSION

kubectl port-forward $(kubectl get pods | grep kube-registry | awk '{print $1;}') 5000:5000 &

sleep 2

docker push $REGISTRY_HOST:$REGISTRY_PORT/deepsea:$VERSION

# kubectl create -f ./deepsea.yaml
kubectl set image deployment/deepsea deepsea=localhost:$REGISTRY_PORT/deepsea:$VERSION
