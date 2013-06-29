package controllers

import (
	"net/http"
	"html/template"
)

type AuthItem struct {
	URL string
	Name string
}

func Index(w http.ResponseWriter, r *http.Request) {
	renderHomePage(w)
}

func isLogged() bool {
	if currentUser != nil {
		return true
	} 
	
	return false
}

func renderHomePage(w http.ResponseWriter) {
	funcs := template.FuncMap{"isLogged": isLogged} 
	t := template.Must(template.New("tmpl_index").Funcs(funcs).ParseFiles("templates/index.html", "templates/footer.html"))
	
	if err := t.Execute(w, currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
