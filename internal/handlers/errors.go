package handlers

import (
	"net/http"

	"github.com/jtorre/qisurChallenge/internal/utils"
)

func RespondError(w http.ResponseWriter, err error) {
	if _, ok := err.(*utils.ValidationError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func RespondNotFound(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusNotFound)
}

func RespondUnauthorized(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusUnauthorized)
}

func RespondForbidden(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusForbidden)
}

func RespondConflict(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusConflict)
}
