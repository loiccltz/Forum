package backend

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const MaxUploadSize = 20 << 20 // 20 MB

func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "L'image est trop grande (max 20MB).", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Erreur lors de l'upload du fichier.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if !isValidImageExtension(ext) {
		http.Error(w, "Format d'image non supporté. Utilisez JPG, PNG ou GIF.", http.StatusBadRequest)
		return
	}

	filename := fmt.Sprintf("%d%s", os.Getpid(), ext)
	filePath := filepath.Join("uploads", filename)

	if err := saveFile(file, filePath); err != nil {
		http.Error(w, "Impossible de sauvegarder l'image.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("/uploads/" + filename))
}

func isValidImageExtension(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}

func saveFile(file multipart.File, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = out.ReadFrom(file)
	return err
}
