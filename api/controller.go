package api

import (
	"encoding/json"
	"net/http"
)

func pingController(w http.ResponseWriter, r *http.Request) {
	var ping struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ping); err != nil {
		message := map[string]string{
			"error": "can not decode payload",
		}
		respondJSON(w, http.StatusBadRequest, message)
		return
	}
	defer r.Body.Close()

	response, err := PingUrl(ping.URL)
	if err != nil {
		message := map[string]string{
			"error": err.Error(),
		}
		respondJSON(w, http.StatusOK, message)
		return
	}

	respondJSON(w, http.StatusOK, response)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
