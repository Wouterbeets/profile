package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Create a new ServeMux to handle requests
	mux := http.NewServeMux()

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
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .request-info {
            background-color: #e9f7fe;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            font-family: monospace;
        }
        .timestamp {
            color: #666;
            font-size: 0.9em;
        }
        .section {
            margin: 20px 0;
            padding: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        .section h2 {
            color: #444;
            margin-top: 0;
        }
        .experience-item, .education-item {
            margin: 10px 0;
            padding: 10px;
            background-color: #f9f9f9;
            border-left: 4px solid #007bff;
        }
        .contact-info {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
        }
        .contact-item {
            flex: 1;
            min-width: 200px;
        }
        .btn {
            background-color: #007bff;
            color: white;
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
            margin: 5px;
        }
        .btn:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>My CV Portfolio</h1>
        
        <div class="request-info">
            <strong>Request Details:</strong><br>
            <span class="timestamp">Timestamp: {{.Timestamp}}</span><br>
            Method: {{.Method}}<br>
            Path: {{.Path}}<br>
            Protocol: {{.Proto}}<br>
            Remote Address: {{.RemoteAddr}}<br>
            User Agent: {{.EscapedUserAgent}}
        </div>

        <div class="section">
            <h2>About Me</h2>
            <p>Hello! I'm a passionate software developer with experience in building web applications using Go, JavaScript, and modern frameworks.</p>
        </div>

        <div class="section" hx-get="/cv/experience" hx-target="#experience-content" hx-trigger="load">
            <h2>Work Experience</h2>
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
                    <strong>Email:</strong> john.doe@example.com
                </div>
                <div class="contact-item">
                    <strong>Phone:</strong> +1 (123) 456-7890
                </div>
                <div class="contact-item">
                    <strong>Location:</strong> San Francisco, CA
                </div>
            </div>
        </div>

        <div class="section">
            <h2>Skills</h2>
            <ul>
                <li>Go Programming</li>
                <li>Web Development (HTML/CSS/JS)</li>
                <li>Database Design</li>
                <li>Cloud Technologies</li>
                <li>DevOps Practices</li>
            </ul>
        </div>

        <p>Request received and logged successfully!</p>
    </div>
</body>
</html>
`

	// Parse the template
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	// Handle all requests with a logging middleware
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

		// Prepare data for template
		data := map[string]interface{}{
			"Timestamp":      time.Now().Format(time.RFC3339),
			"Method":         r.Method,
			"Path":           r.URL.Path,
			"Proto":          r.Proto,
			"RemoteAddr":     r.RemoteAddr,
			"UserAgent":      r.UserAgent(),
			"EscapedUserAgent": template.HTMLEscapeString(r.UserAgent()),
		}

		// Execute template
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Handle experience section
	mux.HandleFunc("/cv/experience", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		experienceHTML := `
        <div class="experience-item">
            <h3>Senior Software Engineer</h3>
            <p><strong>ABC Company</strong> | Jan 2020 - Present</p>
            <ul>
                <li>Lead development of microservices architecture using Go</li>
                <li>Improved system performance by 40%</li>
                <li>Mentored junior developers and conducted code reviews</li>
            </ul>
        </div>
        <div class="experience-item">
            <h3>Software Developer</h3>
            <p><strong>XYZ Solutions</strong> | Jun 2017 - Dec 2019</p>
            <ul>
                <li>Developed web applications using React and Node.js</li>
                <li>Implemented CI/CD pipelines</li>
                <li>Collaborated with UX team to improve user experience</li>
            </ul>
        </div>
        `
		w.Write([]byte(experienceHTML))
	})

	// Handle education section
	mux.HandleFunc("/cv/education", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		educationHTML := `
        <div class="education-item">
            <h3>M.S. Computer Science</h3>
            <p><strong>University of Technology</strong> | 2015 - 2017</p>
            <p>Specialized in Distributed Systems and Web Technologies</p>
        </div>
        <div class="education-item">
            <h3>B.S. Software Engineering</h3>
            <p><strong>State University</strong> | 2011 - 2015</p>
            <p>Graduated with honors, GPA: 3.8/4.0</p>
        </div>
        `
		w.Write([]byte(educationHTML))
	})

	// Create server
	server := &http.Server{
		Addr:    ":33333",
		Handler: mux,
	}

	// Start server and log any errors
	log.Printf("Starting server on port 33333...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
