package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := ":8081"
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		response := map[string]interface{}{
			"status":    "ok",
			"message":   "Request received by backend",
			"path":      r.URL.Path,
			"method":    r.Method,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		
		// Add query parameters if present
		if len(r.URL.Query()) > 0 {
			response["query"] = r.URL.Query()
		}
		
		json.NewEncoder(w).Encode(response)
	})
	
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": []map[string]string{
				{"id": "1", "name": "Alice"},
				{"id": "2", "name": "Bob"},
			},
		})
	})
	
	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query().Get("q")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"query":   query,
			"results": []string{"result1", "result2"},
		})
	})
	
	fmt.Printf("Test backend server starting on http://localhost%s\n", port)
	fmt.Println("This server will receive requests forwarded by the WAF")
	log.Fatal(http.ListenAndServe(port, nil))
}

