package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
)

// .Funcs(template.FuncMap{"markDown": markDowner})
// var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var templates = template.Must(template.New("").Funcs(template.FuncMap{"markDown": markDowner}).ParseFiles("tmpl/edit.html", "tmpl/view.html"))
var validPath = regexp.MustCompile("(^/edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

// This function receives a pointer to a Page struct.
// Writes it to a file, and dumps.
func (p *Page) save() error {
	filename := "files/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Take a page title, and form a filename from it.
// Read that file into a Page struct, and pass the
// pointer back.
func loadPage(title string) (*Page, error) {
	filename := "files/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Page handling stuff done.

// Main func
func main() {
	fmt.Println(reflect.TypeOf(*templates).Kind())
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)

}
