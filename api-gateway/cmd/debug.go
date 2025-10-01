
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func debugRoutes() {
    r := mux.NewRouter()
    
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users/{user_id}/tasks", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID := vars["user_id"]
        fmt.Fprintf(w, "Debug: User ID = %s", userID)
    }).Methods("GET")
    
    fmt.Println("Registered routes:")
    r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        t, err := route.GetPathTemplate()
        if err == nil {
            fmt.Println("Route:", t)
        }
        return nil
    })
}

func main() {
    debugRoutes()
}
