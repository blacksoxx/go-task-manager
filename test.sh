#!/bin/bash

echo "🧪 Testing Notification Service"

# Base URL
BASE_URL="http://localhost:8083"
SERVICE_NAME="notification-service"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Wait for service to be ready
echo "⏳ Waiting for notification service to be ready..."
until curl -s "$BASE_URL/health" > /dev/null; do
    sleep 1
done

# Test 1: Health Check
echo "📋 Test 1: Health Check"
HEALTH_RESPONSE=$(curl -s -w "%{http_code}" "$BASE_URL/health")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Health check passed"
else
    print_result 1 "Health check failed"
fi

# Test 2: Create Notification
echo "📋 Test 2: Create Notification"
CREATE_RESPONSE=$(curl -s -w "%{http_code}" \
    -X POST "$BASE_URL/api/v1/notifications" \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "user-123",
        "title": "Test Notification",
        "message": "This is a test notification",
        "type": "in_app",
        "data": {
            "task_id": "task-456",
            "priority": "high"
        }
    }')

HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$CREATE_RESPONSE" | head -n -1)

if [ "$HTTP_CODE" -eq 201 ]; then
    print_result 0 "Create notification passed"
    NOTIFICATION_ID=$(echo "$RESPONSE_BODY" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "   Created notification ID: $NOTIFICATION_ID"
else
    print_result 1 "Create notification failed"
    echo "   Response: $RESPONSE_BODY"
fi

# Test 3: Get Notification
if [ ! -z "$NOTIFICATION_ID" ]; then
    echo "📋 Test 3: Get Notification by ID"
    GET_RESPONSE=$(curl -s -w "%{http_code}" "$BASE_URL/api/v1/notifications/$NOTIFICATION_ID")
    HTTP_CODE=$(echo "$GET_RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" -eq 200 ]; then
        print_result 0 "Get notification passed"
    else
        print_result 1 "Get notification failed"
    fi
fi

# Test 4: Get User Notifications
echo "📋 Test 4: Get User Notifications"
USER_NOTIFICATIONS_RESPONSE=$(curl -s -w "%{http_code}" "$BASE_URL/api/v1/users/user-123/notifications")
HTTP_CODE=$(echo "$USER_NOTIFICATIONS_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" -eq 200 ]; then
    print_result 0 "Get user notifications passed"
else
    print_result 1 "Get user notifications failed"
fi

# Test 5: Mark as Read
if [ ! -z "$NOTIFICATION_ID" ]; then
    echo "📋 Test 5: Mark Notification as Read"
    MARK_READ_RESPONSE=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/api/v1/notifications/$NOTIFICATION_ID/read")
    HTTP_CODE=$(echo "$MARK_READ_RESPONSE" | tail -n1)
    
    if [ "$HTTP_CODE" -eq 204 ]; then
        print_result 0 "Mark as read passed"
    else
        print_result 1 "Mark as read failed"
    fi
fi

# Test 6: Create Invalid Notification (Validation Test)
echo "📋 Test 6: Create Invalid Notification (Validation)"
INVALID_RESPONSE=$(curl -s -w "%{http_code}" \
    -X POST "$BASE_URL/api/v1/notifications" \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "",
        "title": "",
        "message": "",
        "type": "invalid_type"
    }')

HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" -eq 400 ]; then
    print_result 0 "Validation test passed"
else
    print_result 1 "Validation test failed"
fi

echo ""
echo "🎯 Notification Service Testing Complete!"