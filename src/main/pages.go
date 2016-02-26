package main

import (
	//"encoding/json"
	"fmt"
	"html/template"
	//"log"
	"net/http"

	"github.com/gorilla/mux"
)

/* Index:
Página principal.
*/
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

/* Developers:
Mostra página que lista (paginado) todos os desenvolvedores.
*/
func showAllDevelopers(w http.ResponseWriter, r *http.Request) {
	devs := Developers{
		Developer{Name: "Tufts"},
		Developer{Name: "Caks"},
	}
	//d := Developer{Name: "Tufts"}
	t, err := template.ParseFiles("templates/developers.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(w, devs)
	}

	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(d); err != nil {
	//	panic(err)
	//}
}

/* showDeveloper:
Mostra página de desenvolvedor individual
*/
func showDeveloper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	devId := vars["devId"]
	fmt.Fprintln(w, "Developer #", devId)
}
