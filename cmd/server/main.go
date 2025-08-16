package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	chunkNumber := r.URL.Query().Get("chunk")

	if filePath == "" || chunkNumber == "" {
		http.Error(w, "Missing path or chunk number", http.StatusBadRequest)
		return
	}

	log.Printf(filePath)
	saveDir := filepath.Join("./uploaded_files", filepath.Dir(filePath))
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savePath := filepath.Join("./uploaded_files", filePath)
	f, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Chunk %s for file %s uploaded successfully", chunkNumber, filePath)
}

func checkhealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"health":  "up to date",
	})
}

func main() {
	http.HandleFunc("/uploads", uploadHandler)
	http.HandleFunc("/health", checkhealthHandler)
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
