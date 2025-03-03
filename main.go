package main

import (
	"fmt"
	backend "forum/backend"
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	var fileName = "./frontend/template/home/index.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Erreur pendant le parsing", err)
		return
	}

	t.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	var fileName = "./frontend/template/home/registration/index.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Erreur pendant le parsing", err)
		return
	}

	t.Execute(w, nil)
}

func main() {
	backend.InitDB()
	fs := http.FileServer(http.Dir("./frontend/public/"))
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))
	http.ListenAndServe("", nil)
}