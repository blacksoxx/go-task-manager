package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "taskservice/internal/database"
    "taskservice/internal/models"

    "github.com/gorilla/mux"
)

type TaskHandler struct {
    db *database.DB
}

func NewTaskHandler(db *database.DB) *TaskHandler {
    return &TaskHandler{db: db}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req models.CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    // Validate required fields
    if req.Title == "" || req.UserID == "" {
        http.Error(w, `{"error": "Title and user_id are required"}`, http.StatusBadRequest)
        return
    }

    task := &models.Task{
        Title:       req.Title,
        Description: req.Description,
        UserID:      req.UserID,
        DueDate:     req.DueDate,
        Status:      models.StatusPending,
    }

    if err := h.db.CreateTask(task); err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(models.TaskResponse{Task: task}); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID := vars["id"]

    task, err := h.db.GetTaskByID(taskID)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(models.TaskResponse{Task: task}); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]

    tasks, err := h.db.GetTasksByUserID(userID)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
        return
    }

    response := models.TasksResponse{
        Tasks: tasks,
        Total: len(tasks),
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID := vars["id"]

    var req models.UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
        return
    }

    task, err := h.db.UpdateTask(taskID, &req)
    if err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(models.TaskResponse{Task: task}); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID := vars["id"]

    if err := h.db.DeleteTask(taskID); err != nil {
        http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status":    "healthy",
        "service":   "taskservice",
        "timestamp": time.Now().Format(time.RFC3339),
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, `{"error": "Error encoding response"}`, http.StatusInternalServerError)
    }
}
