#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR/..

VERSION=$(date +%s)
REGISTRY_HOST=$(minikube ip)
REGISTRY_PORT=5000


docker build -t deepsea:$VERSION .
docker tag deepsea:$VERSION $REGISTRY_HOST:$REGISTRY_PORT/deepsea:$VERSION

kubectl port-forward $(kubectl get pods | grep kube-registry | awk '{print $1;}') 5000:5000 &
FORWARD_PID=$!

# Ditry: wait for forwarder to setup connection
sleep 2

docker push $REGISTRY_HOST:$REGISTRY_PORT/deepsea:$VERSION

kill ${FORWARD_PID} > /dev/null 2>&1

echo "==> Update deepsea version in deployment"
kubectl set image \
    -f k8s-provision/common/deepsea.yaml \
    deepsea=localhost:$REGISTRY_PORT/deepsea:$VERSION \
    --local \
    -o yaml | kubectl apply -f -
echo "Deepsea deployment done"

# Take note that the actual yaml is not updated, only the version is replaced
# and the resulting yaml on stdout is piped into kubectl apply