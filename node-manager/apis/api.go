package apis

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dipghoshraj/media-service/node-manager/handler"
	"github.com/gorilla/mux"
)

type NMHandler struct {
	DbManager *handler.DBHandler
}

func SetupRoutes(router *mux.Router, nm *NMHandler) {
	// Node management endpoints
	router.HandleFunc("/api/nodes", nm.RegisterNodeHandler).Methods("POST")
	router.HandleFunc("/api/nodes/stats", nm.GetClusterStatsHandler).Methods("GET")
	router.HandleFunc("/api/all/nodes", nm.GetAllNodesHandler).Methods("GET")
	router.HandleFunc("/api/node", nm.GetNodeDetails).Methods("GET")

	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				respondWithError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
