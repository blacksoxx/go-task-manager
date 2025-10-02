class NotificationApp {
    constructor() {
        this.baseUrl = 'http://localhost:8083/api/v1';
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadNotifications();
    }

    bindEvents() {
        // Refresh button
        document.getElementById('refreshBtn').addEventListener('click', () => {
            this.loadNotifications();
        });

        // Create notification button
        document.getElementById('createBtn').addEventListener('click', () => {
            this.toggleCreateForm();
        });

        // Cancel create form
        document.getElementById('cancelBtn').addEventListener('click', () => {
            this.toggleCreateForm();
        });

        // Submit create form
        document.getElementById('notificationForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createNotification();
        });

        // Close modal
        document.querySelector('.close').addEventListener('click', () => {
            this.closeModal();
        });

        // Close modal when clicking outside
        document.getElementById('notificationModal').addEventListener('click', (e) => {
            if (e.target.id === 'notificationModal') {
                this.closeModal();
            }
        });
    }

    async loadNotifications() {
        const loading = document.getElementById('loading');
        const list = document.getElementById('notificationsList');
        
        loading.style.display = 'block';
        list.innerHTML = '';

        try {
            // For demo purposes, we'll use a fixed user ID
            // In a real app, you'd get this from authentication
            const userId = '123e4567-e89b-12d3-a456-426614174000';
            const response = await fetch(`${this.baseUrl}/users/${userId}/notifications?limit=50`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            this.displayNotifications(data.notifications);
        } catch (error) {
            console.error('Error loading notifications:', error);
            list.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #dc3545;">
                    <h3>❌ Error loading notifications</h3>
                    <p>${error.message}</p>
                    <p>Make sure the notification service is running on port 8083</p>
                </div>
            `;
        } finally {
            loading.style.display = 'none';
        }
    }

    displayNotifications(notifications) {
        const list = document.getElementById('notificationsList');
        
        if (!notifications || notifications.length === 0) {
            list.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #6c757d;">
                    <h3>📭 No notifications found</h3>
                    <p>Create your first notification using the button above!</p>
                </div>
            `;
            return;
        }

        list.innerHTML = notifications.map(notification => `
            <div class="notification-card" onclick="app.showNotificationDetails('${notification.id}')">
                <div class="notification-header">
                    <div>
                        <div class="notification-title">${this.escapeHtml(notification.title)}</div>
                        <div class="notification-message">${this.escapeHtml(notification.message)}</div>
                    </div>
                    <span class="notification-type type-${notification.type}">${notification.type}</span>
                </div>
                <div class="notification-meta">
                    <div>
                        <strong>User:</strong> ${notification.user_id.substring(0, 8)}...
                        <strong>Created:</strong> ${new Date(notification.created_at).toLocaleString()}
                    </div>
                    <span class="notification-status status-${notification.status}">${notification.status}</span>
                </div>
                ${notification.read_at ? `
                    <div class="notification-meta">
                        <small><strong>Read:</strong> ${new Date(notification.read_at).toLocaleString()}</small>
                    </div>
                ` : ''}
            </div>
        `).join('');
    }

    async showNotificationDetails(notificationId) {
        try {
            const response = await fetch(`${this.baseUrl}/notifications/${notificationId}`);
            if (!response.ok) throw new Error('Notification not found');
            
            const data = await response.json();
            this.displayNotificationDetails(data.notification);
        } catch (error) {
            console.error('Error loading notification details:', error);
            alert('Error loading notification details');
        }
    }

    displayNotificationDetails(notification) {
        const modal = document.getElementById('notificationModal');
        const details = document.getElementById('notificationDetails');
        
        details.innerHTML = `
            <div class="detail-item">
                <strong>ID:</strong> ${notification.id}
            </div>
            <div class="detail-item">
                <strong>User ID:</strong> ${notification.user_id}
            </div>
            <div class="detail-item">
                <strong>Title:</strong> ${this.escapeHtml(notification.title)}
            </div>
            <div class="detail-item">
                <strong>Message:</strong> ${this.escapeHtml(notification.message)}
            </div>
            <div class="detail-item">
                <strong>Type:</strong> <span class="notification-type type-${notification.type}">${notification.type}</span>
            </div>
            <div class="detail-item">
                <strong>Status:</strong> <span class="notification-status status-${notification.status}">${notification.status}</span>
            </div>
            <div class="detail-item">
                <strong>Created:</strong> ${new Date(notification.created_at).toLocaleString()}
            </div>
            <div class="detail-item">
                <strong>Updated:</strong> ${new Date(notification.updated_at).toLocaleString()}
            </div>
            ${notification.read_at ? `
                <div class="detail-item">
                    <strong>Read At:</strong> ${new Date(notification.read_at).toLocaleString()}
                </div>
            ` : ''}
            <div class="detail-item">
                <strong>Additional Data:</strong><br>
                <pre>${JSON.stringify(notification.data, null, 2)}</pre>
            </div>
            <div class="form-actions">
                ${notification.status !== 'read' ? `
                    <button onclick="app.markAsRead('${notification.id}')">Mark as Read</button>
                ` : ''}
                <button onclick="app.deleteNotification('${notification.id}')" style="background: #dc3545; color: white;">Delete</button>
            </div>
        `;
        
        modal.style.display = 'flex';
    }

    async createNotification() {
        const form = document.getElementById('notificationForm');
        const formData = new FormData(form);
        
        let data = {};
        const dataField = document.getElementById('data').value;
        if (dataField.trim()) {
            try {
                data = JSON.parse(dataField);
            } catch (e) {
                alert('Invalid JSON in data field');
                return;
            }
        }

        const notificationData = {
            user_id: document.getElementById('user_id').value,
            title: document.getElementById('title').value,
            message: document.getElementById('message').value,
            type: document.getElementById('type').value,
            data: data
        };

        try {
            const response = await fetch(`${this.baseUrl}/notifications`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(notificationData)
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create notification');
            }

            const result = await response.json();
            this.toggleCreateForm();
            this.loadNotifications();
            alert('✅ Notification created successfully!');
            
        } catch (error) {
            console.error('Error creating notification:', error);
            alert(`Error: ${error.message}`);
        }
    }

    async markAsRead(notificationId) {
        try {
            const response = await fetch(`${this.baseUrl}/notifications/${notificationId}/read`, {
                method: 'PUT'
            });

            if (!response.ok) throw new Error('Failed to mark as read');
            
            this.closeModal();
            this.loadNotifications();
            alert('✅ Notification marked as read!');
            
        } catch (error) {
            console.error('Error marking as read:', error);
            alert('Error marking notification as read');
        }
    }

    async deleteNotification(notificationId) {
        if (!confirm('Are you sure you want to delete this notification?')) {
            return;
        }

        try {
            const response = await fetch(`${this.baseUrl}/notifications/${notificationId}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete notification');
            
            this.closeModal();
            this.loadNotifications();
            alert('✅ Notification deleted successfully!');
            
        } catch (error) {
            console.error('Error deleting notification:', error);
            alert('Error deleting notification');
        }
    }

    toggleCreateForm() {
        const form = document.getElementById('createForm');
        form.style.display = form.style.display === 'none' ? 'block' : 'none';
    }

    closeModal() {
        document.getElementById('notificationModal').style.display = 'none';
    }

    escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}

// Initialize the app when the page loads
const app = new NotificationApp();