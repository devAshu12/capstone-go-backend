package utils

import (
	"encoding/json"
	"errors"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"log"
	"net/http"
	"strings"
)

func RespondWithError(w http.ResponseWriter, err error) {

	log.Println("Error:", err)

	var appErr *types.AppError
	if ok := errors.As(err, &appErr); ok {
		w.WriteHeader(appErr.Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data := strings.Split(err.Error(), ",")

	resp := map[string]string{
		"err_code":    data[0],
		"err_message": data[1],
		"error":       data[2],
	}

	json.NewEncoder(w).Encode(resp)
}
