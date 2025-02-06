package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, s string, p *Page) {
	t, err := template.ParseFiles(s + ".html")

	if err != nil {
		fmt.Print("Error parsefile view.html")
	}
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")

	if err != nil {
		fmt.Print("Error parsefile index.html")
	}
	t.Execute(w, "")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}

	renderTemplate(w, "view", p)
}

func edithandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)

	if err != nil {
		fmt.Print("error")
	}
	renderTemplate(w, "edit", p)
}

func savehandler(w http.ResponseWriter, r *http.Request) {
	titel := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: titel, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+titel, http.StatusFound)
}
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", edithandler)
	http.HandleFunc("/save/", savehandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
