package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// Net/http stuff goes below:

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil //the title is the second subexpression
}

func markDowner(args ...interface{}) template.HTML {

	raw := []byte(fmt.Sprintf("%s", args...))
	fmt.Printf("%s", raw)

	// regexp to handle wiki-style link tags
	// matches any alpha string inside square brackets
	link := regexp.MustCompile("\\[[a-zA-Z]+\\]")
	// uses regexp to search slice s for matching regexes.
	// Not sure how ReplaceAllFunc works, but it looks like
	// it takes the main slice (s, in this case), uses a closure
	// (second arg?) to build what we're replacing, and returns the whole
	// shindig.
	s := link.ReplaceAllFunc(raw, func(a []byte) []byte {
		m := string(a[1 : len(a)-1])
		return []byte("<a href=\"/view/" + m + "\">" + m + "</a>")
	})

	t := blackfriday.MarkdownCommon(s)

	sanitized := bluemonday.UGCPolicy().SanitizeBytes(t)
	return template.HTML(sanitized)
}
