cd "/mnt/c/Users/User/Documents/PROJECTS/Cloud-Native Go Task Manager"

# Test both services
echo "=== Testing User Service ==="
curl http://localhost:8081/health

echo ""
echo "=== Testing Task Service ==="
curl http://localhost:8082/health

echo ""
echo "=== Creating Integration Test ==="
# Create a user
USER_JSON=$(curl -s -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "integration@example.com",
    "first_name": "Integration",
    "last_name": "Test",
    "password": "integ123"
  }')

echo "User created: $USER_JSON"

# Extract user ID (simple method - in production use jq)
USER_ID=$(echo $USER_JSON | grep -o '"id":"[^"]*' | cut -d'"' -f4)
if [ -n "$USER_ID" ]; then
    echo "User ID extracted: $USER_ID"
    
    # Create a task for this user
    echo ""
    echo "Creating task for user $USER_ID..."
    curl -X POST http://localhost:8082/api/v1/tasks \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Microservices Integration Task",
        "description": "This task was created via API integration between User and Task services",
        "user_id": "'$USER_ID'"
      }'
else
    echo "Failed to extract user ID"
fi