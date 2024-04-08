package main

import (
	"fmt"
	"net/http"
	"time"

	"jambuster.njvanhaute.com/internal/data"
)

func (app *application) createTuneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new tune")
}

func (app *application) showTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune := data.Tune{
		ID:            id,
		CreatedAt:     time.Now(),
		Title:         "Roanoke",
		Styles:        []string{"Bluegrass"},
		Key:           "G major",
		TimeSignature: "2/2",
		Structure:     "AABB",
		Version:       1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
