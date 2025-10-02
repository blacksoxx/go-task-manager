-- Create database for notification service
CREATE DATABASE notificationdb;

-- Connect to notificationdb
\c notificationdb;

-- Create notifications table (this matches our Go schema)
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('email', 'in_app', 'push')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'read')),
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);

-- Insert sample data for testing
INSERT INTO notifications (user_id, title, message, type, status, data) VALUES
('user-123', 'Welcome to TaskManager!', 'Your account has been created successfully', 'in_app', 'read', '{"action": "welcome"}'),
('user-123', 'Task Assigned', 'You have been assigned a new task: "Design Dashboard"', 'email', 'sent', '{"task_id": "task-789", "action": "task_assigned"}'),
('user-456', 'Reminder', 'Your task "Update Documentation" is due tomorrow', 'push', 'pending', '{"task_id": "task-999", "due_date": "2024-01-15"}');

-- Verify the data
SELECT * FROM notifications;