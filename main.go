package main

import (
	
	"fmt"
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

func main() {
	fs := http.FileServer(http.Dir("../assets/"))
	http.HandleFunc("/", home)
	http.Handle("frontend/public/css", http.StripPrefix("frontend/public/css", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))
	http.ListenAndServe("", nil)
}