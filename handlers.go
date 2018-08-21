package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/lithammer/shortuuid"
)

func AnnotationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var body string
		var list []string
		q, _ := r.URL.Query()["q"]
		rows, err := db.Query("SELECT body FROM annotations where target=?", q[0])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&body)
			if err != nil {
				log.Fatal(err)
			}
			list = append(list, body)
		}
		JoinFunc := template.FuncMap{"StringsJoin": strings.Join}
		tmpl := template.Must(template.New("").Funcs(JoinFunc).Parse(AnnotationListTemplate))
		tmpl.Execute(w, list)

	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var annotation Annotation
		err = json.Unmarshal(body, &annotation)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		annoid := shortuuid.New()
		annotation.ID = "http://annotation-studio.loc/annotation/" + annoid
		AnnotationWithID, _ := json.Marshal(annotation)
		statement, _ := db.Prepare("INSERT INTO annotations (annoid, created_at, target, manifest, body) VALUES (?, ?, ?, ?, ?)")
		statement.Exec(annoid, time.Now(), annotation.Canvas(), annotation.Manifest(), AnnotationWithID)
		fmt.Fprintf(w, string(AnnotationWithID))

	case http.MethodPut:
		http.Error(w, "NOT IMPLEMENTED", 501)

	case http.MethodDelete:
		http.Error(w, "NOT IMPLEMENTED", 501)

	}
}
