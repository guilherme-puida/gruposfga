package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	_ "embed"
)

var port = flag.String("port", ":8080", "the port that the server will listen on")
var data = flag.String("data", "data.json", "data file")

//go:embed index.html
var indexPage string

type item struct {
	Materia   string `json:"materia"`
	Turma     string `json:"turma"`
	Professor string `json:"professor"`
	Horario   string `json:"horario"`
	Link      string `json:"link"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeItems(items []item) {
	b, err := json.Marshal(items)
	check(err)
	os.WriteFile(*data, b, 0644)
}

func updateItems(items *[]item, materia string, turma string, link string) {
	for i, val := range *items {
		if val.Materia == materia && val.Turma == turma {
			(*items)[i].Link = link
			go writeItems(*items)
		}
	}
}

func main() {
	flag.Parse()

	indexPageTemplate := template.Must(template.New("index").Parse(indexPage))

	dat, err := os.ReadFile(*data)
	check(err)

	var items []item
	err = json.Unmarshal(dat, &items)
	check(err)

	sort.Slice(items, func(i, j int) bool {
		a := items[i]
		b := items[j]

		if a.Materia == b.Materia {
			return a.Turma < b.Turma
		}

		return a.Materia < b.Materia
	})

	handleGet := func(w http.ResponseWriter, r *http.Request) {
		indexPageTemplate.Execute(w, items)
	}

	handlePost := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		materia := r.FormValue("Materia")
		turma := r.FormValue("Turma")
		link := r.FormValue("Link")

		updateItems(&items, materia, turma, link)
		http.Redirect(w, r, "/", http.StatusFound)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGet(w, r)
		} else if r.Method == http.MethodPost {
			handlePost(w, r)
		}
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	log.Printf("Listening on %s\n", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}
