package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// Index Page
func Home(w http.ResponseWriter, r *http.Request) {
	var tmpls = template.Must(template.ParseFiles("templates/home.html"))
	data := struct {
		Title  string
		Header string
	}{
		Title:  "Index Page",
		Header: "Hello, World!",
	}

	if err := tmpls.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Signup Page
func Signup(w http.ResponseWriter, r *http.Request) {
	var tmpls = template.Must(template.ParseFiles("templates/signup.html"))
	data := struct {
		Title  string
		Header string
	}{
		Title:  "Index Page",
		Header: "Hello, World!",
	}

	if err := tmpls.ExecuteTemplate(w, "signup.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Login Page
func Login(w http.ResponseWriter, r *http.Request) {
	var tmpls = template.Must(template.ParseFiles("templates/login.html"))
	data := struct {
		Title  string
		Header string
	}{
		Title:  "Index Page",
		Header: "Hello, World!",
	}

	if err := tmpls.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Alums Page
func Alums(w http.ResponseWriter, r *http.Request) {
	var tmpls = template.Must(template.ParseFiles("templates/alums.html"))
	data := struct {
		Title  string
		Header string
	}{
		Title:  "Index Page",
		Header: "Hello, World!",
	}

	if err := tmpls.ExecuteTemplate(w, "alums.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/signup", Signup)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/alums", Alums)
	r.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir("templates/"))))

	http.Handle("/", r)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}