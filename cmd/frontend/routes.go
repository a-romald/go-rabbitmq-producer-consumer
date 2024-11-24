package main

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(app.Home))
	mux.Handle("/handle", http.HandlerFunc(app.HandleForm))

	return mux
}
