package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Create a new Chi router
	router := chi.NewRouter()

	// Define HTML template with HTMX support
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My CV Portfolio</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 1000px;
            margin: 0 auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            min-height: 100vh;
        }
        .container {
            background-color: white;
            padding: 40px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
            margin-top: 20px;
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            font-size: 2.5em;
            margin-bottom: 30px;
            position: relative;
        }
        h1:after {
            content: '';
            display: block;
            width: 100px;
            height: 4px;
            background: linear-gradient(45deg, #667eea, #764ba2);
            margin: 10px auto;
            border-radius: 2px;
        }
        .request-info {
            background-color: #e9f7fe;
            padding: 20px;
            border-radius: 10px;
            margin: 25px 0;
            font-family: monospace;
            border-left: 5px solid #667eea;
        }
        .timestamp {
            color: #666;
            font-size: 0.9em;
        }
        .section {
            margin: 30px 0;
            padding: 25px;
            border-radius: 10px;
            background: #f8f9fa;
            box-shadow: 0 2px 10px rgba(0,0,0,0.05);
        }
        .section h2 {
            color: #2c3e50;
            margin-top: 0;
            padding-bottom: 10px;
            border-bottom: 2px solid #667eea;
        }
        .experience-item, .education-item {
            margin: 20px 0;
            padding: 20px;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 3px 10px rgba(0,0,0,0.08);
            border-left: 4px solid #667eea;
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        .experience-item:hover, .education-item:hover {
            transform: translateY(-5px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        .experience-item h3, .education-item h3 {
            color: #2c3e50;
            margin-top: 0;
        }
        .experience-item .company, .education-item .institution {
            font-weight: bold;
            color: #667eea;
            display: block;
            margin-bottom: 8px;
        }
        .experience-item .period, .education-item .period {
            color: #7f8c8d;
            font-size: 0.9em;
            margin-bottom: 10px;
        }
        .contact-info {
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
            margin-top: 20px;
        }
        .contact-item {
            flex: 1;
            min-width: 200px;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 8px;
            text-align: center;
            box-shadow: 0 2px 5px rgba(0,0,0,0.05);
        }
        .contact-item i {
            font-size: 1.5em;
            margin-bottom: 10px;
            color: #667eea;
        }
        .btn {
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            padding: 12px 25px;
            border: none;
            border-radius: 30px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
            margin: 10px 5px;
            font-weight: bold;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
        }
        .btn:hover {
            transform: translateY(-3px);
            box-shadow: 0 6px 20px rgba(0,0,0,0.25);
        }
        .skills-container {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin-top: 20px;
        }
        .skill-tag {
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            padding: 8px 15px;
            border-radius: 20px;
            font-size: 0.9em;
        }
        .profile-summary {
            text-align: center;
            margin: 30px 0;
            padding: 20px;
            background: linear-gradient(45deg, #667eea, #764ba2);
            color: white;
            border-radius: 10px;
        }
        .profile-summary p {
            font-size: 1.1em;
            line-height: 1.6;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Professional Portfolio</h1>
        
        <div class="request-info">
            <strong>Request Details:</strong><br>
            <span class="timestamp">Timestamp: {{.Timestamp}}</span><br>
            Method: {{.Method}}<br>
            Path: {{.Path}}<br>
            Protocol: {{.Proto}}<br>
            Remote Address: {{.RemoteAddr}}<br>
            User Agent: {{.EscapedUserAgent}}
        </div>

        <div class="profile-summary">
            <p>Passionate software engineer with expertise in building scalable web applications and leading technical teams. 
            Specialized in Go, full-stack development, and modern architecture patterns.</p>
        </div>

        <div class="section" hx-get="/cv/experience" hx-target="#experience-content" hx-trigger="load">
            <h2>Professional Experience</h2>
            <div id="experience-content">
                <!-- Experience will be loaded here via HTMX -->
                <p>Loading experience...</p>
            </div>
        </div>

        <div class="section" hx-get="/cv/education" hx-target="#education-content" hx-trigger="load">
            <h2>Education</h2>
            <div id="education-content">
                <!-- Education will be loaded here via HTMX -->
                <p>Loading education...</p>
            </div>
        </div>

        <div class="section">
            <h2>Contact Information</h2>
            <div class="contact-info">
                <div class="contact-item">
                    <i>📧</i>
                    <div><strong>Email:</strong> john.doe@example.com</div>
                </div>
                <div class="contact-item">
                    <i>📱</i>
                    <div><strong>Phone:</strong> +1 (123) 456-7890</div>
                </div>
                <div class="contact-item">
                    <i>📍</i>
                    <div><strong>Location:</strong> San Francisco, CA</div>
                </div>
            </div>
        </div>

        <div class="section">
            <h2>Skills</h2>
            <div class="skills-container">
                {{range .Skills}}
                <span class="skill-tag">{{.}}</span>
                {{end}}
            </div>
        </div>

        <p>Request received and logged successfully!</p>
    </div>
</body>
</html>
`

	// Parse the template
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

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
			"Timestamp":      time.Now().Format(time.RFC3339),
			"Method":         r.Method,
			"Path":           r.URL.Path,
			"Proto":          r.Proto,
			"RemoteAddr":     r.RemoteAddr,
			"UserAgent":      r.UserAgent(),
			"EscapedUserAgent": template.HTMLEscapeString(r.UserAgent()),
			"Skills": []string{
				"Go Programming", "Python", "C", "React", 
				"Microservices Architecture", "Cloud Technologies",
				"DevOps Practices", "Database Design", "System Design",
				"Leadership", "Team Management", "Technical Vision",
			},
		}

		// Execute template
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Handle experience section
	router.Get("/cv/experience", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		experienceHTML := `
        <div class="experience-item">
            <h3>CTO</h3>
            <span class="company">La Clinique E-Santé</span>
            <span class="period">February 2023 - January 2025, Paris, France</span>
            <ul>
                <li>Led the technical vision and architecture for an online mental health therapy platform offering 24/7 access via message, audio, and video</li>
                <li>Oversaw development and maintenance of mobile and web applications, ensuring secure, reimbursable consultations and seamless patient-psychologist communication</li>
                <li>Managed a team of developers to implement features for anxiety management, brief therapies (e.g., CBT, EMDR, hypnosis), and integrative mental health solutions</li>
                <li>Handled platform scalability and privacy compliance for e-health services, contributing to the company's mission of making mental health accessible until its closure</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Payment Platform Staff Engineer</h3>
            <span class="company">leboncoin</span>
            <span class="period">2022 - January 2023, Paris, France</span>
            <ul>
                <li>Construction de la vision architecturale technique de 4 équipes</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Lead Developer</h3>
            <span class="company">leboncoin</span>
            <span class="period">2019 - 2022, Paris, France</span>
            <ul>
                <li>Accompagnement organisationnel et technique de 100 developers</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Backend Developer</h3>
            <span class="company">leboncoin</span>
            <span class="period">2017 - 2019, Paris, France</span>
            <ul>
                <li>Migration technique du codebase legacy vers une architecture microservices</li>
                <li>Intégration des nouveaux payment service providers</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Fullstack Developer</h3>
            <span class="company">Artefact</span>
            <span class="period">2015 - 2017, Paris, France</span>
            <ul>
                <li>Gestion d’un outil d’analyse big data</li>
                <li>Scrum Master de l’équipe produit</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Entrepreneur service à domicile en restauration</h3>
            <span class="company">Thuis aan Tafel - Pays-Bas</span>
            <span class="period">2012 - 2015, Pays-Bas</span>
            <ul>
                <li>Création et mise à jour d’un logiciel sous MS ACCESS</li>
                <li>Organisation comptable et responsable financier</li>
            </ul>
        </div>
        `
		w.Write([]byte(experienceHTML))
	})

	// Handle education section
	router.Get("/cv/education", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		educationHTML := `
        <div class="education-item">
            <h3>Grande École Numérique</h3>
            <span class="institution">École 42</span>
            <span class="period">September 2013 - 2016, RNCP Niveau 1</span>
        </div>
        <div class="education-item">
            <h3>Communication and Multimedia Design</h3>
            <span class="institution">Hogeschool van Amsterdam</span>
            <span class="period">September 2007 - 2008</span>
        </div>
        <div class="education-item">
            <h3>Engineering, Design and Innovation</h3>
            <span class="institution">Hogeschool van Amsterdam</span>
            <span class="period">September 2006 - 2007</span>
        </div>
        <div class="education-item">
            <h3>Hoger Algemeen Voortgezet Onderwijs</h3>
            <span class="institution">Equivalent du baccalauréat</span>
            <span class="period">September 2001 - 2006</span>
        </div>
        `
		w.Write([]byte(educationHTML))
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
