package main

import (
	"errors"
	"fmt"
	"net/http"

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
		HasLyrics     bool               `json:"has_lyrics"`
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
		HasLyrics:     input.HasLyrics,
	}

	v := validator.New()

	if data.ValidateTune(v, tune); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tunes.Insert(tune)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tunes/%d", tune.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"tune": tune}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune, err := app.models.Tunes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tune, err := app.models.Tunes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title         *string             `json:"title"`
		Styles        []string            `json:"styles"`
		Keys          []data.Key          `json:"keys"`
		TimeSignature *data.TimeSignature `json:"time_signature"`
		Structure     *string             `json:"structure"`
		HasLyrics     *bool               `json:"has_lyrics"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		tune.Title = *input.Title
	}

	if input.Styles != nil {
		tune.Styles = input.Styles
	}

	if input.Keys != nil {
		tune.Keys = input.Keys
	}

	if input.TimeSignature != nil {
		tune.TimeSignature = *input.TimeSignature
	}

	if input.Structure != nil {
		tune.Structure = *input.Structure
	}

	if input.HasLyrics != nil {
		tune.HasLyrics = *input.HasLyrics
	}

	v := validator.New()

	if data.ValidateTune(v, tune); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tunes.Update(tune)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tune": tune}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTuneHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Tunes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "tune successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
