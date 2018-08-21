package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/zserge/webview"
)

var home, _ = homedir.Dir()
var annotationsDirectory = filepath.Join(home, ".annotations")
var annotationsDB = filepath.Join(annotationsDirectory, "annotations.db")
var db *sql.DB

func init() {
	if _, err := os.Stat(annotationsDirectory); os.IsNotExist(err) {
		os.Mkdir(annotationsDirectory, 0755)
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", annotationsDB)
	if err != nil {
		log.Fatal("error db")
	}
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS annotations (
		id INTEGER PRIMARY KEY, 
		annoid VARCHAR, 
		created_at DATETIME, 
		target VARCHAR, 
		manifest VARCHAR, 
		body TEXT)`)
	statement.Exec()

	statement, _ = db.Prepare("CREATE UNIQUE INDEX IF NOT EXISTS annotation_id ON annotations (annoid);")
	statement.Exec()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		http.Handle("/", http.FileServer(assetFS()))
		http.Handle("/iiif/annotation", http.HandlerFunc(AnnotationHandler))
		fmt.Println("running on: http://" + ln.Addr().String())
		log.Fatal(http.Serve(ln, nil))
	}()

	w := webview.New(webview.Settings{
		Width:     1000,
		Height:    700,
		Title:     "IIIF Annotation Studio",
		Resizable: true,
		URL:       "http://" + ln.Addr().String(),
		Debug:     true,
	})
	// w.SetColor(255, 0, 0, 0)
	defer w.Exit()
	w.Run()
}
