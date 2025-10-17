#!/bin/bash

set -e

echo "Deploying Task Manager to AKS..."

# Apply PostgreSQL configuration
echo "Applying PostgreSQL configuration..."
kubectl apply -f ../../k8s/config/

# Deploy services in order (database-dependent first)
echo "Deploying user-service..."
kubectl apply -f ../../k8s/manifests/user-service/

# echo "Deploying auth-service..."
# kubectl apply -f ../../k8s/manifests/auth-service/

# echo "Deploying task-service..."
# kubectl apply -f ../../k8s/manifests/task-service/

# echo "Deploying notification-service..."
# kubectl apply -f ../../k8s/manifests/notification-service/

# echo "Deploying api-gateway..."
# kubectl apply -f ../../k8s/manifests/api-gateway/

# echo "Deploying frontend..."
# kubectl apply -f ../../k8s/manifests/frontend/

# echo "Deploying ingress..."
# kubectl apply -f ../../k8s/manifests/ingress/

# echo "Waiting for deployments to be ready..."
# kubectl wait --for=condition=available --timeout=300s deployment/user-service
# kubectl wait --for=condition=available --timeout=300s deployment/auth-service
# kubectl wait --for=condition=available --timeout=300s deployment/task-service
# kubectl wait --for=condition=available --timeout=300s deployment/notification-service
# kubectl wait --for=condition=available --timeout=300s deployment/api-gateway
# kubectl wait --for=condition=available --timeout=300s deployment/frontend

echo "Deployment completed!"
echo ""
echo "Checking pod status:"
kubectl get pods

echo ""
echo "Services:"
kubectl get services

echo ""
echo "To check logs: kubectl logs -f deployment/<service-name>"