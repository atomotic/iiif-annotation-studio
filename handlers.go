package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"hash/fnv"

	"github.com/julienschmidt/httprouter"
	"github.com/lithammer/shortuuid"
	"github.com/rs/zerolog/log"
)

// Simple hash function used to generate a string hash of the canvas uri
// to have a unique stable annotation list id
func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 16)
}

// List display the annotation list
func List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body json.RawMessage
	canvas, _ := r.URL.Query()["canvas"]

	annotationlisthash := hash(canvas[0])

	list := AnnotationList{
		Context: "http://iiif.io/api/presentation/2/context.json",
		ID:      "https://docuver.se/iiif/annotation/list" + annotationlisthash,
		Type:    "sc:AnnotationList",
	}

	rows, err := db.Query("SELECT body FROM annotations where target=?", canvas[0])
	if err != nil {
		log.Error().Err(err).Str("action", "list").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&body)
		if err != nil {
			log.Error().Err(err).Str("action", "list").Msg("db error")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			return
		}
		list.Resources = append(list.Resources, body)
	}

	if len(list.Resources) == 0 {
		list.Resources = append(list.Resources, json.RawMessage("{}"))
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(list); err != nil {
		log.Error().Err(err).Str("action", "list").Str("canvas", canvas[0]).Msg("annotation list error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

}

// Get retrieve a single annotation with database id
// NOTE: just for test, not used by mirador
func Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var body string
	err := db.QueryRow("SELECT body FROM annotations where id=?", ps.ByName("id")).Scan(&body)

	if err != nil {
		log.Error().Err(err).Str("action", "get").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, body)
}

// Delete an annotation
// the output of this controller is the annotation body that is being deleted
func Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var anno string
	err := db.QueryRow("SELECT body FROM annotations where annoid = ?", ps.ByName("id")).Scan(&anno)
	if err != nil {
		log.Error().Err(err).Str("action", "delete").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	_, err = db.Exec("DELETE FROM annotations where annoid=?", ps.ByName("id"))
	if err != nil {
		log.Error().Err(err).Str("action", "delete").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}
	log.Info().Str("annotation_id", ps.ByName("id")).Str("action", "delete").Msg("")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, anno)
}

// Create an annotation
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error().Err(err).Str("action", "create").Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	var annotation Annotation
	err = json.Unmarshal(body, &annotation)
	if err != nil {
		log.Error().Err(err).Str("action", "create").Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	annoid := shortuuid.New()
	annotation.ID = "https://docuver.se/iiif/annotation/" + annoid
	AnnotationWithID, _ := json.Marshal(annotation)
	statement, _ := db.Prepare("INSERT INTO annotations (annoid, created_at, target, manifest, body) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(annoid, time.Now(), annotation.Canvas(), annotation.Manifest(), AnnotationWithID)
	log.Info().Str("annotation_id", annoid).Str("action", "create").Msg("")

	fmt.Fprintf(w, string(AnnotationWithID))
}

// Update an annotation
func Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	annoid := ps.ByName("id")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error().Err(err).Str("action", "update").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	_, err = db.Exec("UPDATE annotations SET body=? WHERE annoid=?", body, annoid)
	if err != nil {
		log.Error().Err(err).Str("action", "update").Msg("db error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	log.Info().Str("annotation_id", annoid).Str("action", "update").Msg("")

	fmt.Fprintf(w, string(body))
}
