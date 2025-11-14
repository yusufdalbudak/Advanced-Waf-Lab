package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var users = []map[string]string{
	{"id": "1", "name": "Alice", "email": "alice@example.com"},
	{"id": "2", "name": "Bob", "email": "bob@example.com"},
	{"id": "3", "name": "Charlie", "email": "charlie@example.com"},
}

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("cmd/test-website/templates/*.html")
	if err != nil {
		log.Fatal("Failed to parse templates:", err)
	}
}

func main() {
	port := ":8081"

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("cmd/test-website/static"))))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/users", usersAPIHandler)
	http.HandleFunc("/api/search", searchAPIHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/command", commandAPIHandler)
	http.HandleFunc("/api/upload", uploadAPIHandler)
	http.HandleFunc("/search", searchPageHandler)
	http.HandleFunc("/users", usersPageHandler)
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/files", filesHandler)
	http.HandleFunc("/upload", uploadPageHandler)
	http.HandleFunc("/command", commandPageHandler)

	fmt.Printf("Test website starting on http://localhost%s\n", port)
	fmt.Println("This is a vulnerable test application for WAF testing")
	log.Fatal(http.ListenAndServe(port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "WAF Test Website",
		"Time":  time.Now().Format("2006-01-02 15:04:05"),
	}
	renderTemplate(w, "home.html", data)
}

func searchPageHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	data := map[string]interface{}{
		"Title": "Search",
		"Query": query,
		"Results": []string{"Result 1", "Result 2", "Result 3"},
	}
	renderTemplate(w, "search.html", data)
}

func usersPageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	data := map[string]interface{}{
		"Title":  "Users",
		"Users":  users,
		"UserID": userID,
	}
	renderTemplate(w, "users.html", data)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		_ = r.FormValue("password") // Password received but not used (test site)
		data := map[string]interface{}{
			"Title":    "Login",
			"Username": username,
			"Message":  "Login attempt received (this is a test site)",
		}
		renderTemplate(w, "login.html", data)
		return
	}
	data := map[string]interface{}{
		"Title": "Login",
	}
	renderTemplate(w, "login.html", data)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	data := map[string]interface{}{
		"Title":    "Files",
		"FilePath": filePath,
	}
	renderTemplate(w, "files.html", data)
}

func uploadPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			data := map[string]interface{}{
				"Title":    "Upload",
				"FileName": header.Filename,
				"FileSize": header.Size,
				"Message":  "File upload received (test site)",
			}
			renderTemplate(w, "upload.html", data)
			return
		}
	}
	data := map[string]interface{}{
		"Title": "Upload",
	}
	renderTemplate(w, "upload.html", data)
}

func commandPageHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")
	data := map[string]interface{}{
		"Title":   "Command",
		"Command": cmd,
	}
	renderTemplate(w, "command.html", data)
}

func usersAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("id")
	
	if userID != "" {
		for _, user := range users {
			if user["id"] == userID {
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}

func searchAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("q")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"results": []string{"result1", "result2", "result3"},
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	username := r.URL.Query().Get("username")
	_ = r.URL.Query().Get("password") // Password received but not used (test site)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "received",
		"username": username,
		"message":  "Login attempt processed (test site)",
	})
}

func commandAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cmd := r.URL.Query().Get("cmd")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"command": cmd,
		"status":  "received",
		"message": "Command execution simulated (test site)",
	})
}

func uploadAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":   "received",
				"filename": header.Filename,
				"size":     header.Size,
				"message":  "File upload processed (test site)",
			})
			return
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "No file uploaded",
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

