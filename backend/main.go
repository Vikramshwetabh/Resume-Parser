package main

import (
	"fmt"      //for formatting output
	"io"       //for reading and writing data
	"log"      //for logging errors and information
	"net/http" //for handling HTTP requests
	"os"
	"os/exec" // for executing the Python parser
	"path/filepath"
)

// Serves static files (frontend) from ../frontend directory
func staticFileServer() http.Handler {
	frontendPath := filepath.Join("..", "frontend")
	fs := http.FileServer(http.Dir(frontendPath))
	return http.StripPrefix("/", fs)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Only POST method is allowed")
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing form data")
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error retrieving the file")
		return
	}
	defer file.Close()

	// Save the uploaded file to a temp location
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error creating temp file")
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error saving file")
		return
	}

	// Call Python parser here and return its output
	pythonPath := "python" // Assumes 'python' is in PATH
	parserPath := "backend/parser.py"
	cmd := exec.Command(pythonPath, parserPath, tempFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error running parser: %v\nOutput: %s", err, string(output))
		os.Remove(tempFile.Name())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	os.Remove(tempFile.Name())
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", staticFileServer()) // Serve frontend files

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
