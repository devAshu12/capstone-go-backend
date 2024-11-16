package types

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
