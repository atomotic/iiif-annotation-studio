package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
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
	log.Logger = log.Output(os.Stdout)
	db = InitDB(annotationsDB)
	router := httprouter.New()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Error().Msg("server err")
	}
	defer ln.Close()
	go func() {

		router.GET("/annotation/get/:id", Get)
		router.POST("/annotation/update/:id", Update)
		router.POST("/annotation/delete/:id", Delete)
		router.POST("/annotation/create", Create)
		router.GET("/annotation/list", List)

		fileServer := http.FileServer(assetFS())
		router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			r.URL.Path = ps.ByName("filepath")
			fileServer.ServeHTTP(w, r)
		})

		log.Info().Msg("# listening on " + ln.Addr().String())
		if err := http.Serve(ln, router); err != nil {
			log.Fatal().Err(err).Msg("Startup failed")
		}

	}()

	w := webview.New(webview.Settings{
		Width:     1000,
		Height:    700,
		Title:     "IIIF Annotation Studio",
		Resizable: true,
		URL:       "http://" + ln.Addr().String() + "/static/index.html",
		Debug:     true,
	})
	defer w.Exit()
	w.Run()
}
