#!/bin/bash

helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm update

kubectl create namespace nats
helm install nats nats/nats --namespace nats

