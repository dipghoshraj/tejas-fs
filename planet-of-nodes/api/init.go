package api

import (
	"encoding/json"
	"log"
	"net/http"
	handlers "planet-of-node/handler"
	"time"

	"github.com/gorilla/mux"
)

type NApi struct {
	nhm *handlers.HManager
}

func ApiHandler(hm *handlers.HManager) *NApi {
	return &NApi{nhm: hm}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SetUpRouter(router *mux.Router, napi *NApi) {
	router.Use(loggingMiddleware)
	router.HandleFunc("/cluster", napi.CreateCluster).Methods("POST")
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, APIResponse{Success: true, Message: message})
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
