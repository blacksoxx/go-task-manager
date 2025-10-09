#!/bin/bash

set -e

# Configuration
ACR_NAME="acrtaskmanagerdev"
TAG="latest"
SERVICES=("api-gateway" "user-service" "auth-service" "task-service" "notification-service" "frontend")

echo "Starting Docker build and push to ACR..."

# Login to ACR
echo "Logging into ACR..."
az acr login --name $ACR_NAME

# Build and push each service
for SERVICE in "${SERVICES[@]}"; do
    echo "Building $SERVICE..."
    
    # Build the image
    docker build -t $SERVICE:$TAG ./$SERVICE
    
    # Tag for ACR
    docker tag $SERVICE:$TAG $ACR_NAME.azurecr.io/$SERVICE:$TAG
    
    # Push to ACR
    echo "⬆Pushing $SERVICE to ACR..."
    docker push $ACR_NAME.azurecr.io/$SERVICE:$TAG
    
    echo " $SERVICE pushed successfully"
    echo ""
done

echo "All images built and pushed to ACR!"
echo ""
echo "Images in ACR:"
az acr repository list --name $ACR_NAME --output table