package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	style := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", style))

	http.HandleFunc("/", handler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/export", exportHandler) // New endpoint for exporting ASCII art
	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(writing http.ResponseWriter, response *http.Request) {
	templ, _ := template.ParseFiles("templates/index.html")
	writing.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templ.Execute(writing, nil)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	text := r.FormValue("text")
	bannerType := r.FormValue("banner-type")

	if r.FormValue("export") != "" {
		// If the "Export ASCII Art" button was clicked
		exportHandler(w, r)
		return
	}

	result := GenerateASCIIArt(text, bannerType)

	templ, err := template.ParseFiles("templates/index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = templ.Execute(w, struct{ Result string }{Result: result})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result := r.FormValue("exportResult")
	file, err := os.Create("exported_result.txt")
	if err != nil {
		http.Error(w, "Error creating export file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.WriteString(result)
	if err != nil {
		http.Error(w, "Error writing to export file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=exported_result.txt")
	http.ServeFile(w, r, "exported_result.txt")
}
