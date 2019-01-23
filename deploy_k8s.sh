#!/bin/bash

set -ex

echo "sed on k8s/${ENV}/deployment.yaml"
sed -i 's/%BUILD_NUMBER%/'"$GAMECHAIN_BUILD_NUMBER"'/g' k8s/${ENV}/deployment.yaml

echo "kube apply deployment"
kubectl apply -f k8s/${ENV}/deployment.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
echo "kube apply service"
kubectl apply -f k8s/${ENV}/service.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
echo "kube apply ingress"
kubectl apply -f k8s/${ENV}/ingress.yaml  --kubeconfig=/var/lib/jenkins/${ENV}_kube_config.yaml
