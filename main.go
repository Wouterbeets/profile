package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"testserver/templates"
)

type GitHubRepo struct {
	Stars int `json:"stargazers_count"`
	Forks int `json:"forks_count"`
}

func main() {
	// Create a new Chi router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(5)) // GZIP compression
	router.Use(middleware.Recoverer)

	// Static files
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/manifest.json", http.FileServer(http.Dir(".")))
	router.Handle("/sw.js", http.FileServer(http.Dir(".")))

	// Load data from JSON
	experienceData := loadData("data/experience.json")
	educationData := loadData("data/education.json")
	projectsData := loadData("data/projects.json")

	// Handle root route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Skills": []string{
				"Go Programming", "Python", "C", "React",
				"Microservices Architecture", "Cloud Technologies",
				"DevOps Practices", "Database Design", "System Design",
				"Leadership", "Team Management", "Technical Vision",
			},
		}
		templates.IndexTemplate(data).Render(r.Context(), w)
	})

	// Handle experience section
	router.Get("/cv/experience", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templates.ExperienceTemplate(experienceData).Render(r.Context(), w)
	})

	// Handle education section
	router.Get("/cv/education", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templates.EducationTemplate(educationData).Render(r.Context(), w)
	})

	// Handle projects section
	router.Get("/cv/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		templates.ProjectsTemplate(projectsData).Render(r.Context(), w)
	})

	// New: Handle contact form submission
	router.Post("/contact-submit", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")
		if name == "" || email == "" || message == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		// Send email (configure SMTP)
		err := sendEmail(name, email, message)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully!"))
	})

	// New: Handle GitHub stats
	router.Get("/api/github-stats/{repo}", func(w http.ResponseWriter, r *http.Request) {
		repo := chi.URLParam(r, "repo")
		token := os.Getenv("GITHUB_TOKEN") // Set your GitHub token
		url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "token "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var repoData GitHubRepo
		json.Unmarshal(body, &repoData)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"stars": repoData.Stars, "forks": repoData.Forks})
	})

	// New: Handle skills filter
	router.Get("/cv/skills", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		allSkills := []string{
			"Go Programming", "Python", "C", "React",
			"Microservices Architecture", "Cloud Technologies",
			"DevOps Practices", "Database Design", "System Design",
			"Leadership", "Team Management", "Technical Vision",
		}
		filtered := []string{}
		for _, skill := range allSkills {
			if strings.Contains(strings.ToLower(skill), strings.ToLower(query)) {
				filtered = append(filtered, skill)
			}
		}
		w.Header().Set("Content-Type", "text/html")
		for _, skill := range filtered {
			fmt.Fprintf(w, `<span class="skill-tag animate__animated animate__fadeIn">%s</span>`, skill)
		}
	})

	// Create server
	server := &http.Server{
		Addr:    ":33333",
		Handler: router,
	}

	log.Printf("Starting server on port 33333...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loadData(filename string) map[string]interface{} {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v", filename, err)
		return map[string]interface{}{}
	}
	defer file.Close()
	var data map[string]interface{}
	json.NewDecoder(file).Decode(&data)
	return data
}

func sendEmail(name, email, message string) error {
	from := "your-email@example.com" // Configure
	to := "recipient@example.com"
	smtpHost := "smtp.example.com" // Configure
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, "password", smtpHost) // Configure
	msg := fmt.Sprintf("Subject: Contact from %s\n\nName: %s\nEmail: %s\nMessage: %s", name, name, email, message)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
