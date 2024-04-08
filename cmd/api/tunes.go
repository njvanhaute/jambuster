package main

import (
	"fmt"
	"net/http"
)

func (app *application) createTuneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new tune")
}

func (app *application) showTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of tune %d\n", id)
}
