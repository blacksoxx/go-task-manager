// Wait for DOM to be fully loaded before initializing the app
document.addEventListener('DOMContentLoaded', function() {
    // Initialize the app when the page loads
    const app = new TaskManagerApp();
});

class TaskManagerApp {
    constructor() {
        this.authBaseUrl = 'http://localhost:8084/api/v1';
        this.notificationBaseUrl = 'http://localhost:8083/api/v1';
        this.taskBaseUrl = 'http://localhost:8082/api/v1';
        this.currentUser = null;
        this.init();
    }

    init() {
        this.bindAuthEvents();
        this.checkExistingAuth();
    }

    bindAuthEvents() {
        // Tab switching for auth
        const tabButtons = document.querySelectorAll('.tab-button');
        if (tabButtons.length > 0) {
            tabButtons.forEach(button => {
                button.addEventListener('click', (e) => {
                    this.switchTab(e.target.dataset.tab);
                });
            });
        }

        // Login form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.login();
            });
        }

        // Signup form
        const signupForm = document.getElementById('signupForm');
        if (signupForm) {
            signupForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.signup();
            });
        }
    }

    switchTab(tab) {
        // Update active tab buttons
        document.querySelectorAll('.tab-button').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tab === tab);
        });

        // Show active form
        document.querySelectorAll('.auth-form').forEach(form => {
            form.classList.toggle('active', form.id === `${tab}Form`);
        });

        // Clear messages
        this.clearAuthMessage();
    }

    async signup() {
        const firstName = document.getElementById('signupFirstName').value;
        const lastName = document.getElementById('signupLastName').value;
        const email = document.getElementById('signupEmail').value;
        const password = document.getElementById('signupPassword').value;

        try {
            const response = await fetch(`${this.authBaseUrl}/auth/signup`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    first_name: firstName, 
                    last_name: lastName, 
                    email: email, 
                    password: password 
                })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Signup failed');
            }

            const data = await response.json();
            this.handleAuthSuccess(data);

        } catch (error) {
            this.showAuthMessage(error.message, 'error');
        }
    }

    async login() {
        const email = document.getElementById('loginEmail').value;
        const password = document.getElementById('loginPassword').value;

        try {
            const response = await fetch(`${this.authBaseUrl}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Login failed');
            }

            const data = await response.json();
            this.handleAuthSuccess(data);

        } catch (error) {
            this.showAuthMessage(error.message, 'error');
        }
    }

    handleAuthSuccess(authData) {
        this.currentUser = authData.user;
        
        // Store user data in localStorage
        localStorage.setItem('currentUser', JSON.stringify(this.currentUser));
        localStorage.setItem('authToken', authData.token);

        // Switch to dashboard
        this.showDashboard();
        
        // Load user's tasks and notifications
        this.loadTasks();
        this.loadNotifications();
    }

    showDashboard() {
        document.getElementById('authScreen').style.display = 'none';
        document.getElementById('dashboard').style.display = 'block';
        
        // Update user info
        document.getElementById('userName').textContent = 
            `${this.currentUser.first_name} ${this.currentUser.last_name}`;
        
        // Initialize dashboard events
        this.bindDashboardEvents();
        
        // Ensure only tasks tab is visible initially
        this.switchDashboardTab('tasks');
    }

    bindDashboardEvents() {
        // Dashboard tab switching
        document.querySelectorAll('.dashboard-tab').forEach(tab => {
            tab.addEventListener('click', (e) => {
                this.switchDashboardTab(e.target.dataset.tab);
            });
        });

        // Task management
        const createTaskBtn = document.getElementById('createTaskBtn');
        if (createTaskBtn) {
            createTaskBtn.addEventListener('click', () => {
                this.toggleCreateTaskForm();
            });
        }

        const cancelTaskBtn = document.getElementById('cancelTaskBtn');
        if (cancelTaskBtn) {
            cancelTaskBtn.addEventListener('click', () => {
                this.toggleCreateTaskForm();
            });
        }

        const taskForm = document.getElementById('taskForm');
        if (taskForm) {
            taskForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.createTask();
            });
        }

        // Notifications
        const refreshNotificationsBtn = document.getElementById('refreshNotificationsBtn');
        if (refreshNotificationsBtn) {
            refreshNotificationsBtn.addEventListener('click', () => {
                this.loadNotifications();
            });
        }

        // Logout
        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', () => {
                this.logout();
            });
        }

        // Modal close
        const closeBtn = document.querySelector('.close');
        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                this.closeModal();
            });
        }

        const notificationModal = document.getElementById('notificationModal');
        if (notificationModal) {
            notificationModal.addEventListener('click', (e) => {
                if (e.target.id === 'notificationModal') {
                    this.closeModal();
                }
            });
        }
    }

    switchDashboardTab(tab) {
        console.log('Switching to tab:', tab);
        
        // Update active tab buttons
        document.querySelectorAll('.dashboard-tab').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tab === tab);
        });

        // Hide all tab contents first
        document.querySelectorAll('.tab-content').forEach(content => {
            content.style.display = 'none';
            content.classList.remove('active');
        });

        // Show the active tab content
        const activeTab = document.getElementById(`${tab}Tab`);
        if (activeTab) {
            activeTab.style.display = 'block';
            activeTab.classList.add('active');
        }

        // Load data for the active tab
        if (tab === 'tasks') {
            this.loadTasks();
        } else if (tab === 'notifications') {
            this.loadNotifications();
        }
    }

    // Task Management Methods
    async loadTasks() {
        if (!this.currentUser) return;

        const loading = document.getElementById('tasksLoading');
        const list = document.getElementById('tasksList');
        
        if (loading) loading.style.display = 'block';
        if (list) list.innerHTML = '';

        try {
            const response = await fetch(`${this.taskBaseUrl}/users/${this.currentUser.id}/tasks`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            this.displayTasks(data.tasks);
        } catch (error) {
            console.error('Error loading tasks:', error);
            if (list) {
                list.innerHTML = `
                    <div style="text-align: center; padding: 40px; color: #dc3545;">
                        <h3>❌ Error loading tasks</h3>
                        <p>${error.message}</p>
                        <p>Make sure the task service is running on port 8082</p>
                    </div>
                `;
            }
        } finally {
            if (loading) loading.style.display = 'none';
        }
    }

    displayTasks(tasks) {
        const list = document.getElementById('tasksList');
        if (!list) return;
        
        if (!tasks || tasks.length === 0) {
            list.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #6c757d;">
                    <h3>📝 No tasks yet</h3>
                    <p>Create your first task using the button above!</p>
                </div>
            `;
            return;
        }

        list.innerHTML = tasks.map(task => `
            <div class="task-card">
                <div class="task-header">
                    <div class="task-content">
                        <div class="task-title">${this.escapeHtml(task.title)}</div>
                        ${task.description ? `
                            <div class="task-description">${this.escapeHtml(task.description)}</div>
                        ` : ''}
                    </div>
                    <span class="task-status status-${task.status}">${task.status}</span>
                </div>
                <div class="task-meta">
                    <div class="task-dates">
                        <span class="created-date"><strong>Created:</strong> ${new Date(task.created_at).toLocaleString()}</span>
                        ${task.due_date ? `
                            <span class="due-date">• <strong>Due:</strong> ${new Date(task.due_date).toLocaleString()}</span>
                        ` : ''}
                    </div>
                </div>
            </div>
        `).join('');
    }

    async createTask() {
        if (!this.currentUser) return;

        const title = document.getElementById('taskTitle').value;
        const description = document.getElementById('taskDescription').value;
        const dueDate = document.getElementById('taskDueDate').value;

        const taskData = {
            title: title,
            description: description,
            user_id: this.currentUser.id,
            due_date: dueDate ? new Date(dueDate).toISOString() : null
        };

        try {
            const response = await fetch(`${this.taskBaseUrl}/tasks`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(taskData)
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Failed to create task');
            }

            const result = await response.json();
            
            // Create automatic notification for the new task
            await this.createTaskNotification(result.task);
            
            this.toggleCreateTaskForm();
            this.loadTasks();
            this.loadNotifications(); // Refresh notifications to show the new one
            
            alert('✅ Task created successfully!');
            
        } catch (error) {
            console.error('Error creating task:', error);
            alert(`Error: ${error.message}`);
        }
    }

    async createTaskNotification(task) {
        try {
            const notificationData = {
                user_id: this.currentUser.id,
                title: 'New Task Created',
                message: `You created a new task: "${task.title}"`,
                type: 'in_app',
                data: {
                    task_id: task.id,
                    task_title: task.title,
                    action: 'task_created'
                }
            };

            await fetch(`${this.notificationBaseUrl}/notifications`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(notificationData)
            });
        } catch (error) {
            console.error('Error creating task notification:', error);
            // Don't show error to user - notification failure shouldn't block task creation
        }
    }

    toggleCreateTaskForm() {
        const form = document.getElementById('createTaskForm');
        if (!form) return;
        
        const isVisible = form.style.display === 'block';
        form.style.display = isVisible ? 'none' : 'block';
        
        if (!isVisible) {
            // Reset form
            const taskForm = document.getElementById('taskForm');
            if (taskForm) taskForm.reset();
        }
    }

    // Notification Methods
    async loadNotifications() {
        if (!this.currentUser) return;

        const loading = document.getElementById('notificationsLoading');
        const list = document.getElementById('notificationsList');
        
        if (loading) loading.style.display = 'block';
        if (list) list.innerHTML = '';

        try {
            const response = await fetch(`${this.notificationBaseUrl}/users/${this.currentUser.id}/notifications?limit=50`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            this.displayNotifications(data.notifications);
        } catch (error) {
            console.error('Error loading notifications:', error);
            if (list) {
                list.innerHTML = `
                    <div style="text-align: center; padding: 40px; color: #dc3545;">
                        <h3>❌ Error loading notifications</h3>
                        <p>${error.message}</p>
                    </div>
                `;
            }
        } finally {
            if (loading) loading.style.display = 'none';
        }
    }

    displayNotifications(notifications) {
        const list = document.getElementById('notificationsList');
        if (!list) return;
        
        if (!notifications || notifications.length === 0) {
            list.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #6c757d;">
                    <h3>📭 No notifications yet</h3>
                    <p>You'll see notifications here when you create tasks or when tasks are assigned to you.</p>
                </div>
            `;
            return;
        }

        list.innerHTML = notifications.map(notification => `
            <div class="notification-card" onclick="app.showNotificationDetails('${notification.id}')">
                <div class="notification-header">
                    <div class="notification-content">
                        <div class="notification-title">${this.escapeHtml(notification.title)}</div>
                        <div class="notification-message">${this.escapeHtml(notification.message)}</div>
                    </div>
                    <span class="notification-type type-${notification.type}">${notification.type}</span>
                </div>
                <div class="notification-meta">
                    <div class="notification-dates">
                        <span class="created-date"><strong>Created:</strong> ${new Date(notification.created_at).toLocaleString()}</span>
                        ${notification.read_at ? `
                            <span class="read-date">• <strong>Read:</strong> ${new Date(notification.read_at).toLocaleString()}</span>
                        ` : ''}
                    </div>
                    <span class="notification-status status-${notification.status}">${notification.status}</span>
                </div>
            </div>
        `).join('');
    }

    async showNotificationDetails(notificationId) {
        try {
            const response = await fetch(`${this.notificationBaseUrl}/notifications/${notificationId}`);
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
        
        if (!modal || !details) return;
        
        details.innerHTML = `
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

    async markAsRead(notificationId) {
        try {
            const response = await fetch(`${this.notificationBaseUrl}/notifications/${notificationId}/read`, {
                method: 'PUT'
            });

            if (!response.ok) throw new Error('Failed to mark as read');
            
            this.closeModal();
            this.loadNotifications();
            
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
            const response = await fetch(`${this.notificationBaseUrl}/notifications/${notificationId}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete notification');
            
            this.closeModal();
            this.loadNotifications();
            
        } catch (error) {
            console.error('Error deleting notification:', error);
            alert('Error deleting notification');
        }
    }

    closeModal() {
        const modal = document.getElementById('notificationModal');
        if (modal) modal.style.display = 'none';
    }

    // Utility Methods
    checkExistingAuth() {
        const storedUser = localStorage.getItem('currentUser');
        const storedToken = localStorage.getItem('authToken');
        
        if (storedUser && storedToken) {
            this.currentUser = JSON.parse(storedUser);
            this.showDashboard();
            this.loadTasks();
            this.loadNotifications();
        }
    }

    logout() {
        this.currentUser = null;
        localStorage.removeItem('currentUser');
        localStorage.removeItem('authToken');
        
        const dashboard = document.getElementById('dashboard');
        const authScreen = document.getElementById('authScreen');
        
        if (dashboard) dashboard.style.display = 'none';
        if (authScreen) authScreen.style.display = 'flex';
        
        // Clear forms
        const loginForm = document.getElementById('loginForm');
        const signupForm = document.getElementById('signupForm');
        
        if (loginForm) loginForm.reset();
        if (signupForm) signupForm.reset();
        
        this.clearAuthMessage();
    }

    showAuthMessage(message, type) {
        const messageEl = document.getElementById('authMessage');
        if (messageEl) {
            messageEl.textContent = message;
            messageEl.className = `auth-message ${type}`;
            messageEl.style.display = 'block';
        }
    }

    clearAuthMessage() {
        const messageEl = document.getElementById('authMessage');
        if (messageEl) {
            messageEl.textContent = '';
            messageEl.className = 'auth-message';
            messageEl.style.display = 'none';
        }
    }

    escapeHtml(unsafe) {
        if (!unsafe) return '';
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}