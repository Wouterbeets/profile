package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Create a new Chi router
	router := chi.NewRouter()

	// Middleware to log requests
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log request details
			log.Printf(
				"[%s] %s %s %s from %s (User-Agent: %s)",
				time.Now().Format(time.RFC3339),
				r.Method,
				r.URL.Path,
				r.Proto,
				r.RemoteAddr,
				r.UserAgent(),
			)
			next.ServeHTTP(w, r)
		})
	})

	// Handle root route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Prepare data for template
		data := map[string]interface{}{
			"Timestamp":        time.Now().Format(time.RFC3339),
			"Method":           r.Method,
			"Path":             r.URL.Path,
			"Proto":            r.Proto,
			"RemoteAddr":       r.RemoteAddr,
			"UserAgent":        r.UserAgent(),
			"EscapedUserAgent": r.UserAgent(), // Templ handles escaping automatically
			"Skills": []string{
				"Go Programming", "Python", "C", "React",
				"Microservices Architecture", "Cloud Technologies",
				"DevOps Practices", "Database Design", "System Design",
				"Leadership", "Team Management", "Technical Vision",
			},
		}

		// Execute template
		templates.IndexTemplate(data).Render(r.Context(), w)
	})

	// Handle experience section
	router.Get("/cv/experience", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		experienceData := map[string]interface{}{
			"ExperienceItems": []map[string]interface{}{
				{
					"Title":   "CTO",
					"Company": "La Clinique E-Santé",
					"Period":  "February 2023 - January 2025, Paris, France",
					"Description": []string{
						"Led the technical vision and architecture for an online mental health therapy platform offering 24/7 access via message, audio, and video",
						"Oversaw development and maintenance of mobile and web applications, ensuring secure, reimbursable consultations and seamless patient-psychologist communication",
						"Managed a team of developers to implement features for anxiety management, brief therapies (e.g., CBT, EMDR, hypnosis), and integrative mental health solutions",
						"Handled platform scalability and privacy compliance for e-health services, contributing to the company's mission of making mental health accessible until its closure",
					},
				},
				{
					"Title":   "Payment Platform Staff Engineer",
					"Company": "leboncoin",
					"Period":  "2022 - January 2023, Paris, France",
					"Description": []string{
						"Defined the technical architecture vision for 4 teams",
					},
				},
				{
					"Title":   "Lead Developer",
					"Company": "leboncoin",
					"Period":  "2019 - 2022, Paris, France",
					"Description": []string{
						"Provided organizational and technical support to 100 developers",
					},
				},
				{
					"Title":   "Backend Developer",
					"Company": "leboncoin",
					"Period":  "2017 - 2019, Paris, France",
					"Description": []string{
						"Migrated the legacy codebase to a microservices architecture",
						"Integrated new payment service providers",
					},
				},
				{
					"Title":   "Fullstack Developer",
					"Company": "Artefact",
					"Period":  "2015 - 2017, Paris, France",
					"Description": []string{
						"Managed a big data analytics tool",
						"Served as Scrum Master for the product team",
					},
				},
				{
					"Title":   "Home Cooking Service Entrepreneur",
					"Company": "Thuis aan Tafel - Netherlands",
					"Period":  "2012 - 2015, Netherlands",
					"Description": []string{
						"Created and maintained a software solution using MS ACCESS",
						"Managed accounting and financial responsibility",
					},
				},
			},
		}

		ExperienceTemplate(experienceData).Render(r.Context(), w)
	})

	// Handle education section
	router.Get("/cv/education", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		educationData := map[string]interface{}{
			"EducationItems": []map[string]interface{}{
				{
					"Title":       "Grande École Numérique",
					"Institution": "École 42",
					"Period":      "September 2013 - 2016, RNCP Level 1",
				},
				{
					"Title":       "Communication and Multimedia Design",
					"Institution": "Hogeschool van Amsterdam",
					"Period":      "September 2007 - 2008",
				},
				{
					"Title":       "Engineering, Design and Innovation",
					"Institution": "Hogeschool van Amsterdam",
					"Period":      "September 2006 - 2007",
				},
				{
					"Title":       "Hoger Algemeen Voortgezet Onderwijs",
					"Institution": "Equivalent to high school diploma",
					"Period":      "September 2001 - 2006",
				},
			},
		}

		EducationTemplate(educationData).Render(r.Context(), w)
	})

	// Create server
	server := &http.Server{
		Addr:    ":33333",
		Handler: router,
	}

	// Start server and log any errors
	log.Printf("Starting server on port 33333...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
