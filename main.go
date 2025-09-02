package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"testserver/templates"
)

//go:embed data/*.json static/* manifest.json sw.js
var embeddedFS embed.FS

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
	// Parse command-line flags
	port := flag.String("p", "33333", "Port to run the server on")
	flag.Parse()

	// Create a new Chi router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(5)) // GZIP compression
	router.Use(middleware.Recoverer)

	// Static files with correct MIME types
	router.HandleFunc("/static/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		file, err := embeddedFS.Open("static/" + path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		reader := bytes.NewReader(content)
		// Set Content-Type based on extension
		if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.ServeContent(w, r, path, time.Time{}, reader)
	})

	router.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		file, err := embeddedFS.Open("manifest.json")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		reader := bytes.NewReader(content)
		w.Header().Set("Content-Type", "application/json")
		http.ServeContent(w, r, "manifest.json", time.Time{}, reader)
	})

	router.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		file, err := embeddedFS.Open("sw.js")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		reader := bytes.NewReader(content)
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeContent(w, r, "sw.js", time.Time{}, reader)
	})

	// Handle root route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		data := templates.IndexData{
			Skills: []string{
				"Golang", "Python", "C", "React",
				"Stripe", "HubSpot", "PostgreSQL", "Docker", "Kubernetes", "Git", "Agile Methodologies",
				"Dutch (Native)", "English (Fluent)", "French (Fluent)",
				"AI Integration", "Privacy-Conscious AI", "Event Sourcing", "Domain-Driven Design",
			},
			Profile:      loadProfileData(lang),
			Translations: templates.Translations,
			Language:     lang,
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

	// New: Handle experience detail
	router.Get("/cv/experience/detail/{id}", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		data := loadExperienceData(lang)
		if id >= len(data.ExperienceItems) {
			http.Error(w, "ID out of range", http.StatusBadRequest)
			return
		}
		item := data.ExperienceItems[id]
		templates.ExperienceDetailTemplate(item, lang, id).Render(r.Context(), w)
	})

	// New: Handle experience collapse
	router.Get("/cv/experience/collapse/{id}", func(w http.ResponseWriter, r *http.Request) {
		lang := detectLanguage(r)
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		data := loadExperienceData(lang)
		if id >= len(data.ExperienceItems) {
			http.Error(w, "ID out of range", http.StatusBadRequest)
			return
		}
		item := data.ExperienceItems[id]
		templates.ExperienceSummaryTemplate(item, id).Render(r.Context(), w)
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
			"Golang", "Python", "C", "React",
			"Stripe", "HubSpot", "PostgreSQL", "Docker", "Kubernetes", "Git", "Agile Methodologies",
			"Dutch (Native)", "English (Fluent)", "French (Fluent)",
			"AI Integration", "Privacy-Conscious AI", "Event Sourcing", "Domain-Driven Design",
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
		Addr:    ":" + *port,
		Handler: router,
	}

	log.Printf("Starting server on port %s...", *port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loadExperienceData(lang string) templates.ExperienceData {
	filename := fmt.Sprintf("data/experience_%s.json", lang)
	file, err := embeddedFS.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/experience_en.json"
		file, err = embeddedFS.Open(filename)
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
	file, err := embeddedFS.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/education_en.json"
		file, err = embeddedFS.Open(filename)
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
	file, err := embeddedFS.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/projects_en.json"
		file, err = embeddedFS.Open(filename)
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

func loadProfileData(lang string) templates.ProfileData {
	filename := fmt.Sprintf("data/profile_%s.json", lang)
	file, err := embeddedFS.Open(filename)
	if err != nil {
		log.Printf("Error loading %s: %v, falling back to English", filename, err)
		filename = "data/profile_en.json"
		file, err = embeddedFS.Open(filename)
		if err != nil {
			log.Printf("Error loading fallback %s: %v", filename, err)
			return templates.ProfileData{}
		}
	}
	defer file.Close()
	var data templates.ProfileData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Printf("Error decoding %s: %v", filename, err)
		return templates.ProfileData{}
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
