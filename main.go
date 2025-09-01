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

// Helper to detect language from query param or cookie
func detectLanguage(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		if cookie, err := r.Cookie("language"); err == nil {
			lang = cookie.Value
		}
	}
	if lang != "en" && lang != "fr" {
		lang = "en" // Default
	}
	return lang
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

	// Handle root route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		data := templates.IndexData{
			Skills: []string{
				"Go Programming", "Python", "C", "React",
				"Microservices Architecture", "Cloud Technologies",
				"DevOps Practices", "Database Design", "System Design",
				"Leadership", "Team Management", "Technical Vision",
				"Concurrent Programming in Go", "Fyne Framework", "Event-Driven Programming",
				"LLM Integration with Ollama", "Real-Time Audio Transcription", "Event Sourcing & CQRS", "Plugin Architecture",
				"Cross-Language Integration", "Panic Recovery & Error Handling", "File-Based Persistence",
				"Desktop Application Development", "Markdown Parsing for UI", "JSON Data Handling", "Logging & Debugging",
			},
			Translations: templates.Translations,
			Language:    lang,
		}
		templates.IndexTemplate(data).Render(r.Context(), w)
	})

	// Handle experience section
	router.Get("/cv/experience", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		w.Header().Set("Content-Type", "text/html")
		data := loadExperienceData(lang)
		data.Language = lang
		data.Translations = templates.Translations
		templates.ExperienceTemplate(data).Render(r.Context(), w)
	})

	// Handle education section
	router.Get("/cv/education", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		w.Header().Set("Content-Type", "text/html")
		data := loadEducationData(lang)
		data.Language = lang
		data.Translations = templates.Translations
		templates.EducationTemplate(data).Render(r.Context(), w)
	})

	// Handle projects section
	router.Get("/cv/projects", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		w.Header().Set("Content-Type", "text/html")
		data := loadProjectsData(lang)
		data.Language = lang
		data.Translations = templates.Translations
		templates.ProjectsTemplate(data).Render(r.Context(), w)
	})

	// Handle contact form
	router.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		w.Header().Set("Content-Type", "text/html")
		templates.ContactTemplate(lang, templates.Translations).Render(r.Context(), w)
	})

	// New: Handle contact form submission
	router.Post("/contact-submit", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")
		if name == "" || email == "" || message == "" {
			http.Error(w, templates.GetTranslation("all_fields_required", lang), http.StatusBadRequest)
			return
		}
		// Send email (configure SMTP)
		err := sendEmail(name, email, message)
		if err != nil {
			http.Error(w, templates.GetTranslation("failed_send", lang), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(templates.GetTranslation("message_sent", lang)))
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
			"Concurrent Programming in Go", "Fyne Framework", "Event-Driven Programming",
			"LLM Integration with Ollama", "Real-Time Audio Transcription", "Event Sourcing & CQRS", "Plugin Architecture",
			"Cross-Language Integration", "Panic Recovery & Error Handling", "File-Based Persistence",
			"Desktop Application Development", "Markdown Parsing for UI", "JSON Data Handling", "Logging & Debugging",
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

func loadExperienceData(lang string) templates.ExperienceData {
	filename := fmt.Sprintf("data/experience_%s.json", lang)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/experience_en.json"
		file, err = os.Open(filename)
		if err != nil {
			log.Printf("Error loading fallback %s: %v", filename, err)
			return templates.ExperienceData{}
		}
	}
	defer file.Close()
	var data templates.ExperienceData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Printf("Error decoding %s: %v", filename, err)
		return templates.ExperienceData{}
	}
	return data
}

func loadEducationData(lang string) templates.EducationData {
	filename := fmt.Sprintf("data/education_%s.json", lang)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/education_en.json"
		file, err = os.Open(filename)
		if err != nil {
			log.Printf("Error loading fallback %s: %v", filename, err)
			return templates.EducationData{}
		}
	}
	defer file.Close()
	var data templates.EducationData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Printf("Error decoding %s: %v", filename, err)
		return templates.EducationData{}
	}
	return data
}

func loadProjectsData(lang string) templates.ProjectsData {
	filename := fmt.Sprintf("data/projects_%s.json", lang)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/projects_en.json"
		file, err = os.Open(filename)
		if err != nil {
			log.Printf("Error loading fallback %s: %v", filename, err)
			return templates.ProjectsData{}
		}
	}
	defer file.Close()
	var data templates.ProjectsData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Printf("Error decoding %s: %v", filename, err)
		return templates.ProjectsData{}
	}
	return data
}

func sendEmail(name, email, message string) error {
	from := "beetswouter@gmail.com" // Configure
	to := "beetswouter@gmail.com"
	smtpHost := "smtp.example.com" // Configure
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, "password", smtpHost) // Configure
	msg := fmt.Sprintf("Subject: Contact from %s\n\nName: %s\nEmail: %s\nMessage: %s", name, name, email, message)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}
