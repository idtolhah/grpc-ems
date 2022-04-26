#!/bin/bash
kubectl apply -f k8s/

kubectl apply -f bff/bff-deployment.yaml
kubectl apply -f bff/bff-service.yaml

kubectl apply -f comment-cmd/comment-cmd-deployment.yaml
kubectl apply -f comment-cmd/comment-cmd-service.yaml

kubectl apply -f master-query/master-query-deployment.yaml
kubectl apply -f master-query/master-query-service.yaml

kubectl apply -f packing-cmd/packing-cmd-deployment.yaml
kubectl apply -f packing-cmd/packing-cmd-service.yaml

kubectl apply -f packing-query/packing-query-deployment.yaml
kubectl apply -f packing-query/packing-query-service.yaml

kubectl apply -f user-query/user-query-deployment.yaml
kubectl apply -f user-query/user-query-service.yaml