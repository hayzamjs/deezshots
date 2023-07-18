package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func saveUploadedFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	randomFilename := generateRandomHash() + ".png"

	imgFolder := "./img"
	if _, err := os.Stat(imgFolder); os.IsNotExist(err) {
		os.Mkdir(imgFolder, 0755)
	}

	dst, err := os.Create(imgFolder + "/" + randomFilename)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return randomFilename, nil
}

func generateRandomHash() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Error generating random hash: ", err)
	}
	return hex.EncodeToString(bytes)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse the form", http.StatusInternalServerError)
		return
	}

	key, ok := r.Form["api_key"]

	if !ok {
		http.Error(w, "API key not found", http.StatusBadRequest)
		return
	}

	if key[0] != "deez_nuts" {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	file, fileHeader, err := r.FormFile("screenshot")
	if err != nil {
		http.Error(w, "Failed to retrieve the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := saveUploadedFile(fileHeader)
	if err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "http://localhost:8080/view/%s", filename)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Path[len("/view/"):]
	http.ServeFile(w, r, "./img/"+imageName)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view/", viewHandler)

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
