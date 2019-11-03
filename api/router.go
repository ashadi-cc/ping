package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//InitServer just init sever
func InitServer() {
	router := chi.NewMux()
	setupRouter(router)
	log.Println("Serve listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func setupRouter(router *chi.Mux) {
	setupMiddleware(router)
	router.Post("/ping", pingController)
}

func setupMiddleware(router *chi.Mux) {
	router.Use(
		//render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
		middleware.Logger,          //log api request calls
	)
}
