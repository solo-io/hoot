#!/bin/bash

set -e

kubectl create ns spire

kubectl apply -f csidriver.yaml
kubectl apply -f crds.yaml
kubectl apply -f spire-controller-manager-config.yaml
kubectl apply -f spire-controller-manager-webhook.yaml