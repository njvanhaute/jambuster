package main

import (
	"fmt"
	"net/http"
	"time"

	"jambuster.njvanhaute.com/internal/data"
	"jambuster.njvanhaute.com/internal/validator"
)

func (app *application) createTuneHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string             `json:"title"`
		Styles        []string           `json:"styles"`
		Keys          []data.Key         `json:"keys"`
		TimeSignature data.TimeSignature `json:"time_signature"`
		Structure     string             `json:"structure"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tune := &data.Tune{
		Title:         input.Title,
		Styles:        input.Styles,
		Keys:          input.Keys,
		TimeSignature: input.TimeSignature,
		Structure:     input.Structure,
	}

	v := validator.New()

	if data.ValidateTune(v, tune); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
		Keys:          []data.Key{data.Key("G major")},
		TimeSignature: data.TimeSignature("2/2"),
		Structure:     "AABB",
		Version:       1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
