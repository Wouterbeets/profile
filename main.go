package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a new ServeMux to handle requests
	mux := http.NewServeMux()

	// Define HTML template
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Test Webserver</title>
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
    </style>
</head>
<body>
    <div class="container">
        <h1>Simple Test Webserver</h1>
        <p>This is a basic web server for testing purposes.</p>
        
        <div class="request-info">
            <strong>Request Details:</strong><br>
            <span class="timestamp">Timestamp: {{.Timestamp}}</span><br>
            Method: {{.Method}}<br>
            Path: {{.Path}}<br>
            Protocol: {{.Proto}}<br>
            Remote Address: {{.RemoteAddr}}<br>
            User Agent: {{.EscapedUserAgent}}
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
